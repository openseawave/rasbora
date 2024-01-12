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

package callbacks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"openseawaves.com/rasbora/internal/config"
	"openseawaves.com/rasbora/internal/data"
	"openseawaves.com/rasbora/internal/database"
	"openseawaves.com/rasbora/internal/logger"
)

// HttpCallbackManager use http to send callbacks.
type HttpCallbackManager struct {
	Config           *config.Config
	Logger           *logger.Logger
	Database         *database.Database
	workerId         string
	_queueName       string
	_queueable       *data.Queueable
	_callbackPayload *data.Callback
}

// StartCallbackManager start callback worker who listen for new callbacks
func (hcm *HttpCallbackManager) StartCallbackManager(ctx context.Context) {

	// get callback manager queue name
	hcm._queueName = hcm.Config.GetString("Components.CallbackManager.Queue")

	// get callback manager worker id
	hcm.workerId = hcm.Config.GetString("Components.CallbackManager.UniqueID")

	// get callback manager time interval used to check for new callback.
	checkNewCallbackInterval := hcm.Config.GetInt("Components.CallbackManager.CheckNewCallbackInterval")

	hcm.Logger.Info(
		"http_callback_manager",
		"try to init callback manager worker",
		map[string]interface{}{
			"callback_worker_id": hcm.workerId,
		},
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(checkNewCallbackInterval) * time.Second)

			callback, err := hcm.Database.Dequeue(hcm._queueName, hcm.workerId)

			if err != nil {
				hcm.Logger.Debug(
					"http_callback_manager",
					"callbacks queue is empty",
					map[string]interface{}{
						"callback_worker_id": hcm.workerId,
					},
				)
				continue
			}

			hcm.Logger.Success(
				"http_callback_manager",
				"preparing new callback to send",
				map[string]interface{}{
					"callback_id":        callback.ID,
					"callback_worker_id": hcm.workerId,
				},
			)

			hcm._send(callback)
		}
	}
}

// _send trying to send callback
func (hcm *HttpCallbackManager) _send(item data.Queueable) {
	hcm._queueable = &item

	//recover from panic
	defer func() {
		if r := recover(); r != nil {
			jsonData, _ := json.Marshal(r)
			hcm._failed(errors.New(string(jsonData)))
			hcm.Logger.Error(
				"http_callback_manager.send",
				fmt.Sprintf("we got panic: %v", string(jsonData)),
				map[string]interface{}{
					"callback_id":        hcm._queueable.ID,
					"callback_worker_id": hcm.workerId,
				},
			)
		}
	}()

	// get callback payload from queue item.
	callbackAsJsonBytes, errJ := json.Marshal(item.Payload)
	if errJ != nil {
		hcm.Logger.Error(
			"http_callback_manager.send",
			fmt.Sprintf("cannot cast queueable to json string: %v", errJ.Error()),
			map[string]interface{}{
				"callback_id":        hcm._queueable.ID,
				"callback_worker_id": hcm.workerId,
			},
		)
		hcm._failed(errJ)
		return
	}

	// set callback payload
	errU := json.Unmarshal(callbackAsJsonBytes, &hcm._callbackPayload)
	if errU != nil {
		hcm.Logger.Error(
			"http_callback_manager.send",
			fmt.Sprintf("cannot cast queueable to callback struct: %v", errJ.Error()),
			map[string]interface{}{
				"callback_id":        hcm._queueable.ID,
				"callback_worker_id": hcm.workerId,
			},
		)
		hcm._failed(errU)
		return
	}

	hcm.Logger.Debug(
		"http_callback_manager.send",
		"preparing callback",
		map[string]interface{}{
			"callback_id":        hcm._queueable.ID,
			"callback_worker_id": hcm.workerId,
			"callback_data":      &hcm._callbackPayload,
		},
	)

	//generate post request with callback endpoint and payload
	req, err := http.NewRequest("POST", hcm._callbackPayload.URL, bytes.NewBuffer(callbackAsJsonBytes))
	if err != nil {
		hcm.Logger.Error(
			"http_callback_manager.send",
			fmt.Sprintf("cannot create http request: %v", err.Error()),
			map[string]interface{}{
				"callback_id":        hcm._queueable.ID,
				"callback_worker_id": hcm.workerId,
			},
		)
		hcm._failed(errU)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	//trying to send callback throw http post request
	sendTimeout := hcm.Config.GetInt("Components.CallbackManager.Http.SendingTimeout")
	client := &http.Client{
		Timeout: time.Duration(sendTimeout) * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		hcm.Logger.Error(
			"http_callback_manager.send",
			fmt.Sprintf("callback cannot send receiver unreachable: %v", err.Error()),
			map[string]interface{}{
				"callback_id":        hcm._queueable.ID,
				"callback_worker_id": hcm.workerId,
			},
		)
		hcm._failed(err)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	//read response from receiver
	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		hcm.Logger.Error(
			"http_callback_manager.send",
			fmt.Sprintf("cannot read receiver response body: %v", err.Error()),
			map[string]interface{}{
				"callback_id":        hcm._queueable.ID,
				"callback_worker_id": hcm.workerId,
			},
		)
		hcm._failed(err)
		return
	}

	hcm.Logger.Debug(
		"http_callback_manager.send",
		"receiver response with data",
		map[string]interface{}{
			"callback_id":        hcm._queueable.ID,
			"callback_worker_id": hcm.workerId,
			"callback_response":  string(bodyData),
		},
	)

	//receiver handle the callback without issues
	if response.StatusCode == 200 {
		hcm.Logger.Success(
			"http_callback_manager.send",
			"callback successfully sent",
			map[string]interface{}{
				"callback_id":        hcm._queueable.ID,
				"callback_worker_id": hcm.workerId,
			},
		)
		hcm._success()
		return
	}

	//receiver can not handle the callback
	hcm.Logger.Error(
		"http_callback_manager.send",
		"callback receiver response with error",
		map[string]interface{}{
			"callback_id":                   hcm._queueable.ID,
			"callback_worker_id":            hcm.workerId,
			"callback_response":             string(bodyData),
			"callback_response_status_code": response.StatusCode,
		},
	)

	//callback fail to send or receive
	hcm._failed(errors.New(string(bodyData)))
}

// _failed handle failed callback
func (hcm *HttpCallbackManager) _failed(err error) {

	//get retry config for callback manager
	retryCount := hcm.Database.TotalRetry(hcm._queueName, *hcm._queueable)
	retryLimit := hcm.Config.GetInt("Components.CallbackManager.MakeAsFailedAfterRetry")

	//make it fail when arrive to retry limit
	if retryCount >= retryLimit {
		hcm.Logger.Debug(
			"http_callback_manager.failed",
			"failed to send callback after too many retries",
			map[string]interface{}{
				"callback_id":          hcm._queueable.ID,
				"callback_worker_id":   hcm.workerId,
				"callback_retry_count": retryCount,
				"callback_max_retry":   retryLimit,
			},
		)
		_ = hcm.Database.Failed(hcm._queueName, *hcm._queueable, err)
		return
	}

	hcm.Logger.Debug(
		"http_callback_manager.failed",
		"return callback to waiting queue again to retry send callback one more time",
		map[string]interface{}{
			"callback_id":          hcm._queueable.ID,
			"callback_worker_id":   hcm.workerId,
			"callback_retry_count": retryCount,
			"callback_max_retry":   retryLimit,
		},
	)
	_ = hcm.Database.Enqueue(hcm._queueName, *hcm._queueable)
}

// _success handle successful callback
func (hcm *HttpCallbackManager) _success() {
	_ = hcm.Database.Finished(hcm._queueName, *hcm._queueable)
}
