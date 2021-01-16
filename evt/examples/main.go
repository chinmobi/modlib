// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/chinmobi/modlib/evt"
	"github.com/chinmobi/modlib/evt/event"
)

// --- demoMulticaster ---

type demoMulticaster struct{}

func (m demoMulticaster) MulticastEvent(e *event.Event) {
	e.Run()
}

// --- basicListener ---

type basicListener struct {
	Topic, RoutingPath, Source string
	UserID string
	Payload event.Payload
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
	return evt.NewEngine(demoMulticaster{})
}

func tearDown(engine *evt.Engine) {
}

// --- assert helper method ---
func assertEqual(want, got string) int {
	if want != got {
		fmt.Printf("got: %q, want: %q\n", got, want)
		return 1
	}
	return 0
}

// --- demo methods ---

func demoEvtBasicUsage() (errs int) {
	engine := setUp()
	defer tearDown(engine)

	listener := &basicListener{}

	engine.Subscribe("NEW_USER", "/users/:uid", listener)

	publisher := engine.Produce("NEW_USER")

	userInfo := "A test user"

	publisher.PublishEvent("/users/1010", "/test", userInfo)

	errs += assertEqual("NEW_USER", listener.Topic)
	errs += assertEqual("/users/1010", listener.RoutingPath)
	errs += assertEqual("/test", listener.Source)
	errs += assertEqual("1010", listener.UserID)
	errs += assertEqual(userInfo, listener.Payload.(string))

	return
}

func demoEvtReplyUsage() (errs int) {
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

	errs += assertEqual(reply.ReplyTopic, listener.Topic)
	errs += assertEqual("/reply/1010", listener.RoutingPath)
	errs += assertEqual("/users/1010", listener.Source)
	errs += assertEqual("1010", listener.UserID)
	errs += assertEqual(reply.ReplyPayload, listener.Payload.(string))

	return
}

func main() {
	var errs int

	errs += demoEvtBasicUsage()
	errs += demoEvtReplyUsage()

	if errs == 0 {
		fmt.Println("ok")
	}
}
