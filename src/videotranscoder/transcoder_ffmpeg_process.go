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
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

// FfmpegProgressingMonitor hold an instance.
type FfmpegProgressingMonitor struct {
	ffmpegProcessingListener net.Listener
	stopSignal               chan interface{}
	waitGroup                sync.WaitGroup
	engine                   *FfmpegTranscoderEngine
}

// NewFfmpegProgressingMonitor start listen for ffmpeg processing events.
func NewFfmpegProgressingMonitor(ffmpegTranscoderEngine FfmpegTranscoderEngine) *FfmpegProgressingMonitor {

	ffmpegTranscoderEngine.Logger.Debug(
		"ffmpeg_progressing_monitor",
		"started",
		map[string]interface{}{
			"task_id":                    ffmpegTranscoderEngine._queueable.ID,
			"video_transcoder_worker_id": ffmpegTranscoderEngine._videoTranscoderWorkerID,
			"ffmpeg_handler":             ffmpegTranscoderEngine._taskPayload.VideoTranscoder.Output.Handler,
		},
	)

	fpm := &FfmpegProgressingMonitor{
		stopSignal: make(chan interface{}),
		engine:     &ffmpegTranscoderEngine,
	}

	addr := fpm.engine.Config.GetString("Components.VideoTranscoding.Engine.Ffmpeg.ProgressListener")

	ffmpegProcessingListener, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	fpm.ffmpegProcessingListener = ffmpegProcessingListener

	fpm.waitGroup.Add(1)

	go fpm._listenForFfmpegConnections()

	return fpm
}

// _listenForFfmpegConnections listen for ffmpeg engine connection.
func (fpm *FfmpegProgressingMonitor) _listenForFfmpegConnections() {
	defer fpm.waitGroup.Done()

	for {
		ffmpegEngineClient, err := fpm.ffmpegProcessingListener.Accept()
		if err != nil {
			select {
			case <-fpm.stopSignal:
				return
			default:
				return
			}
		} else {
			fpm.waitGroup.Add(1)
			go func() {
				fpm._handleFfmpegEvents(ffmpegEngineClient)
				fpm.waitGroup.Done()
			}()
		}
	}
}

// StopMonitoringFfmpegProgress stops the listener and waits for all goroutines to finish.
func (fpm *FfmpegProgressingMonitor) StopMonitoringFfmpegProgress() {
	close(fpm.stopSignal)

	_ = fpm.ffmpegProcessingListener.Close()

	fpm.waitGroup.Wait()

	fpm.engine.Logger.Debug(
		"ffmpeg_progressing_monitor",
		"stopped",
		map[string]interface{}{
			"task_id":                    fpm.engine._queueable.ID,
			"video_transcoder_worker_id": fpm.engine._videoTranscoderWorkerID,
			"ffmpeg_handler":             fpm.engine._taskPayload.VideoTranscoder.Output.Handler,
		},
	)
}

// _handleFfmpegEvents calculates processing progress, and send real-time updates with the processing status.
func (fpm *FfmpegProgressingMonitor) _handleFfmpegEvents(ffmpegProcessingConnection net.Conn) {

	fpm.engine.Logger.Debug(
		"ffmpeg_progressing_monitor",
		fmt.Sprintf("new connection: %v", ffmpegProcessingConnection.LocalAddr().String()),
		map[string]interface{}{
			"task_id":                    fpm.engine._queueable.ID,
			"video_transcoder_worker_id": fpm.engine._videoTranscoderWorkerID,
			"ffmpeg_handler":             fpm.engine._taskPayload.VideoTranscoder.Output.Handler,
		},
	)

	//handle close ffmpeg processing connection
	defer func(ffmpegProcessingConnection net.Conn) {
		err := ffmpegProcessingConnection.Close()
		if err != nil {
			fpm.engine.Logger.Error(
				"ffmpeg_progressing_monitor",
				fmt.Sprintf("error when closing connection: %v", err.Error()),
				map[string]interface{}{
					"task_id":                    fpm.engine._queueable.ID,
					"video_transcoder_worker_id": fpm.engine._videoTranscoderWorkerID,
					"ffmpeg_handler":             fpm.engine._taskPayload.VideoTranscoder.Output.Handler,
				},
			)
		}
	}(ffmpegProcessingConnection)

	buf := make([]byte, 2048)

	data := map[string]interface{}{}

	for {
		n, err := ffmpegProcessingConnection.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			return
		}

		lines := strings.Split(string(buf[:n]), "\n")

		for _, line := range lines {

			if strings.Contains(line, "fps=") {
				data["fps"] = strings.Split(line, "=")[1]
			}

			if strings.Contains(line, "speed=") {
				data["speed"] = strings.Split(line, "=")[1]
			}

			if strings.Contains(line, "frame=") {
				data["frame"] = strings.Split(line, "=")[1]
			}

			if strings.Contains(line, "bitrate=") {
				data["bitrate"] = strings.Split(line, "=")[1]
			}

			if strings.Contains(line, "out_time_ms=") {
				data["time"] = strings.Split(line, "=")[1]
			}

			if len(data) < 5 {
				continue
			}

		}

		processedTime, err := strconv.ParseFloat(data["time"].(string), 64)
		if err != nil {
			return
		}

		duration := fpm.engine._inputVideoInformation.Format.Duration()

		data["task_id"] = fpm.engine._queueable.ID

		data["duration"] = duration.Microseconds()

		data["percentage"] = fmt.Sprintf("%.2f", processedTime/float64(duration.Microseconds())*100)

		if err := fpm.engine.Database.Processing(fpm.engine._videoTranscoderQueue, data); err != nil {
			fpm.engine.Logger.Error(
				"ffmpeg_progressing_monitor",
				fmt.Sprintf("error when send processing event to stream: %v", err.Error()),
				map[string]interface{}{
					"task_id":                    fpm.engine._queueable.ID,
					"video_transcoder_worker_id": fpm.engine._videoTranscoderWorkerID,
					"ffmpeg_handler":             fpm.engine._taskPayload.VideoTranscoder.Output.Handler,
				},
			)
		}

		fpm.engine.Logger.Debug(
			"ffmpeg_progressing_monitor",
			fmt.Sprintf(
				"[task_id=%v] [trancoder_worker_id=%v] [ffmpeg_handler=%v] %v",
				fpm.engine._queueable.ID,
				fpm.engine._videoTranscoderWorkerID,
				fpm.engine._taskPayload.VideoTranscoder.Output.Handler,
				data,
			),
			map[string]interface{}{
				"task_id":                    fpm.engine._queueable.ID,
				"video_transcoder_worker_id": fpm.engine._videoTranscoderWorkerID,
				"ffmpeg_handler":             fpm.engine._taskPayload.VideoTranscoder.Output.Handler,
				"data":                       data,
			},
		)
	}
}
