// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"go.uber.org/zap"
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
