package gokoa

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

// A Response represents a HTTP response sent back by the Application.
type Response struct {
	// Res is the primitive HTTP response.
	Res http.ResponseWriter

	request *Request
	ctx *Context
	app *Application

	// statusCode represents the HTTP status code that will be sent back
	// to the client.
	//
	// statusCode is default to 404, meaning that the requested resource
	// is NOT found.
	statusCode int

	// body represents the HTTP response body that will be sent back to
	// the client, default to an empty byte array.
	body []byte
}

// NewResponse returns a new empty Response.
func NewResponse() *Response {
	return &Response{
		statusCode: http.StatusNotFound,
	}
}

// GetStatus returns the HTTP response status code.
func (response *Response) GetStatus() int {
	return response.statusCode
}

// SetStatus assigns the given integer to the HTTP status code.
func (response *Response) SetStatus(statusCode int) {
	// TODO: validate whether the given integer is a valid HTTP status code
	response.statusCode = statusCode
}

// GetBody returns the HTTP response body.
func (response *Response) GetBody() []byte {
	return response.body
}

// SetBody assigns the given object to the HTTP response body.
//
// The given object can be either a string, a byte array, an
// io.Reader, or a map containing key-value pairs. Whenever what type
// the given object is, it will be transformed into a byte array.
func (response *Response) SetBody(body interface{}) error {
	var bytes []byte

	if body == nil {
		bytes = nil

		if !isStatusEmpty(response.statusCode) {
			response.SetStatus(http.StatusNoContent)
		}
		response.Remove("Content-Type")
		response.Remove("Content-Length")
		response.Remove("Transfer-Encoding")
	} else {
		response.SetStatus(http.StatusOK)

		typeSet := response.Has("Content-Type")

		switch body := body.(type) {
		case string:
			bytes = []byte(body)

			if !typeSet {
				if matched, err := regexp.MatchString("^[ \f\n\r\t\v]*<", body); err != nil {
					return err
				} else if matched {
					response.SetType("html")
				} else {
					response.SetType("text")
				}
			}
			response.SetLength(len(bytes))
		case []byte:
			bytes = body

			if !typeSet {
				response.SetType("bin")
			}
			response.SetLength(len(bytes))
		case io.Reader:
			var err error
			bytes, err = ioutil.ReadAll(body)
			if err != nil {
				return err
			}

			if !typeSet {
				response.SetType("bin")
			}
			response.Remove("Content-Length")
		case map[string]interface{}:
			var err error
			bytes, err = json.Marshal(body)
			if err != nil {
				return err
			}

			response.SetType("json")
		}
	}

	response.body = bytes
	return nil
}

// GetLength returns the HTTP response Content-Length header.
func (response *Response) GetLength() int {
	if length, err := strconv.Atoi(response.Get("Content-Length")); err != nil {
		panic(err)
	} else {
		return length
	}
}

// SetLength assigns the given integer to the HTTP response
// Content-Length header.
func (response *Response) SetLength(length int) {
	response.Set("Content-Length", strconv.Itoa(length))
}

// SetType assigns the given string to the HTTP response Content-Type
// header.
func (response *Response) SetType(contentType string) {
	response.Set("Content-Type", contentType)
}

// Get returns the value corresponding to the key in the HTTP response
// header.
func (response *Response) Get(field string) string {
	return response.Res.Header().Get(field)
}

// Has returns true when the specific key-value pair exists in the
// HTTP header.
func (response *Response) Has(field string) bool {
	return response.Get(field) != ""
}

// Set assigns the key-value pair to the HTTP response header.
func (response *Response) Set(field string, value string) {
	response.Res.Header().Set(field, value)
}

// Remove deletes the key-value pair from the HTTP response header.
func (response *Response) Remove(field string) {
	response.Res.Header().Del(field)
}
