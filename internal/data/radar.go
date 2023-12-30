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

package data

import "encoding/json"

// Radar holds instances
type Radar struct {
	SystemRadarId     string      `json:"SystemRadarId"`
	MemoryInfo        interface{} `json:"MemoryInfo"`
	CpuInfo           interface{} `json:"CpuInfo"`
	CpuUsageAll       interface{} `json:"CpuUsageAll"`
	DiskUsage         interface{} `json:"DiskUsage"`
	NetworkStat       interface{} `json:"NetworkStat"`
	HostInfo          interface{} `json:"HostInfo"`
	NetworkInterfaces interface{} `json:"NetworkInterfaces"`
}

func (r Radar) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}
