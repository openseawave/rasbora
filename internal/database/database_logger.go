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

package database

import (
	"encoding/json"
	"fmt"

	"openseawave.com/rasbora/internal/config"
)

type LoggerOutput struct {
	Config   *config.Config
	Database *Database
}

// Write implements the io.Writer interface for RedisStreamWriter.
func (lo LoggerOutput) Write(log []byte) (n int, err error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("can not write logs to database", r)
		}
	}()

	var data map[string]interface{}

	err = json.Unmarshal(log, &data)

	if err != nil {
		return 0, err
	}

	data["logger_id"] = lo.Config.GetString("Logger.UniqueID")

	err = lo.Database.SendLogsToDatabase(data)
	if err != nil {
		return 0, err
	}

	return len(log), nil
}
