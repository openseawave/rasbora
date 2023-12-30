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

package taskmanager

import "context"

// Name used as identifier.
const Name = "TaskManagement"

// TaskManager hold an instance
type TaskManager struct {
	manager Interface
}

// Interface defines the methods that task manager should implement.
type Interface interface {
	StartTaskManager(ctx context.Context)
}

// New create new instance
func New(taskManager Interface) *TaskManager {
	return &TaskManager{
		manager: taskManager,
	}
}

// StartTaskManager start task manager to listen for task operations.
func (tm *TaskManager) StartTaskManager(ctx context.Context) {
	tm.manager.StartTaskManager(ctx)
}
