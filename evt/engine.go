// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package evt

import (
	"github.com/chinmobi/modlib/evt/event"
	"github.com/chinmobi/modlib/evt/internal"
)

// --- publisher ---

type publisher struct {
	topic  string
	engine *Engine
}

func newPublisher(topic string, engine *Engine) *publisher {
	p := &publisher{
		topic: topic,
		engine: engine,
	}
	return p
}

func (p *publisher) PublishEvent(routingPath, source string, payload event.Payload) {
	p.engineImpl().PublishEvent(p.topic, routingPath, source, payload)
}

func (p *publisher) engineImpl() *internal.Engine {
	return p.engine.impl
}

// --- Engine ---

type Engine struct {
	impl *internal.Engine
}

func NewEngine(multicaster event.Multicaster) *Engine {
	e := &Engine{
		impl: internal.NewEngine(multicaster),
	}
	return e
}

// --- Broker methods ---

func (e *Engine) Produce(topic string) event.Publisher {
	return newPublisher(topic, e)
}

func (e *Engine) Subscribe(topic, bindingPath string, listener event.Listener) {
	e.engineImpl().BindEventListener(topic, bindingPath, listener)
}

func (e *Engine) Broker() event.Broker {
	return e
}

func (e *Engine) engineImpl() *internal.Engine {
	return e.impl
}
