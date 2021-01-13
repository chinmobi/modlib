// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"github.com/chinmobi/modlib/evt/event"
)

type EventPayload = event.Payload

type EventListener = event.Listener

// --- EventHandler ---

type EventHandler = event.Handler

// --- Event ---

type Event = event.Event

func NewEvent() *Event {
	return event.NewEvent()
}
