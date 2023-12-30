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

package systemradar

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"openseawaves.com/rasbora/internal/config"
	"openseawaves.com/rasbora/internal/data"
	"openseawaves.com/rasbora/internal/database"
	"openseawaves.com/rasbora/internal/logger"
)

// Name used as identifier.
const Name = "SystemRadar"

// SystemRadar holds an instance.
type SystemRadar struct {
	Config   *config.Config
	Logger   *logger.Logger
	Database *database.Database
	Feedback *data.Radar
	workerId string
}

// NewRadar make new system radar.
func NewRadar(cfg *config.Config, log *logger.Logger, db *database.Database) *SystemRadar {
	return &SystemRadar{
		Config:   cfg,
		Logger:   log,
		Database: db,
		Feedback: new(data.Radar),
	}
}

// StartRadar make background thread to scan system info.
func (sr *SystemRadar) StartRadar(ctx context.Context) {
	sr.Logger.Info(
		"system_radar",
		"initializing system radar worker",
		map[string]interface{}{
			"systemradar_worker_id": sr.workerId,
		},
	)

	sr.workerId = sr.Config.GetString("Components.SystemRadar.UniqueID")

	scanSystemInterval := sr.Config.GetInt("Components.SystemRadar.ScanInterval")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			sr._scan()
			time.Sleep(time.Duration(scanSystemInterval) * time.Second)
		}
	}
}

// _scan start searching for data
func (sr *SystemRadar) _scan() {

	//recover from panic
	defer func() {
		if r := recover(); r != nil {
			jsonData, _ := json.Marshal(r)
			sr.Logger.Error(
				"system_radar.scan",
				fmt.Sprintf("we got panic: %v", string(jsonData)),
				map[string]interface{}{
					"systemradar_worker_id": sr.workerId,
				},
			)
		}
	}()

	// search for system data.
	sr._find()
	// try to send radar data to real-time stream.
	sr._send()
}

// _find search for system data.
func (sr *SystemRadar) _find() {

	sr.Logger.Debug(
		"system_radar.find",
		"scanning system",
		map[string]interface{}{
			"systemradar_worker_id": sr.workerId,
		},
	)

	sr.Feedback.SystemRadarId = sr.workerId

	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	memoryInfoJson, _ := json.Marshal(memoryInfo)
	sr.Feedback.MemoryInfo = memoryInfoJson

	cpuInfo, err := cpu.Info()
	if err != nil {
		return
	}
	cpuInfoJson, _ := json.Marshal(cpuInfo)
	sr.Feedback.CpuInfo = cpuInfoJson

	percentageAll, err := cpu.Percent(time.Second, true)
	if err != nil {
		return
	}
	percentageAllJson, _ := json.Marshal(percentageAll)
	sr.Feedback.CpuUsageAll = percentageAllJson

	diskStat, err := disk.Usage(sr.Config.GetString("Components.SystemRadar.DiskStat"))
	if err != nil {
		return
	}
	diskStatJson, _ := json.Marshal(diskStat)
	sr.Feedback.DiskUsage = diskStatJson

	networkStat, err := net.IOCounters(false)
	if err != nil {
		return
	}
	networkStatJson, _ := json.Marshal(networkStat)
	sr.Feedback.NetworkStat = networkStatJson

	networkInterfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	networkInterfacesJson, _ := json.Marshal(networkInterfaces)
	sr.Feedback.NetworkInterfaces = networkInterfacesJson

	hostInfo, err := host.Info()
	if err != nil {
		return
	}
	hostInfoJson, _ := json.Marshal(hostInfo)
	sr.Feedback.HostInfo = hostInfoJson
}

// _send try to send radar data to real-time stream.
func (sr *SystemRadar) _send() {
	sr.Logger.Debug(
		"system_radar.send",
		"sending system info",
		map[string]interface{}{
			"systemradar_worker_id": sr.workerId,
		},
	)

	if err := sr.Database.SendSystemRadarScannerData(map[string]interface{}{
		"SystemRadarId":     sr.Feedback.SystemRadarId,
		"HostInfo":          sr.Feedback.HostInfo,
		"CpuUsageAll":       sr.Feedback.CpuUsageAll,
		"DiskUsage":         sr.Feedback.DiskUsage,
		"NetworkStat":       sr.Feedback.NetworkStat,
		"CpuInfo":           sr.Feedback.CpuInfo,
		"MemoryInfo":        sr.Feedback.MemoryInfo,
		"NetworkInterfaces": sr.Feedback.NetworkInterfaces,
	}); err != nil {
		sr.Logger.Error(
			"system_radar.send",
			err.Error(),
			map[string]interface{}{
				"systemradar_worker_id": sr.workerId,
			},
		)
	}

	sr.Logger.Debug(
		"system_radar.send",
		"data sent without any problems",
		map[string]interface{}{
			"systemradar_worker_id": sr.workerId,
		},
	)
}
