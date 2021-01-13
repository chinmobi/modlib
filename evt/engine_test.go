// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package evt_test

import (
	"testing"

	"github.com/chinmobi/modlib/evt"
	"github.com/chinmobi/modlib/evt/event"

	"github.com/stretchr/testify/assert"
)

// --- testMulticaster ---

type testMulticaster struct{}

func (m testMulticaster) MulticastEvent(e *event.Event) {
	e.Run()
}

// --- basicListener ---

type basicListener struct {
	Topic, RoutingPath, Source string
	UserID string
	Payload event.Payload
}

func (l *basicListener) Reset() {
	l.Topic = ""
	l.RoutingPath = ""
	l.Source = ""
	l.UserID = ""
	l.Payload = nil
}

func (l *basicListener) OnEvent(envelope event.Envelope, payload event.Payload) {
	l.Topic = envelope.Topic()
	l.RoutingPath = envelope.RoutingPath()
	l.Source = envelope.Source()

	l.UserID = envelope.GetParam("uid")

	l.Payload = payload
}

// --- replyListener ---

type replyListener struct{
	ReplyTopic, ReplyPayload string
}

func (l *replyListener) OnEvent(envelope event.Envelope, payload event.Payload) {
	envelope.Reply(l.ReplyTopic, l.ReplyPayload)
}

// --- setUp & tearDown ---

func setUp() *evt.Engine {
	return evt.NewEngine(testMulticaster{})
}

func tearDown(engine *evt.Engine) {
}

// --- Tests ---

func TestEvtBasicUsage(t *testing.T) {
	engine := setUp()
	defer tearDown(engine)

	listener := &basicListener{}

	engine.Subscribe("NEW_USER", "/users/:uid", listener)

	publisher := engine.Produce("NEW_USER")

	userInfo := "A test user"

	publisher.PublishEvent("/users/1010", "/test", userInfo)

	assert := assert.New(t)

	assert.Equal("NEW_USER", listener.Topic)
	assert.Equal("/users/1010", listener.RoutingPath)
	assert.Equal("/test", listener.Source)
	assert.Equal("1010", listener.UserID)
	assert.Equal(userInfo, listener.Payload)
}

func TestEvtReplyUsage(t *testing.T) {
	engine := setUp()
	defer tearDown(engine)

	listener := &basicListener{}
	engine.Subscribe("NEW_USER_REPLY", "/reply/:uid", listener)

	reply := &replyListener{
		ReplyTopic: "NEW_USER_REPLY",
		ReplyPayload: "Handled",
	}
	engine.Subscribe("NEW_USER", "/users/:uid", reply)


	publisher := engine.Produce("NEW_USER")

	publisher.PublishEvent("/users/1010", "/reply/1010", "A test user")

	assert := assert.New(t)

	assert.Equal(reply.ReplyTopic, listener.Topic)
	assert.Equal("/reply/1010", listener.RoutingPath)
	assert.Equal("/users/1010", listener.Source)
	assert.Equal("1010", listener.UserID)
	assert.Equal(reply.ReplyPayload, listener.Payload)
}
