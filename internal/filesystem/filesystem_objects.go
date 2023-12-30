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
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"openseawaves.com/rasbora/internal/config"
	"openseawaves.com/rasbora/internal/data"
)

// Create a new context based on the Background context
var ctx = context.Background()

// ObjectFileSystem hold an instance
type ObjectFileSystem struct {
	Minio *minio.Client
}

// NewObjectClient create new client preconfigured
func NewObjectClient(cfg config.Config) (minioClient *minio.Client, err error) {
	accessKeyID := cfg.GetString("Filesystem.ObjectStorage.AccessKeyID")
	secretAccessKey := cfg.GetString("Filesystem.ObjectStorage.SecretAccessKey")
	sessionToken := cfg.GetString("Filesystem.ObjectStorage.SessionToken")

	var signatureType = credentials.SignatureAnonymous

	if cfg.GetString("Filesystem.ObjectStorage.Signature") == "v4" {
		signatureType = credentials.SignatureV4
	}

	if cfg.GetString("Filesystem.ObjectStorage.Signature") == "v2" {
		signatureType = credentials.SignatureV2
	}

	if cfg.GetString("Filesystem.ObjectStorage.Signature") == "v4Streaming" {
		signatureType = credentials.SignatureV4Streaming
	}

	if cfg.GetString("Filesystem.ObjectStorage.Signature") == "noSignature" {
		signatureType = credentials.SignatureAnonymous
	}

	minioClient, err = minio.New(cfg.GetString("Filesystem.ObjectStorage.Endpoint"), &minio.Options{
		Creds: credentials.NewStatic(
			accessKeyID,
			secretAccessKey,
			sessionToken,
			signatureType,
		),
		Secure: cfg.GetBool("Filesystem.ObjectStorage.UseSSL"),
	})

	return
}

// RemoveAll remove all objects inside bucket.
func (ofs *ObjectFileSystem) RemoveAll(bucket data.File) error {
	return ofs.Minio.RemoveBucket(ctx, bucket.FilePath)
}

// RemoveFile remove object from bucket
func (ofs *ObjectFileSystem) RemoveFile(object data.File) error {
	return ofs.Minio.RemoveObject(
		ctx,
		object.FilePath,
		object.FileName,
		minio.RemoveObjectOptions{},
	)
}

// GetFile get file to other destination.
func (ofs *ObjectFileSystem) GetFile(object data.File, saveAtLocal data.File) error {
	return ofs.Minio.FGetObject(
		ctx,
		object.FilePath,
		object.FileName,
		saveAtLocal.FullPath(),
		minio.GetObjectOptions{},
	)
}

// PutFile put file to other destination.
func (ofs *ObjectFileSystem) PutFile(localFile data.File, saveAsObject data.File) (err error) {
	_, err = ofs.Minio.FPutObject(
		ctx,
		saveAsObject.FilePath,
		saveAsObject.FileName,
		localFile.FullPath(),
		minio.PutObjectOptions{},
	)
	return err
}
