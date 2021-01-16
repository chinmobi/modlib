// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

// --- EqualOpLevel ---

type equalOpLevel struct {
	lvl zapcore.Level
}

func (op equalOpLevel) Enabled(lvl zapcore.Level) bool {
	return lvl == op.lvl
}

// --- LessOpLevel ---

type lessOpLevel struct {
	lvl zapcore.Level
}

func (op lessOpLevel) Enabled(lvl zapcore.Level) bool {
	return lvl < op.lvl
}

// --- OrOpLevel ---

type orOpLevel struct {
	lvls []zapcore.LevelEnabler
}

func (op *orOpLevel) appendLevel(lvl ...zapcore.LevelEnabler) {
	op.lvls = append(op.lvls, lvl...)
}

func (op *orOpLevel) Enabled(lvl zapcore.Level) bool {
	for i, cnt := 0, len(op.lvls); i < cnt; i++ {
		if op.lvls[i].Enabled(lvl) {
			return true
		}
	}
	return false
}

// --- Parse levels ---

func ParseLevels(levels string) zapcore.LevelEnabler {
	lvls := strings.Split(levels, " ")
	cnt := len(lvls)

	if cnt == 1 {
		return parseSingleLevel(lvls[0])
	}

	op := &orOpLevel{}

	for i := 0; i < cnt; i++ {
		op.appendLevel(parseSingleLevel(lvls[i]))
	}

	return op
}

func parseSingleLevel(str string) zapcore.LevelEnabler {
	switch str[0] {
	case '<':
		return lessOpLevel{lvl: zapLevelOf(str[1:]),}
	case '=':
		return equalOpLevel{lvl: zapLevelOf(str[1:]),}
	default:
		return zapLevelOf(str)
	}
}

func zapLevelOf(str string) zapcore.Level {
	switch str {
	case DEBUG_LEVEL:
		return zapcore.DebugLevel
	case INFO_LEVEL:
		return zapcore.InfoLevel
	case WARN_LEVEL:
		return zapcore.WarnLevel
	case ERROR_LEVEL:
		return zapcore.ErrorLevel
	case PANIC_LEVEL:
		return zapcore.PanicLevel
	case FATAL_LEVEL:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
