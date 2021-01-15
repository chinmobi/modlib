# modlib/grpool

A high-performance and low-cost goroutine pool (forked and customed from Ants) utility for Golang.

## Usage

```go
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

```

## License

The source codes about `ants` are forked and customed from the [ants](https://github.com/panjf2000/ants) project.
All of the related source codes are governed by the [LICENSE](https://github.com/panjf2000/ants/blob/master/LICENSE).
