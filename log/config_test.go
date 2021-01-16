// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log_test

import (
	"testing"

	"github.com/chinmobi/modlib/log"
)

func assertTrue(t *testing.T, value bool, name string) {
	if !value {
		t.Errorf("%s got: false, want: true\n", name)
	}
}

func assertFalse(t *testing.T, value bool, name string) {
	if value {
		t.Errorf("%s got: true, want: false\n", name)
	}
}

func assertEqual(t *testing.T, want, got, name string) {
	if want != got {
		t.Errorf("%s got: %q, want: %q\n", name, got, want)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := log.DefaultConfig()

	assertFalse(t, config.File.Enabled, "fileLogger.Enabled")
	assertTrue(t, config.Console.Enabled, "consoleLogger.Enabled")

	assertEqual(t, log.DEFAULT_LEVEL, config.File.Level,
		"fileLogger.Level")
	assertEqual(t, log.DEBUG_LEVEL, config.Console.Level,
		"consoleLogger.Level")
}
