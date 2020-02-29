package gokoa

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// An ApplicationConfig is a container which stores settings for
// configuring an Application.
//
// The ApplicationConfig is organized as key-value pairs, where the
// value is of limited type.
type ApplicationConfig map[string]interface{}

// A Middleware is a single function, which will be registered into an
// Application to be executed during HTTP request handling.
type Middleware func(ctx *Context, next func() error) error

// A composedHandler is a single function composed by middlewares
// registered into the Application.
type composedHandler func(ctx *Context) error

// An ErrorHandler is a function, which is responsible for handling
// error returned by any middleware.
type ErrorHandler func(err error)

// An Application represents a HTTP web server, which handle HTTP
// requests by processing HTTP responses.
type Application struct {
	// middlewares are functions that will be executed during HTTP
	// request handling.
	middlewares []Middleware

	// errorHandler is the function that handles error returned by
	// middlewares.
	errorHandler ErrorHandler

	// Env is the deploying environment variable, default to GOKOA_ENV
	// or "development".
	Env string

	// Keys is the signed cookie keys, which will be used to sign and
	// verify client cookies, default to an empty array.
	Keys []string

	// Proxy is equal to true when fields in the proxy header will be
	// trusted, default to false.
	Proxy bool

	// SubdomainOffset is the offset of subdomain to be ignored, default
	// to 0 (means no ignoring).
	SubdomainOffset int

	// proxyIpHeader is a key in the proxy header, defaulting to
	// "X-Forwarded-For".
	//
	// proxyIpHeader aims at indicating client's ip address.
	proxyIpHeader string

	// maxIpsCount is the maximum of ip addresses read from the proxy
	// header, default to 0 (means infinity).
	maxIpsCount int
}

var (
	// defaultErrorHandler is the default error handler for the
	// Application.
	defaultErrorHandler = func(err error) {
		log.Println()
		log.Println("gokoa: ", err)
		log.Println()
	}
)

// NewApplication returns a new Application initialized with the given
// config.
//
// The config can be nil, which causes the Application to use default
// configuration settings.
//
// Values in key-value pairs must be in the valid type, otherwise
// NewApplication will panic.
func NewApplication(config ApplicationConfig) *Application {
	app := &Application{
		errorHandler: defaultErrorHandler,
	}

	if config == nil {
		config = make(ApplicationConfig)
	}

	if env, ok := config["env"]; ok {
		app.Env = env.(string)
	} else {
		if env, exist := os.LookupEnv("GOKOA_ENV"); exist {
			app.Env = env
		} else {
			app.Env = "development"
		}
	}

	if keys, ok := config["keys"]; ok {
		app.Keys = keys.([]string)
	}

	if proxy, ok := config["proxy"]; ok {
		app.Proxy = proxy.(bool)
	} else {
		app.Proxy = false
	}

	if subdomainOffset, ok := config["subdomainOffset"]; ok {
		app.SubdomainOffset = subdomainOffset.(int)
	} else {
		app.SubdomainOffset = 2
	}

	if proxyIpHeader, ok := config["proxyIpHeader"]; ok {
		app.proxyIpHeader = proxyIpHeader.(string)
	} else {
		app.proxyIpHeader = "X-Forwarded-For"
	}

	if maxIpsCount, ok := config["maxIpsCount"]; ok {
		app.maxIpsCount = maxIpsCount.(int)
	} else {
		app.maxIpsCount = 0
	}

	return app
}

// Listen causes the Application to create a new HTTP server with a
// single handler which composes all middlewares registered in, and
// listen on the given TCP port for incoming connections.
//
// Listen returns the created http.Server when ListenAndServe() does
// NOT returns an error, otherwise returns it.
func (app *Application) Listen(port int) (*http.Server, error) {
	log.Println("listen")

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(app.Callback()),
	}

	err := server.ListenAndServe()
	if err != nil {
		return nil, err
	}

	return &server, nil
}

// Callback returns a function that composes all middlewares
// registered into the Application, creates a new Context, a new
// Request and a new Response for an incoming connection, and handles
// this HTTP request.
func (app *Application) Callback() func(res http.ResponseWriter, req *http.Request) {
	handler := compose(app.middlewares)

	return func(res http.ResponseWriter, req *http.Request) {
		ctx := app.createContext(res, req)
		app.handleRequest(ctx, handler)
	}
}

// compose composes all middlewares in the Application into a single
// composedHandler and returns it.
func compose(middlewares []Middleware) composedHandler {
	return func(ctx *Context) error {
		index := -1
		var dispatch func(int) error
		dispatch = func(i int) error {
			if i <= index {
				ctx.Response.SetStatus(http.StatusInternalServerError)
				return errors.New("next() called multiple times")
			}
			index = i
			if i == len(middlewares) {
				return nil
			}
			middleware := middlewares[i]
			return middleware(ctx, func() error {
				return dispatch(i + 1)
			})
		}
		return dispatch(0)
	}
}

// createContext creates a new Context, a new Request and a new
// Response and binds them together.
//
// createContext returns the newly created Context.
func (app *Application) createContext(res http.ResponseWriter, req *http.Request) *Context {
	ctx := NewContext()
	ctx.app = app

	request := NewRequest()
	request.Req = req
	request.app = app
	request.ctx = ctx

	response := NewResponse()
	response.Res = res
	response.app = app
	response.ctx = ctx

	request.response = response
	response.request = request

	ctx.Request = request
	ctx.Response = response

	return ctx
}

// handleRequest is responsible for handling HTTP request.
func (app *Application) handleRequest(ctx *Context, handler composedHandler) {
	err := handler(ctx)
	if err != nil {
		app.errorHandler(err)
		return
	}

	app.respond(ctx)
}

// respond is responsible for processing HTTP response and sending it
// to the client.
func (app *Application) respond(ctx *Context) {
	statusCode := ctx.Response.GetStatus()
	body := ctx.Response.GetBody()

	ctx.Response.Res.WriteHeader(statusCode)

	// ignore response body
	if isStatusEmpty(statusCode) {
		ctx.Response.Res.Write(nil)
		return
	}

	// HTTP HEAD request
	if ctx.Request.GetMethod() == "HEAD" {
		if ctx.Response.GetLength() != 0 {
			ctx.Response.SetBody(ctx.Response.GetBody()[0:ctx.Response.GetLength()])
		}
		ctx.Response.Res.Write(body)
		return
	}

	// empty response body
	if body == nil {
		body = []byte(strconv.Itoa(statusCode))
	}

	ctx.Response.Res.Write(body)
}

func isStatusEmpty(statusCode int) bool {
	return statusCode == http.StatusNoContent ||
		statusCode == http.StatusResetContent ||
		statusCode == http.StatusNotModified
}

// Use registers the given middleware into the Application.
//
// Use returns the Application itself, which enables chained function
// call instead of function calls in multiple lines.
func (app *Application) Use(middleware Middleware) *Application {
	log.Println("use middleware")
	app.middlewares = append(app.middlewares, middleware)
	return app
}

// OnError registers a new ErrorHandler into the Application.
func (app *Application) OnError(handler ErrorHandler) {
	app.errorHandler = handler
}
