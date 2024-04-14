// Copyright (c) 2022-2023 https://rasbora.openseawave.com
//
// This file is part of Rasbora Distributed Video Transcoding
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package videotranscoder

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/flosch/pongo2/v6"
	"gopkg.in/vansante/go-ffprobe.v2"
	"openseawave.com/rasbora/internal/config"
	"openseawave.com/rasbora/internal/data"
	"openseawave.com/rasbora/internal/database"
	"openseawave.com/rasbora/internal/filesystem"
	"openseawave.com/rasbora/internal/logger"
	"openseawave.com/rasbora/internal/utilities"
)

// FfmpegTranscoderEngine hold an instance
type FfmpegTranscoderEngine struct {
	Config                      *config.Config
	Logger                      *logger.Logger
	Database                    *database.Database
	FileSystem                  *filesystem.FileSystem
	_videoTranscoderQueue       string
	_callbackManagerQueue       string
	_inputVideoInformation      *ffprobe.ProbeData
	_queueable                  *data.Queueable
	_taskPayload                *data.Task
	_videoTranscoderWorkerID    string
	_temporaryWorkingPath       string
	_temporaryInputVideoFile    *data.File
	_temporaryProcessingLogFile *data.File
	_temporaryOutputVideoFiles  *[]data.File
	_sourceInputVideoFile       *data.File
	_finalProcessingLogFile     *data.File
	_finalOutputVideoFiles      *[]data.File
}

//go:embed handlers/*
var handlersFS embed.FS

// StarTranscoderEngine start ffmpeg video transcoder engine.
func (fte *FfmpegTranscoderEngine) StarTranscoderEngine(ctx context.Context) {

	// get video transcoder queue name
	fte._videoTranscoderQueue = fte.Config.GetString("Components.VideoTranscoding.Queue")

	// get callback manager queue name
	fte._callbackManagerQueue = fte.Config.GetString("Components.CallbackManager.Queue")

	// get video transcoder worker id
	fte._videoTranscoderWorkerID = fte.Config.GetString("Components.VideoTranscoding.UniqueID")

	// get video transcoder time interval for pooling new tasks.
	checkNewTaskInterval := fte.Config.GetInt("Components.VideoTranscoding.CheckNewTaskInterval")

	for {
		select {
		case <-ctx.Done():
			return
		default:

			time.Sleep(time.Duration(checkNewTaskInterval) * time.Second)

			queueableItem, err := fte.Database.Dequeue(fte._videoTranscoderQueue, fte._videoTranscoderWorkerID)
			fte._queueable = &queueableItem

			if err != nil {
				fte.Logger.Debug(
					"ffmpeg_transcoder_engine",
					"task queue is empty",
					map[string]interface{}{
						"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
					},
				)
				continue
			}

			fte.Logger.Success(
				"ffmpeg_transcoder_engine",
				"received new task",
				map[string]interface{}{
					"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
					"task_id":                    queueableItem.ID,
				},
			)

			fte.Logger.Debug(
				"ffmpeg_transcoder_engine",
				"received task payload",
				map[string]interface{}{
					"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
					"task_id":                    queueableItem.ID,
					"task_data":                  queueableItem,
				},
			)

			// prepare task for processing
			fte._prepareForProcessingTask()
		}
	}
}

// _prepareForProcessingTask prepare ffmpeg to transcode video base on task settings.
func (fte *FfmpegTranscoderEngine) _prepareForProcessingTask() {

	//recover from panic
	defer func() {
		if r := recover(); r != nil {
			jsonData, _ := json.Marshal(r)
			fte._failedTask(errors.New(string(jsonData)))
			fte.Logger.Error(
				"ffmpeg_transcoder_engine.prepare_for_processing_task",
				fmt.Sprintf("we got panic: %v", string(jsonData)),
				map[string]interface{}{
					"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
					"task_id":                    fte._queueable.ID,
				},
			)
		}
	}()

	// get task payload from queue item.
	jsonData, errJ := json.Marshal(fte._queueable.Payload)
	if errJ != nil {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.prepare_for_processing_task",
			fmt.Sprintf("cannot cast queueable to json string: %v", errJ.Error()),
			map[string]interface{}{
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"task_id":                    fte._queueable.ID,
			},
		)
		fte._failedTask(errJ)
		return
	}

	// set task payload
	errU := json.Unmarshal(jsonData, &fte._taskPayload)
	if errU != nil {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.prepare_for_processing_task",
			fmt.Sprintf("cannot cast queueable to task struct: %v", errJ.Error()),
			map[string]interface{}{
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"task_id":                    fte._queueable.ID,
			},
		)
		fte._failedTask(errU)
		return
	}

	// update task starting time.
	fte._taskPayload.StartedAt = time.Now().UnixMilli()

	// prepare a temporary working path.
	if err := fte._prepareTemporaryWorkingPath(); err != nil {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.prepare_for_processing_task",
			fmt.Sprintf("error when prepare a temporary working path: %v", err.Error()),
			map[string]interface{}{
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"task_id":                    fte._queueable.ID,
			},
		)
		fte._failedTask(err)
		return
	}

	// prepare a temporary processing file.
	if err := fte._prepareTemporaryProcessingLogFile(); err != nil {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.prepare_for_processing_task",
			fmt.Sprintf("error when prepare a temporary processing file: %v", err.Error()),
			map[string]interface{}{
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"task_id":                    fte._queueable.ID,
			},
		)
		fte._failedTask(err)
		return
	}

	// prepare a temporary output video files.
	if err := fte._prepareTemporaryOutputVideoFiles(); err != nil {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.prepare_for_processing_task",
			fmt.Sprintf("error when prepare a temporary output video files: %v", err.Error()),
			map[string]interface{}{
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"task_id":                    fte._queueable.ID,
			},
		)
		fte._failedTask(err)
		return
	}

	// prepare input video file.
	if err := fte._prepareInputVideoFile(); err != nil {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.prepare_for_processing_task",
			fmt.Sprintf("error when prepare a temporary input video file: %v", err.Error()),
			map[string]interface{}{
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"task_id":                    fte._queueable.ID,
			},
		)
		fte._failedTask(errors.New("we cannot make a copy of input video file at to temporary transcoder file"))
		return
	}

	// get all video information about input source
	if err := fte._readInputVideoInformation(); err != nil {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.prepare_for_processing_task",
			fmt.Sprintf("cannot read input video information: %v", err.Error()),
			map[string]interface{}{
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"task_id":                    fte._queueable.ID,
			},
		)
		fte._failedTask(err)
		return
	}

	// prepare a temporary input video file.
	if err := fte._transcodingInputVideoFile(); err != nil {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.prepare_for_processing_task",
			fmt.Sprintf("fail to transcode video files: %v", err.Error()),
			map[string]interface{}{
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"task_id":                    fte._queueable.ID,
			},
		)

		log, readLogError := os.ReadFile(fte._temporaryProcessingLogFile.FullPath())
		if readLogError == nil {
			fte.Logger.Debug(
				"ffmpeg_transcoder_engine.prepare_for_processing_task",
				fmt.Sprintf("cannot transcode input video: %v", string(log)),
				map[string]interface{}{
					"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
					"task_id":                    fte._queueable.ID,
				},
			)
		}

		fte._failedTask(fmt.Errorf("%v \n %v", err, readLogError))
		return
	}

	//if everything okay above then process task as success
	fte._successTask()
}

// _prepareTemporaryWorkingPath prepare a temporary working path.
func (fte *FfmpegTranscoderEngine) _prepareTemporaryWorkingPath() error {

	fte._temporaryWorkingPath = fte.Config.GetString("Filesystem.LocalStorage.Folders.TemporaryWorkingPath")

	if err := os.MkdirAll(filepath.Dir(fte._temporaryWorkingPath), os.ModePerm); err != nil {
		return err
	}

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.prepare_temporary_working_path",
		"temporary working path is ready",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"temporary_working_path":     fte._temporaryWorkingPath,
		},
	)

	return nil
}

// _prepareTemporaryProcessingLogFile prepare a temporary processing file.
func (fte *FfmpegTranscoderEngine) _prepareTemporaryProcessingLogFile() error {

	fte._temporaryProcessingLogFile = &data.File{
		FileMeta: map[string]interface{}{
			"task_id": fte._queueable.ID,
		},
		FileName: fmt.Sprintf("%v%v", fte._queueable.ID, ".log"),
		FilePath: filepath.Join(fte._temporaryWorkingPath),
	}

	_, err := os.Create(fte._temporaryProcessingLogFile.FullPath())
	if err != nil {
		return err
	}

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.prepare_temporary_processing_log_file",
		"temporary processing log file is ready",
		map[string]interface{}{
			"task_id":                       fte._queueable.ID,
			"video_transcoder_worker_id":    fte._videoTranscoderWorkerID,
			"temporary_processing_log_file": fte._temporaryProcessingLogFile.FullPath(),
		},
	)

	return nil
}

// _prepareTemporaryOutputVideoFiles prepare a temporary output video file.
func (fte *FfmpegTranscoderEngine) _prepareTemporaryOutputVideoFiles() error {

	var outputVideoFiles []data.File
	var args []map[string]interface{}
	for _, item := range fte._taskPayload.VideoTranscoder.Output.Args {
		file := data.File{
			FileMeta: map[string]interface{}{
				"task_id": fte._queueable.ID,
				"quality": item["quality"],
			},
			FileName: fmt.Sprintf(
				"%v_%v%v",
				fte._queueable.ID,
				item["quality"].(string),
				fte._taskPayload.VideoTranscoder.Output.Container,
			),
			FilePath: filepath.Join(fte._temporaryWorkingPath),
		}
		outputVideoFiles = append(outputVideoFiles, file)
		item["output"] = file
		args = append(args, item)
	}

	fte._temporaryOutputVideoFiles = &outputVideoFiles
	fte._taskPayload.VideoTranscoder.Output.Args = args

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.prepare_temporary_output_video_files",
		"temporary output video files is ready",
		map[string]interface{}{
			"task_id":                      fte._queueable.ID,
			"video_transcoder_worker_id":   fte._videoTranscoderWorkerID,
			"temporary_output_video_files": fte._temporaryOutputVideoFiles,
		},
	)

	return nil
}

// _prepareInputVideoFile prepare a temporary input video file.
// [IMPORTANT] If there is new file system add this method should be modified
// or this func should re-write to take new file system without change source code.
func (fte *FfmpegTranscoderEngine) _prepareInputVideoFile() error {

	fte.Logger.Info(
		"ffmpeg_transcoder_engine.prepare_input_video_file",
		"copying input source to transcoder working path",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)

	var fs *filesystem.FileSystem

	fte._sourceInputVideoFile = &data.File{
		FileMeta: map[string]interface{}{
			"task_id": fte._taskPayload.ID,
		},
		FileName: fte._taskPayload.VideoTranscoder.InputVideo.FileName,
		FilePath: fte._taskPayload.VideoTranscoder.InputVideo.FilePath,
	}

	fte._temporaryInputVideoFile = &data.File{
		FileName: fmt.Sprintf(
			"%v%v%v",
			fte._queueable.ID,
			"_input",
			filepath.Ext(fte._taskPayload.VideoTranscoder.InputVideo.FileName),
		),
		FilePath: filepath.Join(fte._temporaryWorkingPath),
	}

	//check file system type if its local file system
	if fte._taskPayload.VideoTranscoder.InputVideo.FileSystem == data.LocalFileSystemType {
		fs = filesystem.NewFileSystem(&filesystem.LocalFileSystem{})
	}

	//check file system type if its object file system
	if fte._taskPayload.VideoTranscoder.InputVideo.FileSystem == data.ObjectFileSystemType {
		objectClient, err := filesystem.NewObjectClient(*fte.Config)

		if err != nil {
			return errors.New("we cannot make a client for object filesystem")
		}

		fs = filesystem.NewFileSystem(&filesystem.ObjectFileSystem{
			Minio: objectClient,
		})
	}

	//validate file system type
	if !utilities.InSlice(fte._taskPayload.VideoTranscoder.InputVideo.FileSystem, []data.FileSystemType{
		data.ObjectFileSystemType,
		data.LocalFileSystemType,
	}) {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.prepare_input_video_file",
			"unknown filesystem type",
			map[string]interface{}{
				"task_id":                    fte._queueable.ID,
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"temporary_input_video_file": fte._temporaryInputVideoFile.FullPath(),
				"source_input_video_file":    fte._sourceInputVideoFile.FullPath(),
				"filesystem":                 fte._taskPayload.VideoTranscoder.InputVideo.FileSystem.String(),
			},
		)
	}

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.prepare_input_video_file",
		"filesystem has been selected",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"temporary_input_video_file": fte._temporaryInputVideoFile.FullPath(),
			"source_input_video_file":    fte._sourceInputVideoFile.FullPath(),
			"filesystem":                 fte._taskPayload.VideoTranscoder.InputVideo.FileSystem.String(),
		},
	)

	if err := fs.GetFile(
		*fte._sourceInputVideoFile,
		*fte._temporaryInputVideoFile,
	); err != nil {
		return fmt.Errorf(
			"%s %s",
			"we cannot make a copy of input video file to video transcoder temporary file",
			err.Error(),
		)
	}

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.prepare_input_video_file",
		"input source video file successfully copied to video transcoder working path",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"temporary_input_video_file": fte._temporaryInputVideoFile.FullPath(),
			"source_input_video_file":    fte._sourceInputVideoFile.FullPath(),
		},
	)

	return nil
}

// _transcodingInputVideoFile start transcoding input video file.
func (fte *FfmpegTranscoderEngine) _transcodingInputVideoFile() (err error) {

	var ffmpegHandlerFile interface{}

	// check if handler file exists in config
	if utilities.InSlice(fte._taskPayload.VideoTranscoder.Output.Handler,
		fte.Config.GetStringSlice("Components.VideoTranscoding.Engine.Ffmpeg.Handlers"),
	) {

		// check if handler file from default handlers
		if strings.HasPrefix(fte._taskPayload.VideoTranscoder.Output.Handler, "rasbora:") {
			selectedHandler := strings.Split(fte._taskPayload.VideoTranscoder.Output.Handler, "rasbora:")
			ffmpegHandlerFile, err = handlersFS.ReadFile(filepath.Join("handlers",
				selectedHandler[1],
			))
			if err != nil {
				return err
			}
		}

		// check if handler is custom file
		if strings.HasPrefix(fte._taskPayload.VideoTranscoder.Output.Handler, "custom:") {
			selectedHandler := strings.Split(fte._taskPayload.VideoTranscoder.Output.Handler, "custom:")
			ffmpegHandlerFile, err = os.ReadFile(selectedHandler[1])
			if err != nil {
				return err
			}
		}
	}

	// if there is no handlers should return error
	if ffmpegHandlerFile == nil {
		return fmt.Errorf(
			"%v [task_id=%v] [video_trancoder_worker_id=%v] [ffmpeg_handler=%v] ",
			"unknown rasbora ffmpeg handler",
			fte._taskPayload.ID,
			fte._videoTranscoderWorkerID,
			fte._taskPayload.VideoTranscoder.Output.Handler,
		)
	}

	// parse handler template should have django template style
	ffmpegHandler, err := pongo2.FromBytes(ffmpegHandlerFile.([]byte))
	if err != nil {
		return err
	}

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.transcoding_input_video_file",
		"configuring ffmpeg handler",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"ffmpeg_handler":             fte._taskPayload.VideoTranscoder.Output.Handler,
		},
	)

	// data will be replaced inside handler template
	handlerData := pongo2.Context{
		"ffmpeg":         fte.Config.GetString("Components.VideoTranscoding.Engine.Ffmpeg.Executable"),
		"input":          fte._temporaryInputVideoFile.FullPath(),
		"args":           fte._taskPayload.VideoTranscoder.Output.Args,
		"logfile":        fte._temporaryProcessingLogFile,
		"inputVideoInfo": fte._inputVideoInformation,
		"progressListener": fmt.Sprintf(
			"tcp:%v",
			fte.Config.GetString("Components.VideoTranscoding.Engine.Ffmpeg.ProgressListener"),
		),
	}

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.transcoding_input_video_file",
		"sending data to ffmpeg handler",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"ffmpeg_handler":             fte._taskPayload.VideoTranscoder.Output.Handler,
			"ffmpeg_handler_data":        handlerData,
		},
	)

	// prepare to execute ffmpeg handler
	ffmpegHandlerCmd, err := ffmpegHandler.Execute(handlerData)
	if err != nil {
		return err
	}
	ffmpegHandlerCmd = strings.Join(strings.Fields(strings.ReplaceAll(ffmpegHandlerCmd, "\n", " ")), " ")

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.transcoding_input_video_file",
		"ffmpeg handler is parsed and ready to execute",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"ffmpeg_handler":             fte._taskPayload.VideoTranscoder.Output.Handler,
			"ffmpeg_handler_cmd":         ffmpegHandlerCmd,
		},
	)

	// execute ffmpeg handler and start ffmpeg processing events listener server
	cmd := exec.Command(fte.Config.GetString("Components.VideoTranscoding.Shell"), "-c", ffmpegHandlerCmd)
	monitor := NewFfmpegProgressingMonitor(*fte)
	out, err := cmd.CombinedOutput()
	if err != nil {
		monitor.StopMonitoringFfmpegProgress()
		return errors.New(string(out))
	}
	monitor.StopMonitoringFfmpegProgress()

	fte.Logger.Success(
		"ffmpeg_transcoder_engine.transcoding_input_video_file",
		"ffmpeg transcended input video without any problems",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"ffmpeg_handler":             fte._taskPayload.VideoTranscoder.Output.Handler,
		},
	)

	return nil
}

// _readInputVideoInformation read input video information.
func (fte *FfmpegTranscoderEngine) _readInputVideoInformation() error {

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.read_input_video_information",
		"reading input video information",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)

	videoInputInformationData, err := ffprobe.ProbeURL(context.Background(), fte._temporaryInputVideoFile.FullPath())
	if err != nil {
		return err
	}
	fte._inputVideoInformation = videoInputInformationData

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.read_input_video_information",
		"source video input information is ready",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)

	return nil
}

// _moveTranscoderOutputVideos move temporary output videos to main file system.
// [IMPORTANT] If there is new file system add this method should be modified
// or this func should re-write to take new file system without change source code.
func (fte *FfmpegTranscoderEngine) _moveTranscoderOutputVideos() error {

	fte.Logger.Info(
		"ffmpeg_transcoder_engine.move_transcoder_output_videos",
		"moving temporary output videos to main file system",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)

	var transcoderOutputVideos string

	if fte.Config.GetString("Filesystem.Type") == data.LocalFileSystemType.String() {
		transcoderOutputVideos = fte.Config.GetString("Filesystem.LocalStorage.Folders.TranscoderOutputVideos")

		if err := os.MkdirAll(filepath.Dir(filepath.Join(transcoderOutputVideos, "rasbora.tmp")), os.ModePerm); err != nil {
			return err
		}

		fte.Logger.Debug(
			"ffmpeg_transcoder_engine.move_transcoder_output_videos",
			"moving output video files from temporary working folder",
			map[string]interface{}{
				"task_id":                    fte._queueable.ID,
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"move_to_filesystem":         data.LocalFileSystemType.String(),
			},
		)
	}

	if fte.Config.GetString("Filesystem.Type") == data.ObjectFileSystemType.String() {
		transcoderOutputVideos = fte.Config.GetString("Filesystem.ObjectStorage.Buckets.TranscoderOutputVideos")
		fte.Logger.Debug(
			"ffmpeg_transcoder_engine.move_transcoder_output_videos",
			"moving output video files from temporary working folder",
			map[string]interface{}{
				"task_id":                    fte._queueable.ID,
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"move_to_filesystem":         data.ObjectFileSystemType.String(),
			},
		)
	}

	var finalOutputVideoFiles []data.File

	for _, temporaryVideoOutputFile := range *fte._temporaryOutputVideoFiles {

		finalOutputVideoFile := data.File{
			FileMeta: temporaryVideoOutputFile.FileMeta,
			FileName: temporaryVideoOutputFile.FileName,
			FilePath: transcoderOutputVideos,
		}

		finalOutputVideoFiles = append(finalOutputVideoFiles, finalOutputVideoFile)

		if err := fte.FileSystem.PutFile(
			temporaryVideoOutputFile,
			finalOutputVideoFile,
		); err != nil {
			return err
		}

		fte.Logger.Debug(
			"ffmpeg_transcoder_engine.move_transcoder_output_videos",
			"moving output video file",
			map[string]interface{}{
				"task_id":                     fte._queueable.ID,
				"video_transcoder_worker_id":  fte._videoTranscoderWorkerID,
				"temporary_video_output_file": temporaryVideoOutputFile.FullPath(),
				"final_output_video_file":     finalOutputVideoFile.FullPath(),
			},
		)
	}

	fte._finalOutputVideoFiles = &finalOutputVideoFiles

	fte.Logger.Info(
		"ffmpeg_transcoder_engine.move_transcoder_output_videos",
		"all output videos moved to main file system",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)

	return nil
}

// _moveTranscoderProcessingLog move video transcoder processing logs main file system.
// [IMPORTANT] If there is new file system add this method should be modified
// or this func should re-write to take new file system without change source code.
func (fte *FfmpegTranscoderEngine) _moveTranscoderProcessingLog() error {
	fte.Logger.Info(
		"ffmpeg_transcoder_engine.move_transcoder_processing_log",
		"move video transcoder processing logs main file system",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)

	var transcoderProcessingLogsPath string

	if fte.Config.GetString("Filesystem.Type") == data.LocalFileSystemType.String() {
		transcoderProcessingLogsPath = fte.Config.GetString("Filesystem.LocalStorage.Folders.TranscoderProcessingLogs")

		if err := os.MkdirAll(filepath.Dir(filepath.Join(transcoderProcessingLogsPath, "rasbora.tmp")), os.ModePerm); err != nil {
			return err
		}

		fte.Logger.Debug(
			"ffmpeg_transcoder_engine.move_transcoder_processing_log",
			"select file system for processing log",
			map[string]interface{}{
				"task_id":                    fte._queueable.ID,
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"filesystem_type":            data.LocalFileSystemType.String(),
				"folder":                     transcoderProcessingLogsPath,
			},
		)
	}

	if fte.Config.GetString("Filesystem.Type") == data.ObjectFileSystemType.String() {
		transcoderProcessingLogsPath = fte.Config.GetString("Filesystem.ObjectStorage.Buckets.TranscoderProcessingLogs")
		fte.Logger.Debug(
			"ffmpeg_transcoder_engine.move_transcoder_processing_log",
			"select file system for processing log",
			map[string]interface{}{
				"task_id":                    fte._queueable.ID,
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"filesystem_type":            data.ObjectFileSystemType.String(),
				"bucket":                     transcoderProcessingLogsPath,
			},
		)
	}

	fte._finalProcessingLogFile = &data.File{
		FileMeta: map[string]interface{}{
			"task_id": fte._taskPayload.ID,
		},
		FileName: fte._temporaryProcessingLogFile.FileName,
		FilePath: transcoderProcessingLogsPath,
	}

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.move_transcoder_processing_log",
		"move video transcoder processing file",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"move_from":                  *fte._temporaryProcessingLogFile,
			"move_to":                    *fte._finalProcessingLogFile,
		},
	)

	return fte.FileSystem.PutFile(
		*fte._temporaryProcessingLogFile,
		*fte._finalProcessingLogFile,
	)
}

// _cleanAndPrepareForNextTask clean up after finish transcoding
func (fte *FfmpegTranscoderEngine) _cleanAndPrepareForNextTask() {

	// remove temporary input video file
	_ = os.RemoveAll(fte._temporaryInputVideoFile.FullPath())

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.clean_and_prepare_next_task",
		"remove temporary input video file",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"temporary_input_video_file": fte._temporaryInputVideoFile.FullPath(),
		},
	)

	// remove temporary video transcoder processing log file
	_ = os.RemoveAll(fte._temporaryProcessingLogFile.FullPath())

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.clean_and_prepare_next_task",
		"remove temporary input video file",
		map[string]interface{}{
			"task_id":                       fte._queueable.ID,
			"video_transcoder_worker_id":    fte._videoTranscoderWorkerID,
			"temporary_processing_log_file": fte._temporaryProcessingLogFile.FullPath(),
		},
	)

	// remove temporary video output files
	for _, temporaryOutputVideoFile := range *fte._temporaryOutputVideoFiles {
		_ = os.RemoveAll(temporaryOutputVideoFile.FullPath())
	}

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.clean_and_prepare_next_task",
		"remove temporary video output files",
		map[string]interface{}{
			"task_id":                      fte._queueable.ID,
			"video_transcoder_worker_id":   fte._videoTranscoderWorkerID,
			"temporary_video_output_files": fte._temporaryOutputVideoFiles,
		},
	)

	// remove temporary transcoding working space
	_ = os.RemoveAll(fte._temporaryWorkingPath)

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.clean_and_prepare_next_task",
		"remove temporary transcoding working space",
		map[string]interface{}{
			"task_id":                      fte._queueable.ID,
			"video_transcoder_worker_id":   fte._videoTranscoderWorkerID,
			"temporary_video_output_files": fte._temporaryWorkingPath,
		},
	)

	fte.Logger.Info(
		"ffmpeg_transcoder_engine.clean_and_prepare_next_task",
		"clean up working space is done",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)
}

// _createNewCallback create new callback
func (fte *FfmpegTranscoderEngine) _createNewCallback(err error) {

	var queueable = new(data.Queueable)
	var callback = new(data.Callback)

	callback.TaskId = fte._taskPayload.ID
	callback.Priority = fte._taskPayload.Priority
	callback.URL = fte._taskPayload.Callback.URL
	callback.Data = fte._taskPayload.Callback.Data

	if err == nil {
		callback.Error = false
		callback.Message = "video transcended without any problems"
		callback.VideoOutputFiles = *fte._finalOutputVideoFiles
		callback.ProcessingLogFile = *fte._finalProcessingLogFile
	} else {
		callback.Error = true
		callback.Message = err.Error()
	}

	callback.TaskTimeline.Add = fte._taskPayload.CreatedAt
	callback.TaskTimeline.Failed = fte._taskPayload.FailedAt
	callback.TaskTimeline.Finished = fte._taskPayload.FinishedAt
	callback.TaskTimeline.Started = fte._taskPayload.StartedAt

	queueable.ID = fte._queueable.ID
	queueable.Payload = fte._queueable.Priority
	queueable.Payload = callback

	_ = fte.Database.Enqueue(fte._callbackManagerQueue, *queueable)

	fte.Logger.Info(
		"ffmpeg_transcoder_engine.create_new_callback",
		"create new callback and send it to waiting queue",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.create_new_callback",
		"callback data",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"callback":                   callback,
		},
	)

}

// _failedTask inform queue about failed task
func (fte *FfmpegTranscoderEngine) _failedTask(err error) {

	// move video transcoder processing logs to main filesystem.
	if err := fte._moveTranscoderProcessingLog(); err != nil {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.failed_task",
			fmt.Sprintf("error when move video transcoder processing logs to main filesystem: %v", err),
			map[string]interface{}{
				"task_id":                    fte._queueable.ID,
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			},
		)
	}

	fte.Logger.Info(
		"ffmpeg_transcoder_engine.failed_task",
		"informing queue about failed task",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)

	customError := map[string]interface{}{
		"msg":   err.Error(),
		"debug": debug.Stack(),
	}

	jsonError, _ := json.Marshal(customError)

	fte._taskPayload.FailedAt = time.Now().UnixMilli()

	fte._queueable.Payload = fte._taskPayload

	//get retry config for video video transcoder
	retryCount := fte.Database.TotalRetry(fte._videoTranscoderQueue, *fte._queueable)
	retryLimit := fte.Config.GetInt("Components.VideoTranscoding.MakeAsFailedAfterRetry")

	//make it fail when arrive to retry limit
	if retryCount >= retryLimit {
		fte.Logger.Debug(
			"ffmpeg_transcoder_engine.failed_task",
			"failed to transcode task after too many retries",
			map[string]interface{}{
				"task_id":                    fte._queueable.ID,
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
				"transcoder_retry_count":     retryCount,
				"transcoder_max_retry":       retryLimit,
			},
		)
		fte._createNewCallback(err)
		_ = fte.Database.Failed(fte._videoTranscoderQueue, *fte._queueable, errors.New(string(jsonError)))
		fte._cleanAndPrepareForNextTask()
		return
	}

	fte.Logger.Debug(
		"ffmpeg_transcoder_engine.failed_task",
		"sending back task to waiting queue again to retry transcoding one more time",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			"transcoder_retry_count":     retryCount,
			"transcoder_max_retry":       retryLimit,
		},
	)

	_ = fte.Database.Enqueue(fte._videoTranscoderQueue, *fte._queueable)
	fte._cleanAndPrepareForNextTask()
}

// _successTask inform queue about success task
func (fte *FfmpegTranscoderEngine) _successTask() {

	// move video transcoder processing logs to main filesystem.
	if err := fte._moveTranscoderProcessingLog(); err != nil {
		fte.Logger.Error(
			"ffmpeg_transcoder_engine.success_task",
			fmt.Sprintf("error when move video transcoder processing logs to main filesystem: %v", err),
			map[string]interface{}{
				"task_id":                    fte._queueable.ID,
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			},
		)
	}

	// move output files to main filesystem.
	if err := fte._moveTranscoderOutputVideos(); err != nil {
		fte.Logger.Error("ffmpeg_transcoder_engine.success_task",
			fmt.Sprintf("error when move output files to main filesystem: %v", err.Error()),
			map[string]interface{}{
				"task_id":                    fte._queueable.ID,
				"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
			},
		)
		fte._failedTask(err)
		return
	}

	fte.Logger.Info(
		"ffmpeg_transcoder_engine.success_task",
		"inform queue about success task",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)

	fte._taskPayload.FinishedAt = time.Now().UnixMilli()

	fte._queueable.Payload = fte._taskPayload

	_ = fte.Database.Finished(fte._videoTranscoderQueue, *fte._queueable)

	fte.Logger.Success(
		"ffmpeg_transcoder_engine.success_task",
		"task finished processing without any problems",
		map[string]interface{}{
			"task_id":                    fte._queueable.ID,
			"video_transcoder_worker_id": fte._videoTranscoderWorkerID,
		},
	)

	fte._createNewCallback(nil)

	fte._cleanAndPrepareForNextTask()
}
