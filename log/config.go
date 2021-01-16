// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

const (
	DEFAULT_LEVEL = "INFO"
	DEBUG_LEVEL = "DEBUG"
	INFO_LEVEL  = "INFO"
	WARN_LEVEL  = "WARN"
	ERROR_LEVEL = "ERROR"
	PANIC_LEVEL = "PANIC"
	FATAL_LEVEL = "FATAL"
)

// Config file logger
type FileLogger struct {
	Enabled    bool   `json:"enabled"`
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

func (config *FileLogger) init() {
	config.Enabled = false
	config.Level = DEFAULT_LEVEL
	config.Filename = "/tmp/logs/default.log"
	config.MaxSize = 128 // megabytes
	config.MaxBackups = 3
	config.MaxAge = 28 // days
	config.Compress = false
}

// Config console logger
type ConsoleLogger struct {
	Enabled    bool   `json:"enabled"`
	Level      string `json:"level"`
}

func (config *ConsoleLogger) init() {
	config.Enabled = true
	config.Level = DEBUG_LEVEL
}

// Logger configuration
type LoggerConfig struct {
	File    FileLogger     `json:"file"`
	Console ConsoleLogger  `json:"console"`
	Level   string         `json:"level"`
}

// Initial the LoggerConfig with default value.
func (config *LoggerConfig) Init() {
	config.File.init()
	config.Console.init()
	config.Level = DEFAULT_LEVEL
}

// Create a LoggerConfig with default value.
func DefaultConfig() *LoggerConfig {
	config := &LoggerConfig{
	}

	config.Init()

	return config
}
