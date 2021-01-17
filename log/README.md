# modlib/log

A fast, leveled, configurable logging (based on [zap](https://pkg.go.dev/go.uber.org/zap)) in Go.

## Usage

```go
package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/chinmobi/modlib/log"
)

func setUp() {
	config := log.DefaultConfig()

	tempDir, err := ioutil.TempDir("", "*-logs")
	if err == nil {
		config.File.Enabled = true
		config.File.Filename = filepath.Join(tempDir, "info.log")
		config.File.Level = "<WARN"

		warnFilename := filepath.Join(tempDir, "warn.log")
		warnFileCfg := log.NewFileConfig("=WARN", warnFilename)
		config.AppendFileConfig(warnFileCfg)

		errFilename := filepath.Join(tempDir, "error.log")
		errFileCfg := log.NewFileConfig("ERROR", errFilename)
		config.AppendFileConfig(errFileCfg)
	} else {
		fmt.Errorf("%+v", err)
	}

	log.SetUpLogger(config)
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			// discard
		}
	}()

	setUp()

	log.Debug("A debug msg")
	log.Debugf("A %s msg", "debug")

	log.Info("A info msg")
	log.Infof("A %s msg", "info")

	log.Warn("A warn msg")
	log.Warnf("A %s msg", "warn")

	log.Error("A error msg")
	log.Errorf("A %s msg", "error")

	log.Panic("A panic msg")
}

```
