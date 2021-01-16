// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log_test

import (
	"testing"

	"github.com/chinmobi/modlib/log"

	"go.uber.org/zap/zapcore"
)

func TestEqualOpLevel(t *testing.T) {
	level := log.ParseLevels("=INFO")

	assertTrue(t, level.Enabled(zapcore.InfoLevel), "INFO =INFO")
	assertFalse(t, level.Enabled(zapcore.DebugLevel), "DEBUG =INFO")
	assertFalse(t, level.Enabled(zapcore.WarnLevel), "WARN =INFO")
}

func TestNotOpLevel(t *testing.T) {
	level := log.ParseLevels("!INFO")

	assertFalse(t, level.Enabled(zapcore.InfoLevel), "INFO !INFO")
	assertTrue(t, level.Enabled(zapcore.DebugLevel), "DEBUG !INFO")
	assertTrue(t, level.Enabled(zapcore.WarnLevel), "WARN !INFO")
}

func TestLessOpLevel(t *testing.T) {
	level := log.ParseLevels("<INFO")

	assertFalse(t, level.Enabled(zapcore.InfoLevel), "INFO <INFO")
	assertTrue(t, level.Enabled(zapcore.DebugLevel), "DEBUG <INFO")
	assertFalse(t, level.Enabled(zapcore.WarnLevel), "WARN <INFO")
}

func TestOrOpLevel(t *testing.T) {
	level := log.ParseLevels("=DEBUG ERROR")

	assertTrue(t, level.Enabled(zapcore.DebugLevel), "DEBUG in [=DEBUG, ERROR]")
	assertFalse(t, level.Enabled(zapcore.InfoLevel), "INFO in [=DEBUG, ERROR]")
	assertFalse(t, level.Enabled(zapcore.WarnLevel), "WARN in [=DEBUG, ERROR]")
	assertTrue(t, level.Enabled(zapcore.ErrorLevel), "ERROR in [=DEBUG, ERROR]")
	assertTrue(t, level.Enabled(zapcore.PanicLevel), "PANIC in [=DEBUG, ERROR]")
}

func TestNormalLevel(t *testing.T) {
	level := log.ParseLevels("ERROR")

	assertFalse(t, level.Enabled(zapcore.InfoLevel), "INFO >=ERROR")
	assertFalse(t, level.Enabled(zapcore.WarnLevel), "WARN >=ERROR")
	assertTrue(t, level.Enabled(zapcore.ErrorLevel), "ERROR >=ERROR")
	assertTrue(t, level.Enabled(zapcore.PanicLevel), "PANIC >=ERROR")
}
