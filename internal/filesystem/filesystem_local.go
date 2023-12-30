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
	"io"
	"openseawaves.com/rasbora/internal/data"
	"os"
)

type LocalFileSystem struct{}

// RemoveAll remove all files included inside folder and folder itself.
func (lfs *LocalFileSystem) RemoveAll(folder data.File) error {
	return os.RemoveAll(folder.FullPath())
}

// RemoveFile remove single file
func (lfs *LocalFileSystem) RemoveFile(file data.File) error {
	return os.Remove(file.FullPath())
}

// GetFile put file to other destination.
func (lfs *LocalFileSystem) GetFile(file data.File, saveAt data.File) error {

	// Open the source file for reading
	sourceFile, err := os.Open(file.FullPath())
	if err != nil {
		return err
	}

	defer func(sourceFile *os.File) {
		_ = sourceFile.Close()
	}(sourceFile)

	// Create or open the destination file for writing
	destinationFile, err := os.Create(saveAt.FullPath())
	if err != nil {
		return err
	}

	defer func(destinationFile *os.File) {
		_ = destinationFile.Close()
	}(destinationFile)

	// Copy the content from the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// PutFile put file to other destination.
func (lfs *LocalFileSystem) PutFile(file data.File, saveAt data.File) error {
	return os.Rename(file.FullPath(), saveAt.FullPath())
}
