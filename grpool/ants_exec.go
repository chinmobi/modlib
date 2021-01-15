// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package grpool

import (
	"github.com/chinmobi/modlib/grpool/ants"
	"github.com/chinmobi/modlib/grpool/gr"
)

type antsExecutor struct {
	pool *ants.Pool
}

func (exec *antsExecutor) Execute(task gr.Runnable) error {
	return exec.pool.Submit(task)
}

func (exec *antsExecutor) Shutdown() {
	exec.pool.Release()
}

// Create an ExecutorService based on ants.Pool with the poolSize and options.
func NewAntsExecutor(poolSize int, opts *ants.Options) (gr.ExecutorService, error) {
	pool, err := ants.NewPool(poolSize, ants.WithOptions(opts))
	if err != nil {
		return nil, err
	}

	return &antsExecutor{pool: pool,}, nil
}

// Get the default ants options
func DefaultAntsOptions() *ants.Options {
	opts := &ants.Options{
		ExpiryDuration: ants.DefaultCleanIntervalTime,
		PreAlloc: true,
		MaxBlockingTasks: 0,
		Nonblocking: false,
	}
	return opts
}
