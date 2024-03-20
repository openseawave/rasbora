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

package config

import (
	"github.com/spf13/viper"
)

// ViperConfigManager holds an instance
type ViperConfigManager struct {
	Viper *viper.Viper
}

// GetIntSlice get list of int values by config key
func (v *ViperConfigManager) GetIntSlice(key string) []int {
	var _value []int
	if err := v.Viper.UnmarshalKey(key, &_value); err != nil {
		return []int{}
	}
	return _value
}

// GetStringSlice get list of strings values by config key
func (v *ViperConfigManager) GetStringSlice(key string) []string {
	var _value []string
	if err := v.Viper.UnmarshalKey(key, &_value); err != nil {
		return []string{}
	}
	return _value
}

// GetString get string value by config key
func (v *ViperConfigManager) GetString(key string) string {
	return v.Viper.GetString(key)
}

// GetBool bet boolean value by config key
func (v *ViperConfigManager) GetBool(key string) bool {
	return v.Viper.GetBool(key)
}

// GetInt get int value by config key
func (v *ViperConfigManager) GetInt(key string) int {
	return v.Viper.GetInt(key)
}
