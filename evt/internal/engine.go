// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style license that can be found
// at https://github.com/gin-gonic/gin/blob/master/LICENSE

// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"sync"
)

type Engine struct {
	RouterGroup

	multicaster    EventMulticaster

	pool           sync.Pool
	eventPool      sync.Pool

	trees          topicTreeS
	maxParams      uint16
}

var _ IRouter = &Engine{}

func NewEngine(multicaster EventMulticaster) *Engine {
	engine := newEngine()

	engine.multicaster = multicaster

	engine.eventPool.New = func() interface{} {
		return engine.allocateEvent()
	}

	return engine
}

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

// --- Event handle, publish, ... ---

func (engine *Engine) HandleEvent(event *Event) {
	c := engine.pool.Get().(*Context)

	engine.trees.handle(c, event)

	engine.pool.Put(c)
}

func (engine *Engine) PublishEvent(topic, routingPath, source string, payload EventPayload) {
	event := engine.eventPool.Get().(*Event)

	event.Init(topic, routingPath, source, payload)
	event.Handler = engine

	engine.multicaster.MulticastEvent(event)
}

func (engine *Engine) ReplyEvent(event *Event, ack EventPayload) {
	engine.PublishEvent(event.Topic, event.Source, event.RoutingPath, ack)
}
