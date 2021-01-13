// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"sync"
)

type Engine struct {
	RouterGroup

	pool           sync.Pool
	eventPool      sync.Pool

	trees          topicTreeS
	maxParams      uint16
}

var _ IRouter = &Engine{}

func newEngine() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers:  nil,
			basePath:  "/",
			root:      true,
		},
		trees:       make(topicTreeS, 0, 2),
	}

	engine.RouterGroup.engine = engine

	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}

	return engine
}

func (engine *Engine) allocateContext() *Context {
	return newContext(engine.maxParams)
}

func (engine *Engine) allocateEvent() *Event {
	return NewEvent()
}

// Use attaches a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	return engine
}

func (engine *Engine) addRoute(topic, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(len(handlers) > 0, "there must be at least one handler")

	engine.trees = engine.trees.addRoute(topic, path, handlers)

	// Update maxParams
	if paramsCount := countParams(path); paramsCount > engine.maxParams {
		engine.maxParams = paramsCount
	}
}
