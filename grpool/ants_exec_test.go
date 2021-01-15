// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package grpool_test

import (
	"sync"
	"sync/atomic"
	"testing"
	//"time"

	"github.com/chinmobi/modlib/grpool"

	"github.com/stretchr/testify/assert"
)

func demoFunc(sum *int32) {
	//time.Sleep(time.Duration(10) * time.Millisecond)
	atomic.AddInt32(sum, 1)
}

type demoRunnable struct {
	Func func()
}

func (r demoRunnable) Run() {
	r.Func()
}

func TestAntsExecutorUsage(t *testing.T) {
	var wg sync.WaitGroup

	exec, err := grpool.NewAntsExecutor(10, grpool.DefaultAntsOptions())
	assert.NoError(t, err)

	defer exec.Shutdown()

	var errs error

	var sum, count int32 = 0, 15

	for i := 0; i < int(count); i++ {
		wg.Add(1)
		err := exec.Execute(demoRunnable{Func: func() {
			demoFunc(&sum)
			wg.Done()
		},})
		if err != nil {
			errs = err
		}
	}
	wg.Wait()

	assert.NoError(t, errs)
	assert.Equal(t, count, sum)
}
