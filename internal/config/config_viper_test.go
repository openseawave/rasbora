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
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestViperConfigManager_GetStringSlice(t *testing.T) {
	// Initialize ViperConfigManager with a mock Viper instance
	mockViper := viper.New()
	mockViper.Set("key", []string{"value1", "value2"})
	configManager := &ViperConfigManager{Viper: mockViper}

	// Test GetStringSlice
	result := configManager.GetStringSlice("key")
	expected := []string{"value1", "value2"}
	assert.Equal(t, expected, result)
}

func TestViperConfigManager_GetString(t *testing.T) {
	// Initialize ViperConfigManager with a mock Viper instance
	mockViper := viper.New()
	mockViper.Set("key", "value")
	configManager := &ViperConfigManager{Viper: mockViper}

	// Test GetString
	result := configManager.GetString("key")
	expected := "value"
	assert.Equal(t, expected, result)
}

func TestViperConfigManager_GetBool(t *testing.T) {
	// Initialize ViperConfigManager with a mock Viper instance
	mockViper := viper.New()
	mockViper.Set("key", true)
	configManager := &ViperConfigManager{Viper: mockViper}

	// Test GetBool
	result := configManager.GetBool("key")
	assert.Equal(t, true, result)
}

func TestViperConfigManager_GetInt(t *testing.T) {
	// Initialize ViperConfigManager with a mock Viper instance
	mockViper := viper.New()
	mockViper.Set("key", 42)
	configManager := &ViperConfigManager{Viper: mockViper}

	// Test GetInt
	result := configManager.GetInt("key")
	expected := 42
	assert.Equal(t, expected, result)
}
