// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLastChar(t *testing.T) {
	assert.Equal(t, uint8('a'), lastChar("hola"))
	assert.Equal(t, uint8('s'), lastChar("adios"))
	assert.Panics(t, func() { lastChar("") })
}

func TestJoinPaths(t *testing.T) {
	assert.Equal(t, "", joinPaths("", ""))
	assert.Equal(t, "/", joinPaths("", "/"))
	assert.Equal(t, "/a", joinPaths("/a", ""))
	assert.Equal(t, "/a/", joinPaths("/a/", ""))
	assert.Equal(t, "/a/", joinPaths("/a/", "/"))
	assert.Equal(t, "/a/", joinPaths("/a", "/"))
	assert.Equal(t, "/a/hola", joinPaths("/a", "/hola"))
	assert.Equal(t, "/a/hola", joinPaths("/a/", "/hola"))
	assert.Equal(t, "/a/hola/", joinPaths("/a/", "/hola/"))
	assert.Equal(t, "/a/hola/", joinPaths("/a/", "/hola//"))
}
