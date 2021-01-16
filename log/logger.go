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

var (
	Logger *zap.Logger = zap.L()
	SLogger *zap.SugaredLogger = zap.S()
)

// --- Convenient SugaredLogger methods ---

func Debug(args ...interface{}) {
	SLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	SLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	SLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	SLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	SLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	SLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	SLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	SLogger.Errorf(template, args...)
}

func Panic(args ...interface{}) {
	SLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	SLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	SLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	SLogger.Fatalf(template, args...)
}

// --- Get global logger instance ---

func L() *zap.Logger {
	return Logger
}

func S() *zap.SugaredLogger {
	return SLogger
}

// --- Set up logger ---

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
