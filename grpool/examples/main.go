// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/chinmobi/modlib/grpool"
)

func demoFunc(sum *int32) {
	atomic.AddInt32(sum, 1)
}

type demoRunnable struct {
	Func func()
}

func (r demoRunnable) Run() {
	r.Func()
}

func main() {
	var wg sync.WaitGroup

	exec, err := grpool.NewAntsExecutor(10, grpool.DefaultAntsOptions())
	if err != nil {
		fmt.Errorf("%+v\n", err)
		return
	}
	defer exec.Shutdown()

	var sum, count int32 = 0, 15

	for i := 0; i < int(count); i++ {
		wg.Add(1)
		err = exec.Execute(demoRunnable{Func: func() {
			demoFunc(&sum)
			wg.Done()
		},})
		if err != nil {
			wg.Done()
			break
		}
	}
	wg.Wait()

	if err != nil {
		fmt.Errorf("%+v\n", err)
	} else if count != sum {
		fmt.Errorf("got: %d, want: %d\n", sum, count)
	} else {
		fmt.Println("ok")
	}
}
