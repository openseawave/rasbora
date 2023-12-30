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

package filesystem

import (
	"openseawaves.com/rasbora/internal/data"
)

// FileSystem holds an instance
type FileSystem struct {
	fileManager Interface
}

// Interface defines the methods that a file system should implement.
type Interface interface {
	RemoveFile(file data.File) error
	RemoveAll(path data.File) error
	GetFile(file data.File, saveAt data.File) error
	PutFile(file data.File, saveAt data.File) error
}

// NewFileSystem create new file system instance.
func NewFileSystem(fileManager Interface) *FileSystem {
	return &FileSystem{
		fileManager: fileManager,
	}
}

// RemoveFile file from file system
func (f *FileSystem) RemoveFile(file data.File) error {
	return f.fileManager.RemoveFile(file)
}

// RemoveAll remove all files included inside
func (f *FileSystem) RemoveAll(path data.File) error {
	return f.fileManager.RemoveAll(path)
}

// GetFile will move get file from outside system
func (f *FileSystem) GetFile(file data.File, saveAt data.File) error {
	return f.fileManager.GetFile(file, saveAt)
}

// PutFile will put file inside file system
func (f *FileSystem) PutFile(file data.File, saveAt data.File) error {
	return f.fileManager.PutFile(file, saveAt)
}
