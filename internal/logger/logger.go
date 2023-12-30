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

package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"openseawaves.com/rasbora/internal/data"
	"os"
	"runtime"
	"strings"
	"time"

	"openseawaves.com/rasbora/internal/utilities"
)

// Levels
const (
	SUCCESS int = 1
	INFO    int = 2
	WARN    int = 3
	ERROR   int = 4
	DEBUG   int = 5
)

// Color constants
const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
	ColorGrey   = "\033[90m"
)

// Logger is a struct representing a logger instance.
type Logger struct {
	Options Options
}

// Options contains configuration options for the logger.
type Options struct {
	Level  []int
	Output []data.LoggerOutputConfig
}

// Interface defines the methods that a logger should implement.
type Interface interface {
	Error(label, msg string)
	Success(label, msg string)
	Info(label, msg string)
	Warn(label, msg string)
	Debug(label, msg string)
}

// New Create new logger instance
func New() *Logger {
	return &Logger{
		Options: Options{
			Output: []data.LoggerOutputConfig{
				{
					OutputType:   data.StdoutLoggerOutputType,
					OutputWriter: os.Stdout,
				},
			},
			Level: []int{ERROR, SUCCESS, WARN, INFO, DEBUG},
		},
	}
}

// NewWithConfig Creates and returns a new Logger with custom config.
func NewWithConfig(loggerOptions Options) *Logger {
	return &Logger{
		Options: loggerOptions,
	}
}

// Error print error log message
func (l *Logger) Error(label, msg string, extra interface{}) {
	if !utilities.InSlice(ERROR, l.Options.Level) {
		return
	}

	l._processLog("ERROR", label, msg, ColorRed, extra)
}

// Success print success log message
func (l *Logger) Success(label, msg string, extra interface{}) {
	if !utilities.InSlice(SUCCESS, l.Options.Level) {
		return
	}

	l._processLog("SUCCESS", label, msg, ColorGreen, extra)
}

// Info print info log message
func (l *Logger) Info(label, msg string, extra interface{}) {
	if !utilities.InSlice(INFO, l.Options.Level) {
		return
	}

	l._processLog("INFO", label, msg, ColorBlue, extra)
}

// Warn print warning log message
func (l *Logger) Warn(label, msg string, extra interface{}) {
	if !utilities.InSlice(WARN, l.Options.Level) {
		return
	}

	l._processLog("WARN", label, msg, ColorYellow, extra)
}

// Debug print debug log message
func (l *Logger) Debug(label, msg string, extra interface{}) {
	if !utilities.InSlice(DEBUG, l.Options.Level) {
		return
	}

	l._processLog("DEBUG", label, msg, ColorGrey, extra)
}

// logMessage formats and prints the log message with appropriate styles.
func (l *Logger) _processLog(prefix, label, msg, colorCode string, extra interface{}) {
	resetColor := ColorReset

	// Retrieve caller information (file and line number)
	// Extract the base file name from the full path
	_, file, line, _ := runtime.Caller(2)
	fileParts := strings.Split(file, "/")
	fileName := fileParts[len(fileParts)-1]
	formattedTime := time.Now().UTC().Format("[2006/01/02 15:04:05]")

	//print output or use io.writer to write logs anywhere.
	for _, output := range l.Options.Output {

		switch output.OutputType {
		case data.FileLoggerOutputType:
			// save logs on file
			extraAsString := new(bytes.Buffer)

			if extra != nil {
				for key, value := range extra.(map[string]interface{}) {
					_, _ = fmt.Fprintf(extraAsString, "[%s=%s]", key, value)
				}
			} else {
				_, _ = fmt.Fprintf(extraAsString, "")
			}

			_, _ = fmt.Fprintln(output.OutputWriter, fmt.Sprintf(
				"%v %v [%s] [%v:%v] %v %v",
				formattedTime,
				prefix,
				strings.ToLower(label),
				fileName,
				line,
				msg,
				extraAsString,
			))

		case data.DatabaseLoggerOutputType:
			// save logs in database
			extraAsJson, _ := json.Marshal(extra)

			log := map[string]interface{}{
				"type":    prefix,
				"label":   strings.ToLower(label),
				"file":    fileName,
				"line":    line,
				"message": msg,
				"extra":   extraAsJson,
			}

			logAsJson, _ := json.Marshal(log)

			_, _ = fmt.Fprintln(output.OutputWriter, string(logAsJson))

			break
		default:
			// print log on terminal
			extraAsString := new(bytes.Buffer)

			if extra != nil {
				for key, value := range extra.(map[string]interface{}) {
					_, _ = fmt.Fprintf(extraAsString, "[%s=%s]", key, value)
				}
			} else {
				_, _ = fmt.Fprintf(extraAsString, "")
			}

			_, _ = fmt.Fprintln(output.OutputWriter, fmt.Sprintf(
				"%s %s%s %s[%s] [%v:%v]%s %s %v%v%s",
				formattedTime,
				colorCode,
				prefix,
				ColorReset,
				strings.ToLower(label),
				fileName, line,
				resetColor,
				msg,
				ColorYellow,
				extraAsString,
				ColorReset,
			))

		}

	}
}
