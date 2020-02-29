package gokoa

import (
	"net/http"
)

// A Request represents a HTTP request received by the Application.
type Request struct {
	// Req is the primitive HTTP request.
	Req *http.Request

	response *Response
	ctx *Context
	app *Application
}

// NewRequest returns a new empty Request.
func NewRequest() *Request {
	return &Request{}
}

// GetMethod returns the HTTP request method.
func (request *Request) GetMethod() string {
	return request.Req.Method
}
