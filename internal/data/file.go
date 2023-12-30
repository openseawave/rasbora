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

import (
	"encoding/json"
	"path/filepath"
)

// File holds instances
type File struct {
	FileName string      `json:"file_name"`
	FilePath string      `json:"file_path"`
	FileMeta interface{} `json:"file_meta,omitempty"`
}

func (f File) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f File) FullPath() string {
	return filepath.Join(f.FilePath, f.FileName)
}
