// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Set up logger with config
func SetUpLogger(config *LoggerConfig) {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeTime  = zapcore.ISO8601TimeEncoder

	cores := []zapcore.Core{}

	if config.File.Enabled {
		core := createFileLoggerCore(&config.File, zapcore.NewJSONEncoder(cfg))
		cores = append(cores, core)
	}

	if config.Console.Enabled {
		consoleLevel := ParseLevels(config.Console.Level)
		core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.Lock(os.Stdout), consoleLevel)
		cores = append(cores, core)
	}

	extras := len(config.extraFiles)
	if extras > 0 {
		for i := 0; i < extras; i++ {
			fileCfg := config.extraFiles[i]
			if fileCfg.Enabled {
				core := createFileLoggerCore(fileCfg, zapcore.NewJSONEncoder(cfg))
				cores = append(cores, core)
			}
		}
	}

	if len(cores) > 0 {
		core := zapcore.NewTee(cores...)

		logger := zap.New(core)
		Logger = logger

		SLogger = logger.Sugar()
	}
}

func createFileLoggerCore(config *FileConfig, enc zapcore.Encoder) zapcore.Core {
	hook := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	fileLevel := ParseLevels(config.Level)

	return zapcore.NewCore(enc, zapcore.AddSync(hook), fileLevel)
}
