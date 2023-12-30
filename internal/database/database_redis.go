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

package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"openseawaves.com/rasbora/internal/config"
	"openseawaves.com/rasbora/internal/data"
	"openseawaves.com/rasbora/internal/logger"
)

// Create a new context based on the Background context
var ctx = context.Background()

// RedisDatabaseManager holds an instance
type RedisDatabaseManager struct {
	Redis  *redis.Client
	Config *config.Config
	Logger *logger.Logger
}

// SendLogsToDatabase save rasbora logs at database.
func (rdm *RedisDatabaseManager) SendLogsToDatabase(log map[string]interface{}) error {
	res := rdm.Redis.XAdd(ctx, &redis.XAddArgs{
		Stream: rdm.Config.GetString("Database.Redis.Structure.Logger"),
		Values: log,
	})

	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

// SendHeartbeat send heartbeat to update cluster status.
func (rdm *RedisDatabaseManager) SendHeartbeat(workerId, workerType string) error {
	clusterHeartbeatList := rdm.Config.GetString("Database.Redis.Structure.Cluster.Heartbeat")

	workerWithType := fmt.Sprintf("%s:%s", workerType, workerId)

	if zAddResult := rdm.Redis.ZAdd(
		ctx,
		clusterHeartbeatList,
		redis.Z{
			Score:  float64(time.Now().UnixMilli()),
			Member: workerWithType,
		},
	); zAddResult.Err() != nil {
		return zAddResult.Err()
	}

	return nil
}

// Enqueue add item to waiting queue.
func (rdm *RedisDatabaseManager) Enqueue(queueName string, item data.Queueable) error {
	scoreWithID := fmt.Sprintf("%d:%s", time.Now().UnixMilli(), item.ID)
	waiting, status, worker, _, retry, items, _ := rdm._queueStructures(queueName)

	tx := rdm.Redis.TxPipeline()
	tx.ZAdd(ctx, waiting, redis.Z{Score: item.Priority, Member: scoreWithID})
	tx.HSet(ctx, items, item.ID, item)
	tx.HSet(ctx, status, item.ID, "waiting")
	tx.HSet(ctx, worker, item.ID, nil)
	tx.HIncrBy(ctx, retry, item.ID, 1)

	if _, err := tx.Exec(ctx); err != nil {
		return err
	}

	return nil
}

// Dequeue fetch item from waiting queue.
func (rdm *RedisDatabaseManager) Dequeue(queueName string, workerId string) (item data.Queueable, err error) {
	waiting, status, worker, _, _, items, _ := rdm._queueStructures(queueName)

	// find stopped tasks and make them as failed
	if workerList, err := rdm.Redis.HGetAll(ctx, worker).Result(); err == nil {
		for tID, wID := range workerList {
			if wID == workerId {
				var item data.Queueable
				itemAsJsonString, hGetError := rdm.Redis.HGet(ctx, items, tID).Result()
				if hGetError != nil {
					continue
				}
				jsonParserError := json.Unmarshal([]byte(itemAsJsonString), &item)
				if jsonParserError != nil {
					continue
				}
				if err := rdm.Failed(
					queueName,
					item,
					errors.New("previous worker encountered issues, task returned to waiting queue for redistribution"),
				); err != nil {
					continue
				}
			}
		}
	}

	ids, zPopMinError := rdm.Redis.ZPopMin(ctx, waiting, 1).Result()
	if zPopMinError != nil {
		return item, zPopMinError
	}

	if len(ids) <= 0 {
		return item, errors.New("there is no items in waiting queue")
	}

	var itemID = strings.Split(ids[0].Member.(string), ":")[1]
	itemAsJsonString, hGetError := rdm.Redis.HGet(ctx, items, itemID).Result()
	if hGetError != nil {
		return item, hGetError
	}

	jsonParserError := json.Unmarshal([]byte(itemAsJsonString), &item)
	if jsonParserError != nil {
		return item, jsonParserError
	}

	rdm.Redis.HSet(ctx, status, item.ID, "working")
	rdm.Redis.HSet(ctx, worker, item.ID, workerId)

	return item, nil
}

// Failed change item status to failed.
func (rdm *RedisDatabaseManager) Failed(queueName string, item data.Queueable, err error) error {
	_, status, worker, processing, _, items, logs := rdm._queueStructures(queueName)

	tx := rdm.Redis.TxPipeline()

	tx.Del(ctx, fmt.Sprintf("%v:%v", processing, item.ID))
	tx.HSet(ctx, items, item.ID, item)
	tx.HSet(ctx, status, item.ID, "failed")
	tx.HSet(ctx, logs, item.ID, err.Error())
	tx.HDel(ctx, worker, item.ID)

	if _, err := tx.Exec(ctx); err != nil {
		return err
	}

	return nil
}

// Finished change item status to finished.
func (rdm *RedisDatabaseManager) Finished(queueName string, item data.Queueable) error {
	_, status, worker, _, _, items, _ := rdm._queueStructures(queueName)

	tx := rdm.Redis.TxPipeline()

	tx.HSet(ctx, items, item.ID, item)
	tx.HSet(ctx, status, item.ID, "finished")
	tx.HDel(ctx, worker, item.ID)

	if _, err := tx.Exec(ctx); err != nil {
		return err
	}

	return nil
}

// Processing send single content a processing status.
func (rdm *RedisDatabaseManager) Processing(queueName string, data map[string]interface{}) error {
	_, _, _, processing, _, _, _ := rdm._queueStructures(queueName)

	res := rdm.Redis.XAdd(ctx, &redis.XAddArgs{
		Stream: fmt.Sprintf("%v:%v", processing, data["task_id"]),
		Values: data,
	})

	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

// TotalRetry get total failed retry.
func (rdm *RedisDatabaseManager) TotalRetry(queueName string, item data.Queueable) int {
	_, _, _, _, retry, _, _ := rdm._queueStructures(queueName)

	resH, hGetError := rdm.Redis.HGet(ctx, retry, item.ID).Result()
	if hGetError != nil {
		return -1
	}

	resP, err := strconv.ParseInt(resH, 10, 64)
	if err != nil {
		return -1
	}

	return int(resP)
}

// SendSystemRadarScannerData send system radar scanning data content full information about running node.
func (rdm *RedisDatabaseManager) SendSystemRadarScannerData(scanner map[string]interface{}) error {
	res := rdm.Redis.XAdd(ctx, &redis.XAddArgs{
		Stream: rdm.Config.GetString("Database.Redis.Structure.Cluster.Radar"),
		Values: scanner,
	})

	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

// _queueStructures shortcut to fetch all key names
func (rdm *RedisDatabaseManager) _queueStructures(queueName string) (waiting, status, worker, processing, retry, items, logs string) {
	waiting = strings.Replace(rdm.Config.GetString("Database.Redis.Structure.Queue.Waiting"), "{{name}}", queueName, 1)
	status = strings.Replace(rdm.Config.GetString("Database.Redis.Structure.Queue.Status"), "{{name}}", queueName, 1)
	worker = strings.Replace(rdm.Config.GetString("Database.Redis.Structure.Queue.Worker"), "{{name}}", queueName, 1)
	processing = strings.Replace(rdm.Config.GetString("Database.Redis.Structure.Queue.Processing"), "{{name}}", queueName, 1)
	retry = strings.Replace(rdm.Config.GetString("Database.Redis.Structure.Queue.Retry"), "{{name}}", queueName, 1)
	items = strings.Replace(rdm.Config.GetString("Database.Redis.Structure.Queue.Items"), "{{name}}", queueName, 1)
	logs = strings.Replace(rdm.Config.GetString("Database.Redis.Structure.Queue.Logs"), "{{name}}", queueName, 1)
	return
}
