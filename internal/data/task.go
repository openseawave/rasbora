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

// Task holds instances
type Task struct {
	// Unique identifier for the task.
	ID string `json:"task_id"`

	// Label for task.
	Label string `json:"task_label" validate:"required"`

	// Priority level assigned to the task.
	Priority *float64 `json:"task_priority" validate:"required"`

	// Callback struct holds details for a callback associated with the task.
	Callback struct {
		// URL to send callback.
		URL string `json:"callback_url" validate:"required"`
		// Data to be sent as part of the callback.
		Data interface{} `json:"callback_data" validate:"required"`
	} `json:"callback"`

	// VideoTranscoder contains details about video transcoding for the task.
	VideoTranscoder struct {
		InputVideo struct {
			// File system type.
			FileSystem FileSystemType `json:"input_file_system" validate:"required"`
			// Name of the input video file.
			FileName string `json:"input_file_name" validate:"required"`
			// Path to the input video file.
			FilePath string `json:"input_file_path" validate:"required"`
		} `json:"input"`

		//holds information how should be the video output.
		Output struct {
			Handler   string                   `json:"handler" validate:"required"`
			Container string                   `json:"container" validate:"required"`
			Args      []map[string]interface{} `json:"args" validate:"required"`
		} `json:"output"`
	} `json:"video_transcoder"`

	// Timestamp indicating when the task was created.
	CreatedAt int64 `json:"created_at,omitempty"`
	// Timestamp indicating when the task started.
	StartedAt int64 `json:"started_at,omitempty"`
	// Timestamp indicating when the task finished.
	FinishedAt int64 `json:"finished_at,omitempty"`
	// Timestamp indicating when the task failed.
	FailedAt int64 `json:"failed_at,omitempty"`
}

func (i Task) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
