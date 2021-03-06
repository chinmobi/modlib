// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gr

type Executor interface {
	Execute(task Runnable) error
}

type ExecutorService interface {
	Executor

	Shutdown()
}
