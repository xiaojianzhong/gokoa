package gokoa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContext(t *testing.T) {
	var ctx *Context

	// act
	ctx = NewContext()

	// assert
	assert.Equal(t, (*Request)(nil), ctx.Request)
	assert.Equal(t, (*Response)(nil), ctx.Response)
	assert.Equal(t, (*Application)(nil), ctx.app)
	assert.Equal(t, make(map[string]interface{}), ctx.State)
}
