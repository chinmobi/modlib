// Copyright 2020 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

type Payload = interface{}

type Publisher interface {
	PublishEvent(routingPath, source string, payload Payload)
}

type Envelope interface {
	Topic() string

	GetParam(name string) string

	RoutingPath() string

	Source() string

	Reply(topic string, ack Payload)
}

type Listener interface {
	OnEvent(envelope Envelope, payload Payload)
}
