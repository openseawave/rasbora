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

package logger

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"openseawave.com/rasbora/internal/data"
)

func TestNewLogger(t *testing.T) {
	logger := New()
	if logger == nil {
		t.Error("New() returned a nil logger")
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.BufferLoggerOutputType,
				OutputWriter: &buf,
			},
		},
		Level: []int{ERROR},
	})

	logger.Error("error_tag", "an error occurred", map[string]interface{}{})

	output := buf.String()

	if !strings.Contains(output, "ERROR") {
		t.Errorf("Log level not as expected: %s", output)
	}

	if !strings.Contains(output, "[error_tag]") {
		t.Errorf("Log tag not as expected: %s", output)
	}

	if !strings.Contains(output, "an error occurred") {
		t.Errorf("Log msg not as expected: %s", output)
	}
}

func TestLogger_Success(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.BufferLoggerOutputType,
				OutputWriter: &buf,
			},
		},
		Level: []int{SUCCESS},
	})

	logger.Success("success_tag", "operation successful", map[string]interface{}{})

	output := buf.String()

	if !strings.Contains(output, "SUCCESS") {
		t.Errorf("Log level not as expected: %s", output)
	}

	if !strings.Contains(output, "[success_tag]") {
		t.Errorf("Log tag not as expected: %s", output)
	}

	if !strings.Contains(output, "operation successful") {
		t.Errorf("Log msg not as expected: %s", output)
	}
}

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.BufferLoggerOutputType,
				OutputWriter: &buf,
			},
		},
		Level: []int{INFO},
	})

	logger.Info("info_tag", "information message", map[string]interface{}{})

	output := buf.String()

	if !strings.Contains(output, "INFO") {
		t.Errorf("Log level not as expected: %s", output)
	}

	if !strings.Contains(output, "[info_tag]") {
		t.Errorf("Log tag not as expected: %s", output)
	}

	if !strings.Contains(output, "information message") {
		t.Errorf("Log msg not as expected: %s", output)
	}
}

func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer

	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.BufferLoggerOutputType,
				OutputWriter: &buf,
			},
		},
		Level: []int{WARN},
	})

	logger.Warn("warn_tag", "warning message", map[string]interface{}{})

	output := buf.String()

	if !strings.Contains(output, "WARN") {
		t.Errorf("Log level not as expected: %s", output)
	}

	if !strings.Contains(output, "[warn_tag]") {
		t.Errorf("Log tag not as expected: %s", output)
	}

	if !strings.Contains(output, "warning message") {
		t.Errorf("Log msg not as expected: %s", output)
	}
}

func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer

	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.BufferLoggerOutputType,
				OutputWriter: &buf,
			},
		},
		Level: []int{DEBUG},
	})

	logger.Debug("debug_tag", "debug message", map[string]interface{}{})

	output := buf.String()

	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Log level not as expected: %s", output)
	}

	if !strings.Contains(output, "[debug_tag]") {
		t.Errorf("Log tag not as expected: %s", output)
	}

	if !strings.Contains(output, "debug message") {
		t.Errorf("Log msg not as expected: %s", output)
	}
}

func TestNewLoggerWithConfig(t *testing.T) {
	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.StdoutLoggerOutputType,
				OutputWriter: os.Stdout,
			},
		},
		Level: []int{ERROR, SUCCESS, WARN, INFO},
	})

	if logger == nil {
		t.Error("NewWithConfig() returned a nil logger")
	}
}

func BenchmarkLogger(b *testing.B) {
	buf := &bytes.Buffer{}

	logOptions := Options{
		Level: []int{ERROR, SUCCESS, WARN, INFO},
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.BufferLoggerOutputType,
				OutputWriter: buf,
			},
		},
	}

	logger := NewWithConfig(logOptions)

	label := "BENCHMARK_FUNC"

	b.ResetTimer()
	for _, bm := range []struct {
		level string
		msg   string
	}{
		{"INFO", "BENCHMARK INFO MESSAGE"},
		{"ERROR", "BENCHMARK ERROR MESSAGE"},
		{"SUCCESS", "BENCHMARK SUCCESS MESSAGE"},
		{"WARN", "BENCHMARK WARN MESSAGE"},
		{"DEBUG", "BENCHMARK DEBUG MESSAGE"},
	} {
		b.Run(fmt.Sprintf("Benchmark-%s", bm.level), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if bm.level == "INFO" {
					logger.Info(label, bm.msg, map[string]interface{}{})
				}
				if bm.level == "ERROR" {
					logger.Error(label, bm.msg, map[string]interface{}{})
				}

				if bm.level == "SUCCESS" {
					logger.Success(label, bm.msg, map[string]interface{}{})
				}

				if bm.level == "WARN" {
					logger.Warn(label, bm.msg, map[string]interface{}{})
				}

				if bm.level == "DEBUG" {
					logger.Debug(label, bm.msg, map[string]interface{}{})
				}
			}
		})
	}
}

func ExampleLogger_Info() {
	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.StdoutLoggerOutputType,
				OutputWriter: os.Stdout,
			},
		},
		Level: []int{INFO},
	})
	logger.Info("internal.logger.testing", "this is info log message.", map[string]interface{}{})
}

func ExampleLogger_Error() {
	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.StdoutLoggerOutputType,
				OutputWriter: os.Stdout,
			},
		},
		Level: []int{ERROR},
	})
	logger.Error("internal.logger.testing", "this is error log message.", map[string]interface{}{})
}

func ExampleLogger_Success() {
	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.StdoutLoggerOutputType,
				OutputWriter: os.Stdout,
			},
		},
		Level: []int{SUCCESS},
	})
	logger.Success("internal.logger.testing", "this is success log message.", map[string]interface{}{})
}

func ExampleLogger_Warn() {
	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.StdoutLoggerOutputType,
				OutputWriter: os.Stdout,
			},
		},
		Level: []int{WARN},
	})
	logger.Warn("internal.logger.testing", "this is warn log message.", map[string]interface{}{})
}

func ExampleLogger_Debug() {
	logger := NewWithConfig(Options{
		Output: []data.LoggerOutputConfig{
			{
				OutputType:   data.StdoutLoggerOutputType,
				OutputWriter: os.Stdout,
			},
		},
		Level: []int{DEBUG},
	})
	logger.Warn("internal.logger.testing", "this is debug log message.", map[string]interface{}{})
}
