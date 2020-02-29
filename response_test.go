package gokoa

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResponse(t *testing.T) {
	var response *Response

	// act
	response = NewResponse()

	// assert
	assert.Equal(t, http.ResponseWriter(nil), response.Res)
	assert.Equal(t, (*Request)(nil), response.request)
	assert.Equal(t, (*Context)(nil), response.ctx)
	assert.Equal(t, (*Application)(nil), response.app)
	assert.Equal(t, http.StatusNotFound, response.statusCode)
	assert.Equal(t, []byte(nil), response.body)
}

func TestResponse_GetStatus(t *testing.T) {
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux
	var statusCode int

	// arrange
	app = NewApplication(nil)
	app.middlewares = []Middleware {
		func(ctx *Context, next func() error) error {
			statusCode = ctx.Response.GetStatus()
			return nil
		},
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	mux = http.NewServeMux()

	// act
	mux.HandleFunc("/", app.Callback())
	mux.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, http.StatusNotFound, statusCode)
}

func TestResponse_SetStatus(t *testing.T) {
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux
	var res *http.Response

	// arrange
	app = NewApplication(nil)
	app.middlewares = []Middleware {
		func(ctx *Context, next func() error) error {
			ctx.Response.SetStatus(400)
			return nil
		},
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	mux = http.NewServeMux()

	// act
	mux.HandleFunc("/", app.Callback())
	mux.ServeHTTP(rec, req)
	res = rec.Result()

	// assert
	assert.Equal(t, 400, res.StatusCode)
}

func TestResponse_GetBody(t *testing.T) {
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux
	var body []byte

	// arrange
	app = NewApplication(nil)
	app.middlewares = []Middleware {
		func(ctx *Context, next func() error) error {
			body = ctx.Response.GetBody()
			return nil
		},
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	mux = http.NewServeMux()

	// act
	mux.HandleFunc("/", app.Callback())
	mux.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, []byte(nil), body)
}
