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

package config

// Config holds an instance
type Config struct {
	configManager Interface
}

// Interface defines the methods that a config should implement.
type Interface interface {
	GetStringSlice(key string) []string
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetIntSlice(key string) []int
}

// New create config instance.
func New(configManager Interface) *Config {
	return &Config{
		configManager: configManager,
	}
}

// GetIntSlice get list of int values by config key
func (c *Config) GetIntSlice(key string) []int {
	return c.configManager.GetIntSlice(key)
}

// GetStringSlice get list of strings values by config key
func (c *Config) GetStringSlice(key string) []string {
	return c.configManager.GetStringSlice(key)
}

// GetString get string value by config key
func (c *Config) GetString(key string) string {
	return c.configManager.GetString(key)
}

// GetBool bet boolean value by config key
func (c *Config) GetBool(key string) bool {
	return c.configManager.GetBool(key)
}

// GetInt get int value by config key
func (c *Config) GetInt(key string) int {
	return c.configManager.GetInt(key)
}
