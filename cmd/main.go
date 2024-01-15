// Copyright (c) 2022-2023 OpenSeaWaves.com/Rasbora
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

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"openseawaves.com/rasbora/internal/config"
	"openseawaves.com/rasbora/internal/data"
	"openseawaves.com/rasbora/internal/database"
	"openseawaves.com/rasbora/internal/filesystem"
	"openseawaves.com/rasbora/internal/logger"
	"openseawaves.com/rasbora/internal/utilities"
	"openseawaves.com/rasbora/src/callbacks"
	"openseawaves.com/rasbora/src/heartbeat"
	"openseawaves.com/rasbora/src/systemradar"
	"openseawaves.com/rasbora/src/taskmanager"
	"openseawaves.com/rasbora/src/videotranscoder"
)

var (
	// It is used to wait for all active components to finish execution before the program exits.
	wg sync.WaitGroup
	// It provides a centralized mechanism for logging messages with various log levels.
	log *logger.Logger
	// It holds the application's configuration settings
	cfg *config.Config
	// It represents the connection to the chosen database backend and provides methods for interacting with it.
	db *database.Database
	// It holds the file system instance.
	fs *filesystem.FileSystem
)

// Its holds rasbora info.
var (
	Version   = "undefined"
	BuildTime = "undefined"
	GitHash   = "undefined"
	ctx       = context.Background()
)

// Initializing Rasbora.
func init() {
	//step 1 show info
	initShowRasboraInfo()
	//step 2 load config
	initInternalConfigManager()
	//step 3 open database connection
	initInternalDatabaseConnection()
	//step 4 load file system
	initInternalFileSystem()
	//step 5 start rasbora logger
	initInternalSystemLogger()
}

// main starting point for Rasbora.
func main() {
	activeComponents := cfg.GetStringSlice("Components.Active")

	log.Debug("main", "preloaded components", map[string]interface{}{
		"active_components": activeComponents,
	})

	for _, component := range activeComponents {

		//start callback manager component
		if component == callbacks.Name {
			_startComponent(
				cfg.GetString("Components.CallbackManager.UniqueID"),
				callbacks.Name,
				initComponentCallbackManager,
			)
		}

		//start task manager component
		if component == taskmanager.Name {
			_startComponent(
				cfg.GetString("Components.TaskManagement.UniqueID"),
				taskmanager.Name,
				initComponentTaskManager,
			)
		}

		//start video transcoder component
		if component == videotranscoder.Name {
			_startComponent(
				cfg.GetString("Components.VideoTranscoding.UniqueID"),
				videotranscoder.Name,
				initComponentVideoTranscoder,
			)
		}

		//start system radar component
		if component == systemradar.Name {
			_startComponent(
				cfg.GetString("Components.SystemRadar.UniqueID"),
				systemradar.Name,
				initComponentSystemRadar,
			)
		}
	}

	wg.Wait()
}

// showRasboraInfo show rasbora info
func initShowRasboraInfo() {
	fmt.Println("\nRasbora Video Transcoding")
	fmt.Printf("Version    : %s\n", Version)
	fmt.Printf("GitHash    : %s\n", GitHash)
	fmt.Printf("BuildTime  : %s\n", BuildTime)
	fmt.Printf("Copyright  : %s\n", "2022-2023 https://openseawaves.com/rasbora")
	fmt.Printf("License    : %s\n", "GNU Affero General Public License\n")
}

// initInternalSystemLogger initializes the system logger based on configuration settings.
func initInternalSystemLogger() {
	var outputs []data.LoggerOutputConfig

	for _, outputType := range cfg.GetStringSlice("Logger.Output") {

		if outputType == data.StdoutLoggerOutputType.String() {
			outputs = append(outputs, data.LoggerOutputConfig{
				OutputType:   data.StdoutLoggerOutputType,
				OutputWriter: os.Stdout,
			})
		}

		if outputType == data.DatabaseLoggerOutputType.String() {
			outputs = append(outputs, data.LoggerOutputConfig{
				OutputType: data.DatabaseLoggerOutputType,
				OutputWriter: database.LoggerOutput{
					Config:   cfg,
					Database: db,
				},
			})
		}

		if outputType == data.FileLoggerOutputType.String() {
			logFilePath := cfg.GetString("Filesystem.LocalStorage.Folders.LoggerFilePath")
			logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				fmt.Printf("error while opening log file: %v\n", err.Error())
				os.Exit(1)
			}
			outputs = append(outputs, data.LoggerOutputConfig{
				OutputType:   data.FileLoggerOutputType,
				OutputWriter: logFile,
			})
		}
	}

	log = logger.NewWithConfig(logger.Options{
		Output: outputs,
		Level:  cfg.GetIntSlice("Logger.Level"),
	})

	log.Success(
		"main.init.logger",
		"initialized successfully.",
		nil,
	)
}

// initInternalConfigManager load and parse config using viper.
func initInternalConfigManager() {
	l := logger.NewWithConfig(logger.Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.StdoutLoggerOutputType,
				OutputWriter: os.Stdout,
			},
		},
		Level: []int{1, 2, 3, 4},
	})

	l.Info(
		"main.init.config",
		"initializing",
		nil,
	)

	_viper := viper.New()

	// Configuration settings for Viper
	_viper.SetConfigName("config")
	_viper.SetConfigType("yaml")
	_viper.SetEnvPrefix("RASBORA")
	_viper.AutomaticEnv()
	_viper.AddConfigPath("/etc/rasbora/")
	_viper.AddConfigPath("./")
	_viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	// Attempt to read the configuration file
	if err := _viper.ReadInConfig(); err != nil {
		l.Error("main.init.config", err.Error(), nil)
		os.Exit(1)
	}

	// Create a new configuration manager using Viper
	cfg = config.New(&config.ViperConfigManager{Viper: _viper})

	l.Success(
		"main.init.config",
		"initialized successfully.",
		nil,
	)
}

// initInternalDatabaseConnection prepare the database connection based on configuration settings.
func initInternalDatabaseConnection() {

	l := logger.NewWithConfig(logger.Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.StdoutLoggerOutputType,
				OutputWriter: os.Stdout,
			},
		},
		Level: []int{1, 2, 3, 4},
	})

	dbType := cfg.GetString("Database.Type")

	l.Info(
		"main.init.database",
		"initializing database",
		map[string]interface{}{
			"database_type": dbType,
		},
	)

	if dbType == "Redis" {
		l.Info(
			"main.init.database",
			"testing the connection to database",
			map[string]interface{}{
				"database_type": dbType,
			},
		)

		rds := redis.NewClient(&redis.Options{
			Addr:     cfg.GetString("Database.Redis.Connection.Address"),
			Password: cfg.GetString("Database.Redis.Connection.Password"),
			DB:       cfg.GetInt("Database.Redis.Connection.DatabaseIndex"),
		})

		if _, err := rds.Ping(context.Background()).Result(); err != nil {
			l.Error(
				"main.init.database",
				fmt.Sprintf("error when open connection to database: %s", err.Error()),
				map[string]interface{}{
					"database_type": dbType,
				},
			)
			os.Exit(1)
		}

		db = database.New(&database.RedisDatabaseManager{
			Redis:  rds,
			Config: cfg,
			Logger: log,
		})

		l.Success(
			"main.init.database",
			"database successfully has been started",
			map[string]interface{}{
				"database_type": dbType,
			},
		)
		return
	}

	l.Error(
		"main.init.database",
		"There is no database founded in config file",
		nil,
	)

	os.Exit(1)
}

// initInternalFileSystem prepare file system based on configuration settings.
func initInternalFileSystem() {

	l := logger.NewWithConfig(logger.Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.StdoutLoggerOutputType,
				OutputWriter: os.Stdout,
			},
		},
		Level: []int{1, 2, 3, 4},
	})

	fileSystemType := cfg.GetString("Filesystem.Type")

	l.Info(
		"main.init.file_system",
		"initializing filesystem",
		map[string]interface{}{
			"filesystem_type": fileSystemType,
		},
	)

	if fileSystemType == data.ObjectFileSystemType.String() {
		l.Info(
			"main.init.filesystem",
			"testing the connection to filesystem",
			map[string]interface{}{
				"filesystem_type": fileSystemType,
			},
		)

		minioClient, err := filesystem.NewObjectClient(*cfg)

		if err != nil {
			l.Error(
				"main.init.filesystem",
				fmt.Sprintf("error when open connection to filesystem: %s", err.Error()),
				map[string]interface{}{
					"filesystem_type": fileSystemType,
				},
			)
			os.Exit(1)
		}

		fs = filesystem.NewFileSystem(&filesystem.ObjectFileSystem{
			Minio: minioClient,
		})

		l.Success(
			"main.init.filesystem",
			"filesystem successfully has been started",
			map[string]interface{}{
				"filesystem_type": fileSystemType,
			},
		)

		return
	}

	if fileSystemType == data.LocalFileSystemType.String() {

		fs = filesystem.NewFileSystem(&filesystem.LocalFileSystem{})

		l.Success(
			"main.init.filesystem",
			"filesystem successfully has been started",
			map[string]interface{}{
				"filesystem_type": fileSystemType,
			},
		)

		return
	}

	l.Error(
		"main.init.filesystem",
		"There is no filesystem founded in config file",
		nil,
	)

	os.Exit(1)
}

// initInternalTaskManagerComponent  this func used to start task manager component.
func initComponentTaskManager() {
	log.Info(
		"main.init.task_manager_component",
		"initializing",
		nil,
	)

	if cfg.GetString("Components.TaskManagement.Active") == "Restful" {

		log.Info(
			"main.init.task_manager_component",
			"default protocol has been selected for task manager",
			map[string]interface{}{
				"protocol_type": "Restful",
			},
		)

		taskmanager.New(&taskmanager.RestfulTaskManager{
			Config:   cfg,
			Logger:   log,
			Database: db,
		}).StartTaskManager(ctx)

		log.Success(
			"main.init.task_manager_component",
			"has been started successfully",
			map[string]interface{}{
				"protocol_type": "Restful",
			},
		)

		return
	}

	log.Error(
		"main.init.task_manager_component",
		"There is no protocol founded in config file",
		nil,
	)

	os.Exit(1)
}

// initVideoTranscoderComponent this func used to start video transcoder component.
func initComponentVideoTranscoder() {
	log.Info(
		"main.init.video_transcoder_component",
		"initializing",
		nil,
	)

	if cfg.GetString("Components.VideoTranscoding.Engine.Type") == "ffmpeg" {

		log.Info(
			"main.init.video_transcoder_component",
			"configuring transcoder engine",
			map[string]interface{}{
				"transcoder_engine_type": "ffmpeg",
			},
		)

		videotranscoder.New(&videotranscoder.FfmpegTranscoderEngine{
			Config:     cfg,
			Logger:     log,
			Database:   db,
			FileSystem: fs,
		}).StarTranscoderEngine(ctx)

		log.Success(
			"main.init.video_transcoder_component",
			"transcoder engine started successfully",
			map[string]interface{}{
				"transcoder_engine_type": "ffmpeg",
			},
		)

		return
	}

	log.Error(
		"main.init.video_transcoder_component",
		"there is no engine founded in config file",
		nil,
	)

	os.Exit(1)
}

// initCallbackManagerComponent this func used to start callback manager.
func initComponentCallbackManager() {
	log.Info(
		"main.init.callback_manager_component",
		"initializing",
		nil,
	)

	if cfg.GetString("Components.CallbackManager.Active") == "http" {

		log.Info(
			"main.init.callback_manager_component",
			"configuring callback manager protocol",
			map[string]interface{}{
				"callback_protocol_type": "http",
			},
		)

		callbacks.New(&callbacks.HttpCallbackManager{
			Config:   cfg,
			Logger:   log,
			Database: db,
		}).StartCallbackManager(ctx)

		log.Success(
			"main.init.callback_manager_component",
			"has been successfully started",
			map[string]interface{}{
				"callback_protocol_type": "http",
			},
		)
	}

	log.Error(
		"main.init.callback_manager_component",
		"There is no protocol founded in config file",
		nil,
	)

	os.Exit(1)
}

// initSystemRadarComponent this func used to start system radar.
func initComponentSystemRadar() {
	log.Info(
		"main.init.system_radar_component",
		"initializing",
		nil,
	)

	systemradar.NewRadar(cfg, log, db).StartRadar(ctx)

	log.Success(
		"main.init.system_radar_component",
		"started successfully",
		nil,
	)
}

// _monitorComponent sending heartbeat single about component.
func _monitorComponent(workerId, workerType string) {

	log.Debug(
		"main.monitor_component",
		"start send heartbeat",
		map[string]interface{}{
			"worker_type": workerType,
			"worker_id":   workerId,
		},
	)

	updateClusterStatus := &heartbeat.Heartbeat{
		Config:     cfg,
		Database:   db,
		Logger:     log,
		WorkerId:   workerId,
		WorkerType: workerType,
	}

	updateClusterStatus.Start(ctx)
}

// _startComponent this func used to load and start components.
func _startComponent(workerId, workerType string, loader func()) {
	activeComponents := cfg.GetStringSlice("Components.Active")

	wg.Add(1)
	go func() {
		loader()
		defer wg.Done()
	}()

	if utilities.InSlice(heartbeat.Name, activeComponents) {
		wg.Add(1)
		go _monitorComponent(
			workerId,
			workerType,
		)
	}
}
