// Copyright 2018 Andy Pan. All rights reserved.
// Use of this source code is governed by a MIT style license that can be found
// at https://github.com/panjf2000/ants/blob/master/LICENSE

package ants

import (
	"time"

	"github.com/chinmobi/modlib/grpool/gr"
)

// goWorker is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type goWorker struct {
	// pool who owns this worker.
	pool *Pool

	// task is a job should be done.
	task chan gr.Runnable

	// recycleTime will be update when putting a worker back into queue.
	recycleTime time.Time
}

// run starts a goroutine to repeat the process
// that performs the function calls.
func (w *goWorker) run() {
	w.pool.incRunning()
	go func() {
		defer func() {
			w.pool.decRunning()
			w.pool.workerCache.Put(w)
			if p := recover(); p != nil {
				if ph := w.pool.options.PanicHandler; ph != nil {
					ph(p)
				} else {
					defaultPanicHandle(w.pool.options, p)
				}
			}
		}()

		for t := range w.task {
			if t == nil {
				return
			}
			t.Run()
			if ok := w.pool.revertWorker(w); !ok {
				return
			}
		}
	}()
}
