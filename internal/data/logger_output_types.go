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

import "io"

// LoggerOutputType holds instances
type LoggerOutputType string

const (
	StdoutLoggerOutputType   LoggerOutputType = "stdout"
	DatabaseLoggerOutputType LoggerOutputType = "database"
	FileLoggerOutputType     LoggerOutputType = "file"
	BufferLoggerOutputType   LoggerOutputType = "buffer"
)

// String returns the string representation of LoggerOutputWriter.
func (low LoggerOutputType) String() string {
	return string(low)
}

// LoggerOutputConfig holds instances
type LoggerOutputConfig struct {
	OutputWriter io.Writer
	OutputType   LoggerOutputType
}
