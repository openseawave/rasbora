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

package heartbeat

import (
	"context"
	"fmt"
	"time"

	"openseawaves.com/rasbora/internal/config"
	"openseawaves.com/rasbora/internal/database"
	"openseawaves.com/rasbora/internal/logger"
)

// Name used as identifier.
const Name = "Heartbeat"

// Heartbeat hold an instance.
type Heartbeat struct {
	Config     *config.Config
	Logger     *logger.Logger
	Database   *database.Database
	WorkerId   string
	WorkerType string
}

// Start sending heart beats to update cluster status.
func (hb *Heartbeat) Start(ctx context.Context) {

	// get send heartbeat interval in seconds
	heartbeatSendInterval := hb.Config.GetInt("Components.Heartbeat.SendInterval")

	hb.Logger.Info(
		"heartbeat",
		"new heartbeat has been started to update cluster status",
		map[string]interface{}{
			"worker_id":   hb.WorkerId,
			"worker_type": hb.WorkerType,
		},
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:

			if err := hb.Database.SendHeartbeat(hb.WorkerId, hb.WorkerType); err != nil {
				hb.Logger.Warn(
					"heartbeat",
					fmt.Sprintf("cannot send heartbeat: %v", err.Error()),
					map[string]interface{}{
						"worker_id":   hb.WorkerId,
						"worker_type": hb.WorkerType,
					},
				)
			}

			time.Sleep(time.Duration(heartbeatSendInterval) * time.Second)

			hb.Logger.Debug(
				"heartbeat",
				"new heartbeat pulse has been sent",
				map[string]interface{}{
					"worker_id":   hb.WorkerId,
					"worker_type": hb.WorkerType,
				},
			)

		}
	}
}
