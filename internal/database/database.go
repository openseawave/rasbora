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
	"openseawave.com/rasbora/internal/data"
)

// Database holds an instance.
type Database struct {
	databaseManager Interface
}

// Interface defines the methods that a database should implement.
type Interface interface {
	SendHeartbeat(workerId, workerType string) error
	Enqueue(queueName string, item data.Queueable) error
	Dequeue(queueName string, workerId string) (item data.Queueable, err error)
	Failed(queueName string, item data.Queueable, err error) error
	Finished(queueName string, item data.Queueable) error
	Processing(queueName string, data map[string]interface{}) error
	TotalRetry(queueName string, item data.Queueable) int
	SendSystemRadarScannerData(data map[string]interface{}) error
	SendLogsToDatabase(log map[string]interface{}) error
}

// New create new database instance.
func New(databaseManager Interface) *Database {
	return &Database{
		databaseManager: databaseManager,
	}
}

// SendHeartbeat send heartbeat to update cluster status.
func (d *Database) SendHeartbeat(workerId, workerType string) error {
	return d.databaseManager.SendHeartbeat(workerId, workerType)
}

// Enqueue add item to waiting queue.
func (d *Database) Enqueue(queueName string, item data.Queueable) error {
	return d.databaseManager.Enqueue(queueName, item)
}

// Dequeue fetch item from waiting queue.
func (d *Database) Dequeue(queueName string, workerId string) (item data.Queueable, err error) {
	return d.databaseManager.Dequeue(queueName, workerId)
}

// Failed change item status to failed.
func (d *Database) Failed(queueName string, item data.Queueable, err error) error {
	return d.databaseManager.Failed(queueName, item, err)
}

// Finished change item status to finished.
func (d *Database) Finished(queueName string, item data.Queueable) error {
	return d.databaseManager.Finished(queueName, item)
}

// Processing send single content a processing status.
func (d *Database) Processing(queueName string, data map[string]interface{}) error {
	return d.databaseManager.Processing(queueName, data)
}

// TotalRetry get total failed retry.
func (d *Database) TotalRetry(queueName string, item data.Queueable) int {
	return d.databaseManager.TotalRetry(queueName, item)
}

// SendSystemRadarScannerData send system radar scanning data content full information about running node.
func (d *Database) SendSystemRadarScannerData(data map[string]interface{}) error {
	return d.databaseManager.SendSystemRadarScannerData(data)
}

// SendLogsToDatabase  save rasbora logs at database.
func (d *Database) SendLogsToDatabase(log map[string]interface{}) error {
	return d.databaseManager.SendLogsToDatabase(log)
}
