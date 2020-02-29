package gokoa

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRequest(t *testing.T) {
var request *Request

	// act
	request = NewRequest()

	// assert
	assert.Equal(t, (*http.Request)(nil), request.Req)
	assert.Equal(t, (*Response)(nil), request.response)
	assert.Equal(t, (*Context)(nil), request.ctx)
	assert.Equal(t, (*Application)(nil), request.app)
}

func TestRequest_GetMethod(t *testing.T) {
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux
	var method string

	// arrange
	app = NewApplication(nil)
	app.middlewares = []Middleware {
		func(ctx *Context, next func() error) error {
			method = ctx.Request.GetMethod()
			return nil
		},
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	mux = http.NewServeMux()

	// act
	mux.HandleFunc("/", app.Callback())
	mux.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, http.MethodPost, method)
}
