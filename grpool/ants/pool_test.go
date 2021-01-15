// Copyright 2018 Andy Pan. All rights reserved.
// Use of this source code is governed by a MIT style license that can be found
// at https://github.com/panjf2000/ants/blob/master/LICENSE

package ants

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func demoFunc() {
	time.Sleep(time.Duration(10) * time.Millisecond)
}

func demoPoolFunc(args interface{}) {
	n := args.(int)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func longRunningFunc() {
	for {
		runtime.Gosched()
	}
}

type testRunnable struct {
	Func func()
}

func (r testRunnable) Run() {
	r.Func()
}

// -----------------------------------------------------------------------------

func TestPanicHandler(t *testing.T) {
	var panicCounter int64
	var wg sync.WaitGroup
	p0, err := NewPool(10, WithPanicHandler(func(p interface{}) {
		defer wg.Done()
		atomic.AddInt64(&panicCounter, 1)
		t.Logf("catch panic with PanicHandler: %v", p)
	}))
	assert.NoErrorf(t, err, "create new pool failed: %v", err)
	defer p0.Release()
	wg.Add(1)
	_ = p0.Submit(testRunnable{Func: func() {
		panic("Oops!")
	},})
	wg.Wait()
	c := atomic.LoadInt64(&panicCounter)
	assert.EqualValuesf(t, 1, c, "panic handler didn't work, panicCounter: %d", c)
	assert.EqualValues(t, 0, p0.Running(), "pool should be empty after panic")
}

func TestPanicHandlerPreMalloc(t *testing.T) {
	var panicCounter int64
	var wg sync.WaitGroup
	p0, err := NewPool(10, WithPreAlloc(true), WithPanicHandler(func(p interface{}) {
		defer wg.Done()
		atomic.AddInt64(&panicCounter, 1)
		t.Logf("catch panic with PanicHandler: %v", p)
	}))
	assert.NoErrorf(t, err, "create new pool failed: %v", err)
	defer p0.Release()
	wg.Add(1)
	_ = p0.Submit(testRunnable{Func: func() {
		panic("Oops!")
	},})
	wg.Wait()
	c := atomic.LoadInt64(&panicCounter)
	assert.EqualValuesf(t, 1, c, "panic handler didn't work, panicCounter: %d", c)
	assert.EqualValues(t, 0, p0.Running(), "pool should be empty after panic")
}

// -----------------------------------------------------------------------------

func TestPurge(t *testing.T) {
	p, err := NewPool(10)
	assert.NoErrorf(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	_ = p.Submit(testRunnable{Func: demoFunc,})
	time.Sleep(3 * DefaultCleanIntervalTime)
	assert.EqualValues(t, 0, p.Running(), "all p should be purged")
}

func TestPurgePreMalloc(t *testing.T) {
	p, err := NewPool(10, WithPreAlloc(true))
	assert.NoErrorf(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	_ = p.Submit(testRunnable{Func: demoFunc,})
	time.Sleep(3 * DefaultCleanIntervalTime)
	assert.EqualValues(t, 0, p.Running(), "all p should be purged")
}

// -----------------------------------------------------------------------------

func TestNonblockingSubmit(t *testing.T) {
	poolSize := 10
	p, err := NewPool(poolSize, WithNonblocking(true))
	assert.NoErrorf(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	for i := 0; i < poolSize-1; i++ {
		assert.NoError(t,
			p.Submit(testRunnable{Func: longRunningFunc,}),
			"nonblocking submit when pool is not full shouldn't return error")
	}
	ch := make(chan struct{})
	ch1 := make(chan struct{})
	f := func() {
		<-ch
		close(ch1)
	}
	// p is full now.
	assert.NoError(t,
		p.Submit(testRunnable{Func: f,}),
		"nonblocking submit when pool is not full shouldn't return error")
	assert.EqualError(t,
		p.Submit(testRunnable{Func: demoFunc,}), ErrPoolOverload.Error(),
		"nonblocking submit when pool is full should get an ErrPoolOverload")
	// interrupt f to get an available worker
	close(ch)
	<-ch1
	assert.NoError(t,
		p.Submit(testRunnable{Func: demoFunc,}),
		"nonblocking submit when pool is not full shouldn't return error")
}

// -----------------------------------------------------------------------------

func TestMaxBlockingSubmit(t *testing.T) {
	poolSize := 10
	p, err := NewPool(poolSize, WithMaxBlockingTasks(1))
	assert.NoErrorf(t, err, "create TimingPool failed: %v", err)
	defer p.Release()
	for i := 0; i < poolSize-1; i++ {
		assert.NoError(t,
			p.Submit(testRunnable{Func: longRunningFunc,}),
			"submit when pool is not full shouldn't return error")
	}
	ch := make(chan struct{})
	f := func() {
		<-ch
	}
	// p is full now.
	assert.NoError(t,
		p.Submit(testRunnable{Func: f,}),
		"submit when pool is not full shouldn't return error")
	var wg sync.WaitGroup
	wg.Add(1)
	errCh := make(chan error, 1)
	go func() {
		// should be blocked. blocking num == 1
		if err := p.Submit(testRunnable{Func: demoFunc,}); err != nil {
			errCh <- err
		}
		wg.Done()
	}()
	time.Sleep(1 * time.Second)
	// already reached max blocking limit
	assert.EqualError(t,
		p.Submit(testRunnable{Func: demoFunc,}), ErrPoolOverload.Error(),
		"blocking submit when pool reach max blocking submit should return ErrPoolOverload")
	// interrupt f to make blocking submit successful.
	close(ch)
	wg.Wait()
	select {
	case <-errCh:
		t.Fatalf("blocking submit when pool is full should not return error")
	default:
	}
}

// -----------------------------------------------------------------------------

func TestRebootPool(t *testing.T) {
	var wg sync.WaitGroup
	p, err := NewPool(10)
	assert.NoErrorf(t, err, "create Pool failed: %v", err)
	defer p.Release()
	wg.Add(1)
	_ = p.Submit(testRunnable{Func: func() {
		demoFunc()
		wg.Done()
	},})
	wg.Wait()
	p.Release()
	assert.EqualError(t, p.Submit(nil), ErrPoolClosed.Error(), "pool should be closed")
	p.Reboot()
	wg.Add(1)
	assert.NoError(t,
		p.Submit(testRunnable{Func: func() { wg.Done() },}),
		"pool should be rebooted")
	wg.Wait()
}

// -----------------------------------------------------------------------------

func TestInfinitePool(t *testing.T) {
	c := make(chan struct{})
	p, _ := NewPool(-1)
	_ = p.Submit(testRunnable{Func: func() {
		_ = p.Submit(testRunnable{Func: func() {
			<-c
		},})
	},})
	c <- struct{}{}
	if n := p.Running(); n != 2 {
		t.Errorf("expect 2 workers running, but got %d", n)
	}
	p.Tune(10)
	if capacity := p.Cap(); capacity != -1 {
		t.Fatalf("expect capacity: -1 but got %d", capacity)
	}
	var err error
	p, err = NewPool(-1, WithPreAlloc(true))
	if err != ErrInvalidPreAllocSize {
		t.Errorf("expect ErrInvalidPreAllocSize but got %v", err)
	}
}
