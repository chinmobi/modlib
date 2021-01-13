// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style license that can be found
// at https://github.com/gin-gonic/gin/blob/master/LICENSE

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEngine(t *testing.T) {
	router := newEngine()
	assert.Equal(t, "/", router.basePath)
	assert.Equal(t, router.engine, router)
	assert.Empty(t, router.Handlers)
}

func TestAddRoute(t *testing.T) {
	router := newEngine()
	router.addRoute("TOPIC_1", "/", HandlersChain{func(_ *Context) {}})

	assert.Len(t, router.trees, 1)
	assert.NotNil(t, router.trees.get("TOPIC_1"))
	assert.Nil(t, router.trees.get("TOPIC_2"))

	router.addRoute("TOPIC_2", "/", HandlersChain{func(_ *Context) {}})

	assert.Len(t, router.trees, 2)
	assert.NotNil(t, router.trees.get("TOPIC_1"))
	assert.NotNil(t, router.trees.get("TOPIC_2"))

	router.addRoute("TOPIC_2", "/users", HandlersChain{func(_ *Context) {}})
	assert.Len(t, router.trees, 2)
}
