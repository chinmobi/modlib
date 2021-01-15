// Copyright 2018 Andy Pan. All rights reserved.
// Use of this source code is governed by a MIT style license that can be found
// at https://github.com/panjf2000/ants/blob/master/LICENSE

package ants

import (
	"log"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
)

const (
	Param    = 100
	AntsSize = 1000
	TestSize = 10000
	n        = 10
)

var curMem uint64

// -----------------------------------------------------------------------------

// TestAntsPoolWaitToGetWorker is used to test waiting to get worker.
func TestAntsPoolWaitToGetWorker(t *testing.T) {
	var wg sync.WaitGroup
	p, _ := NewPool(AntsSize)
	defer p.Release()

	for i := 0; i < n; i++ {
		wg.Add(1)
		_ = p.Submit(testRunnable{Func: func() {
			demoPoolFunc(Param)
			wg.Done()
		},})
	}
	wg.Wait()
	t.Logf("pool, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestAntsPoolWaitToGetWorkerPreMalloc(t *testing.T) {
	var wg sync.WaitGroup
	p, _ := NewPool(AntsSize, WithPreAlloc(true))
	defer p.Release()

	for i := 0; i < n; i++ {
		wg.Add(1)
		_ = p.Submit(testRunnable{Func: func() {
			demoPoolFunc(Param)
			wg.Done()
		},})
	}
	wg.Wait()
	t.Logf("pool, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

// -----------------------------------------------------------------------------

// TestAntsPoolGetWorkerFromCache is used to test getting worker from sync.Pool.
func TestAntsPoolGetWorkerFromCache(t *testing.T) {
	p, _ := NewPool(TestSize)
	defer p.Release()

	for i := 0; i < AntsSize; i++ {
		_ = p.Submit(testRunnable{Func: demoFunc,})
	}
	time.Sleep(2 * DefaultCleanIntervalTime)
	_ = p.Submit(testRunnable{Func: demoFunc,})
	t.Logf("pool, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

// -----------------------------------------------------------------------------

func TestRestCodeCoverage(t *testing.T) {
	_, err := NewPool(-1, WithExpiryDuration(-1))
	t.Log(err)
	_, err = NewPool(1, WithExpiryDuration(-1))
	t.Log(err)

	options := Options{}
	options.ExpiryDuration = time.Duration(10) * time.Second
	options.Nonblocking = true
	options.PreAlloc = true
	poolOpts, _ := NewPool(1, WithOptions(options))
	t.Logf("Pool with options, capacity: %d", poolOpts.Cap())

	p0, _ := NewPool(TestSize, WithLogger(log.New(os.Stderr, "", log.LstdFlags)))
	defer func() {
		_ = p0.Submit(testRunnable{Func: demoFunc,})
	}()
	defer p0.Release()
	for i := 0; i < n; i++ {
		_ = p0.Submit(testRunnable{Func: demoFunc,})
	}
	t.Logf("pool, capacity:%d", p0.Cap())
	t.Logf("pool, running workers number:%d", p0.Running())
	t.Logf("pool, free workers number:%d", p0.Free())
	p0.Tune(TestSize)
	p0.Tune(TestSize / 10)
	t.Logf("pool, after tuning capacity, capacity:%d, running:%d", p0.Cap(), p0.Running())

	pprem, _ := NewPool(TestSize, WithPreAlloc(true))
	defer func() {
		_ = pprem.Submit(testRunnable{Func: demoFunc,})
	}()
	defer pprem.Release()
	for i := 0; i < n; i++ {
		_ = pprem.Submit(testRunnable{Func: demoFunc,})
	}
	t.Logf("pre-malloc pool, capacity:%d", pprem.Cap())
	t.Logf("pre-malloc pool, running workers number:%d", pprem.Running())
	t.Logf("pre-malloc pool, free workers number:%d", pprem.Free())
	pprem.Tune(TestSize)
	pprem.Tune(TestSize / 10)
	t.Logf("pre-malloc pool, after tuning capacity, capacity:%d, running:%d",
		pprem.Cap(), pprem.Running())
}
