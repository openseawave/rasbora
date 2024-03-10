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
	"fmt"
	"testing"
)

// MockConfigManager implements the Interface for testing purposes.
type MockConfigManager struct {
	data map[string]interface{}
}

func (m *MockConfigManager) GetIntSlice(key string) []int {
	if val, ok := m.data[key].([]int); ok {
		return val
	}
	return nil
}

func (m *MockConfigManager) GetStringSlice(key string) []string {
	if val, ok := m.data[key].([]string); ok {
		return val
	}
	return nil
}

func (m *MockConfigManager) GetString(key string) string {
	if val, ok := m.data[key].(string); ok {
		return val
	}
	return ""
}

func (m *MockConfigManager) GetBool(key string) bool {
	if val, ok := m.data[key].(bool); ok {
		return val
	}
	return false
}

func (m *MockConfigManager) GetInt(key string) int {
	if val, ok := m.data[key].(int); ok {
		return val
	}
	return 0
}

func TestConfig_GetIntSlice(t *testing.T) {
	mockData := map[string]interface{}{
		"intSliceKey": []int{1, 2, 3},
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	result := config.GetIntSlice("intSliceKey")

	expected := []int{1, 2, 3}
	if len(result) != len(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected, result)
			break
		}
	}
}

func TestConfig_GetStringSlice(t *testing.T) {
	mockData := map[string]interface{}{
		"stringSliceKey": []string{"value1", "value2", "value3"},
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	result := config.GetStringSlice("stringSliceKey")

	expected := []string{"value1", "value2", "value3"}
	if len(result) != len(expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected %v, but got %v", expected, result)
			break
		}
	}
}

func TestConfig_GetString(t *testing.T) {
	mockData := map[string]interface{}{
		"stringKey": "testValue",
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	result := config.GetString("stringKey")

	expected := "testValue"
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestConfig_GetBool(t *testing.T) {
	mockData := map[string]interface{}{
		"boolKey": true,
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	result := config.GetBool("boolKey")

	if result != true {
		t.Errorf("Expected %v, but got %v", true, result)
	}
}

func TestConfig_GetInt(t *testing.T) {
	mockData := map[string]interface{}{
		"intKey": 42,
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	result := config.GetInt("intKey")

	expected := 42
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func BenchmarkConfig_GetIntSlice(b *testing.B) {
	mockData := map[string]interface{}{
		"intSliceKey": []int{1, 2, 3},
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GetIntSlice("intSliceKey")
	}
}

func BenchmarkConfig_GetStringSlice(b *testing.B) {
	mockData := map[string]interface{}{
		"stringSliceKey": []string{"value1", "value2", "value3"},
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GetStringSlice("stringSliceKey")
	}
}

func BenchmarkConfig_GetString(b *testing.B) {
	mockData := map[string]interface{}{
		"stringKey": "testValue",
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GetString("stringKey")
	}
}

func BenchmarkConfig_GetBool(b *testing.B) {
	mockData := map[string]interface{}{
		"boolKey": true,
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GetBool("boolKey")
	}
}

func BenchmarkConfig_GetInt(b *testing.B) {
	mockData := map[string]interface{}{
		"intKey": 42,
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.GetInt("intKey")
	}
}

func ExampleConfig_GetIntSlice() {
	mockData := map[string]interface{}{
		"intSliceKey": []int{1, 2, 3},
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	result := config.GetIntSlice("intSliceKey")
	fmt.Println(result)
	// Output: [1 2 3]
}

func ExampleConfig_GetStringSlice() {
	mockData := map[string]interface{}{
		"stringSliceKey": []string{"value1", "value2", "value3"},
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	result := config.GetStringSlice("stringSliceKey")
	fmt.Println(result)
	// Output: [value1 value2 value3]
}

func ExampleConfig_GetString() {
	mockData := map[string]interface{}{
		"stringKey": "testValue",
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	result := config.GetString("stringKey")
	fmt.Println(result)
	// Output: testValue
}

func ExampleConfig_GetBool() {
	mockData := map[string]interface{}{
		"boolKey": true,
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	result := config.GetBool("boolKey")
	fmt.Println(result)
	// Output: true
}

func ExampleConfig_GetInt() {
	mockData := map[string]interface{}{
		"intKey": 42,
	}
	mockConfigManager := &MockConfigManager{data: mockData}
	config := New(mockConfigManager)

	result := config.GetInt("intKey")
	fmt.Println(result)
	// Output: 42
}
