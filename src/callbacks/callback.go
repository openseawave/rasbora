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

package callbacks

import "context"

// Name used as identifier.
const Name = "CallbackManager"

// CallbackManager holds an instance
type CallbackManager struct {
	manager Interface
}

// Interface defines the methods that callback manager should implement.
type Interface interface {
	StartCallbackManager(ctx context.Context)
}

// New create new callback manager instance.
func New(callbackManager Interface) *CallbackManager {
	return &CallbackManager{
		manager: callbackManager,
	}
}

// StartCallbackManager start process callbacks.
func (cm *CallbackManager) StartCallbackManager(ctx context.Context) {
	cm.manager.StartCallbackManager(ctx)
}
