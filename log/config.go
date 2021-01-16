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

// Config for file logger
type FileConfig struct {
	Enabled    bool   `json:"enabled"`
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// New a FileConfig
func NewFileConfig(level, filename string) *FileConfig {
	config := &FileConfig{}

	config.init()

	config.Enabled = true
	config.Level = level
	config.Filename = filename

	return config
}

func (config *FileConfig) init() {
	config.Enabled = false
	config.Level = DEFAULT_LEVEL
	config.Filename = "/tmp/logs/default.log"
	config.MaxSize = 128 // megabytes
	config.MaxBackups = 3
	config.MaxAge = 28 // days
	config.Compress = false
}

// Config for console logger
type ConsoleConfig struct {
	Enabled    bool   `json:"enabled"`
	Level      string `json:"level"`
}

func (config *ConsoleConfig) init() {
	config.Enabled = true
	config.Level = DEBUG_LEVEL
}

// Logger configuration
type LoggerConfig struct {
	File       FileConfig    `json:"file"`
	Console    ConsoleConfig `json:"console"`
	Level      string        `json:"level"`
	extraFiles []*FileConfig `json:"-"`
}

// Initial the LoggerConfig with default value.
func (config *LoggerConfig) Init() {
	config.File.init()
	config.Console.init()
	config.Level = DEFAULT_LEVEL
}

// Append extra FileConfig(s)
func (config *LoggerConfig) AppendFileConfig(cfg ...*FileConfig) {
	config.extraFiles = append(config.extraFiles, cfg...)
}

// Create a LoggerConfig with default value.
func DefaultConfig() *LoggerConfig {
	config := &LoggerConfig{
	}

	config.Init()

	return config
}
