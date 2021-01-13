// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

// --- Handler ---

type Handler interface {
	HandleEvent(event *Event)
}

// --- Multicaster ---

type Multicaster interface {
	MulticastEvent(event *Event)
}

// --- Event ---

type Event struct {
	Topic, RoutingPath, Source string
	Payload Payload
	Handler Handler
}

func NewEvent() *Event {
	event := &Event{
	}
	return event
}

func (e *Event) Init(topic, routingPath, source string, payload Payload) *Event {
	e.Topic = topic
	e.RoutingPath = routingPath
	e.Source = source
	e.Payload = payload

	return e
}

func (e *Event) Reset() *Event {
	e.Topic = ""
	e.RoutingPath = ""
	e.Source = ""
	e.Payload = nil

	e.Handler = nil

	return e
}

func (e *Event) Run() {
	if e.Handler != nil {
		e.Handler.HandleEvent(e)
	}
}
