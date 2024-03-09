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

package videotranscoder

import "context"

// Name used as identifier.
const Name = "VideoTranscoding"

// Transcoder holds an instance.
type Transcoder struct {
	engine Interface
}

// Interface defines the methods that a video transcoder engine should implement.
type Interface interface {
	StarTranscoderEngine(ctx context.Context)
}

// New create video transcoder instance
func New(transcoderEngine Interface) *Transcoder {
	return &Transcoder{
		engine: transcoderEngine,
	}
}

// StarTranscoderEngine listen for new task to transcode them.
func (tm *Transcoder) StarTranscoderEngine(ctx context.Context) {
	tm.engine.StarTranscoderEngine(ctx)
}
