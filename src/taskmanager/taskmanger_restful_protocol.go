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

package taskmanager

import (
	"context"
	"os"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"openseawave.com/rasbora/internal/config"
	"openseawave.com/rasbora/internal/data"
	"openseawave.com/rasbora/internal/database"
	"openseawave.com/rasbora/internal/logger"

	// Auto-generated swagger documentation
	_ "openseawave.com/rasbora/src/taskmanager/docs"
)

// RestfulTaskManager hold an instance
type RestfulTaskManager struct {
	Config                *config.Config
	Logger                *logger.Logger
	Database              *database.Database
	_videoTranscoderQueue string
	_taskManagerWorkerID  string
	app                   *fiber.App
}

// @title Rasbora Task Manager API
// @version 1.0
// @description Task Manager API for Rasbora Distributed Video Transcoding.
// @contact.name Rasbora
// @contact.url	https://rasbora.openseawave.com
// @contact.email rasbora.support@openseawave.com
// @license.name GNU Affero General Public License
// @license.url	http://www.gnu.org/licenses/
// @host localhost:3701
// @BasePath /v1.0
// @schemes http
func (rtm *RestfulTaskManager) StartTaskManager(ctx context.Context) {

	// get video transcoder queue name
	rtm._videoTranscoderQueue = rtm.Config.GetString("Components.VideoTranscoding.Queue")

	// get task manager worker id
	rtm._taskManagerWorkerID = rtm.Config.GetString("Components.TaskManagement.UniqueID")

	// get listen addr
	listenAddress := rtm.Config.GetString("Components.TaskManagement.Protocols.Restful.ListenAddress")

	rtm.Logger.Debug(
		"restful_task_manager",
		"starting",
		map[string]interface{}{
			"task_manager_worker_id": rtm._taskManagerWorkerID,
			"protocol":               "restful",
			"address":                listenAddress,
		},
	)

	rtm.app = fiber.New(fiber.Config{
		AppName:               "Rasbora",
		DisableStartupMessage: true,
		ReduceMemoryUsage:     true,
	})

	rtm._prepareHttpServer()

	select {
	case <-ctx.Done():
		return
	default:
		if err := rtm.app.Listen(listenAddress); err != nil {
			rtm.Logger.Error(
				"restful_task_manager.Run",
				err.Error(),
				map[string]interface{}{
					"task_manager_worker_id": rtm._taskManagerWorkerID,
					"protocol":               "restful",
					"address":                listenAddress,
				},
			)
			os.Exit(1)
		}
	}
}

// _prepareHttpServer return errors as json response.
func (rtm *RestfulTaskManager) _prepareHttpServer() {
	// Load recover middleware.
	rtm.app.Use(recover.New())

	// Load json errors middleware.
	rtm.app.Use(rtm._middlewareJsonErrors)

	// Add endpoint to server creating new tasks.
	rtm.app.Post("/v1.0/tasks/create", rtm._endpointCreateNewTask)

	// Add endpoint to serve swagger documentation.
	rtm.app.Get("/swagger/*", swagger.HandlerDefault)
}

// _middlewareJsonErrors return errors as json response.
func (rtm *RestfulTaskManager) _middlewareJsonErrors(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {

		rtm.Logger.Error(
			"restful_task_manager.Run",
			err.Error(),
			map[string]interface{}{
				"task_manager_worker_id": rtm._taskManagerWorkerID,
				"protocol":               "restful",
			},
		)

		return c.Status(fiber.StatusInternalServerError).JSON(data.Response{
			Error:   true,
			Message: err.Error(),
			Payload: nil,
		})

	}

	return nil
}

// CreateTask godoc
// @Summary Create new task for video transcoding.
// @Description Create new task for video transcoding.
// @Tags tasks
// @Param task body data.Task true "Task data"
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} data.Response
// @Failure 400 {object} data.Response
// @Failure 404 {object} data.Response
// @Failure 500 {object} data.Response
// @Router /tasks/create [post]
func (rtm *RestfulTaskManager) _endpointCreateNewTask(c *fiber.Ctx) error {

	rtm.Logger.Info(
		"restful_task_manager.create_new_task",
		"received create new task request",
		map[string]interface{}{
			"task_manager_worker_id": rtm._taskManagerWorkerID,
		},
	)

	var task data.Task
	var queueable data.Queueable

	if err := c.BodyParser(&task); err != nil {
		rtm.Logger.Error(
			"restful_task_manager.create_new_task",
			"error reading request body",
			map[string]interface{}{
				"task_manager_worker_id": rtm._taskManagerWorkerID,
			},
		)
		return err
	}

	rtm.Logger.Debug(
		"restful_task_manager.create_new_task",
		"task data",
		map[string]interface{}{
			"task_manager_worker_id": rtm._taskManagerWorkerID,
			"task_data":              task,
		},
	)

	check := validator.New()
	if err := check.Struct(task); err != nil {
		rtm.Logger.Error(
			"restful_task_manager.create_new_task",
			"error json input is not correct",
			map[string]interface{}{
				"task_manager_worker_id": rtm._taskManagerWorkerID,
				"task_id":                task.ID,
			},
		)
		return err
	}

	if len(task.ID) <= 0 {
		task.ID = uuid.NewString()
	}

	task.CreatedAt = time.Now().UnixMilli()

	queueable.ID = task.ID
	queueable.Priority = *task.Priority
	queueable.Payload = task

	if err := rtm.Database.Enqueue(rtm._videoTranscoderQueue, queueable); err != nil {
		rtm.Logger.Error(
			"restful_task_manager.create_new_task",
			"error when saving task in database",
			map[string]interface{}{
				"task_manager_worker_id": rtm._taskManagerWorkerID,
				"task_id":                task.ID,
			},
		)
		return err
	}

	rtm.Logger.Success(
		"restful_task_manager.create_new_task",
		"task created without any problems",
		map[string]interface{}{
			"task_manager_worker_id": rtm._taskManagerWorkerID,
			"task_id":                task.ID,
		},
	)

	_ = c.JSON(data.Response{Error: false, Message: "task created without any problems",
		Payload: struct {
			TaskId string `json:"task_id"`
		}{
			TaskId: task.ID,
		},
	})

	return nil
}
