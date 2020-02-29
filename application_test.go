package gokoa

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestNewApplication_WithoutConfig(t *testing.T) {
	var app *Application

	// arrange
	os.Unsetenv("GOKOA_ENV")

	// act
	app = NewApplication(nil)

	// assert
	assert.Equal(t, "development", app.Env)
	assert.Equal(t, []string(nil), app.Keys)
	assert.Equal(t, false, app.Proxy)
	assert.Equal(t, 2, app.SubdomainOffset)
	assert.Equal(t, "X-Forwarded-For", app.proxyIpHeader)
	assert.Equal(t, 0, app.maxIpsCount)

	// arrange
	os.Setenv("GOKOA_ENV", "test")

	// act
	app = NewApplication(nil)

	// assert
	assert.Equal(t, "test", app.Env)
}

func TestNewApplication_WithConfig(t *testing.T) {
	var config ApplicationConfig
	var app *Application

	// arrange
	config = make(ApplicationConfig)
	config["env"] = "test"
	config["keys"] = []string{ "1", "2", "3" }
	config["proxy"] = true
	config["subdomainOffset"] = 5
	config["proxyIpHeader"] = ""
	config["maxIpsCount"] = 10

	// act
	app = NewApplication(config)

	// assert
	assert.Equal(t, "test", app.Env)
	assert.Equal(t, []string{ "1", "2", "3" }, app.Keys)
	assert.Equal(t, true, app.Proxy)
	assert.Equal(t, 5, app.SubdomainOffset)
	assert.Equal(t, "", app.proxyIpHeader)
	assert.Equal(t, 10, app.maxIpsCount)
}

func TestApplication_Listen(t *testing.T) {
	var app *Application
	var res *http.Response
	var err error

	// arrange
	app = NewApplication(nil)

	// act
	go app.Listen(34567)
	// TODO: find a more elegant way to test a blocking function
	for {
		time.Sleep(1 * time.Second)
		res, err = http.Get("http://localhost:34567/")
		if err == nil {
			break
		}
	}

	// assert
	assert.NotNil(t, res)
}

func TestApplication_Callback_WithoutMiddleware(t *testing.T) {
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux
	var res *http.Response

	// arrange
	app = NewApplication(nil)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	mux = http.NewServeMux()

	// act
	mux.HandleFunc("/", app.Callback())
	mux.ServeHTTP(rec, req)
	res = rec.Result()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	defer res.Body.Close()

	// assert
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, []byte(strconv.Itoa(http.StatusNotFound)), body)
}

func TestApplication_Callback_WithMiddleware(t *testing.T) {
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux
	var res *http.Response
	var body []byte
	var err error

	// arrange
	app = NewApplication(nil)
	app.middlewares = []Middleware {
		func(ctx *Context, next func() error) error {
			ctx.SetBody("response body 1")
			return next()
		},
		func(ctx *Context, next func() error) error {
			ctx.SetBody(string(ctx.GetBody()) + " " + "response body 2")
			return nil
		},
		func(ctx *Context, next func() error) error {
			ctx.SetBody(string(ctx.GetBody()) + " " + "response body 3")
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
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	defer res.Body.Close()

	// assert
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, []byte("response body 1 response body 2"), body)
}

func TestApplication_Callback_MultipleNext(t *testing.T) {
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux
	var res *http.Response
	var body []byte
	var err error

	// arrange
	app = NewApplication(nil)
	app.middlewares = []Middleware {
		func(ctx *Context, next func() error) error {
			next()
			next()
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
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	defer res.Body.Close()

	// assert
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, []byte(strconv.Itoa(http.StatusInternalServerError)), body)
}

func TestApplication_Callback_IgnoringBody(t *testing.T) {
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux
	var res *http.Response
	var body []byte
	var err error

	emptyStatusCodes := []int{
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusNotModified,
	}

	for _, code := range emptyStatusCodes {
		// arrange
		app = NewApplication(nil)
		app.middlewares = []Middleware {
			func(ctx *Context, next func() error) error {
				ctx.SetStatus(code)
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
		body, err = ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		defer res.Body.Close()

		// assert
		assert.Empty(t, body)
	}
}

func TestApplication_Callback_EmptyBody(t *testing.T) {
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux
	var res *http.Response
	var body []byte
	var err error

	// arrange
	app = NewApplication(nil)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.ProtoMajor = 2
	mux = http.NewServeMux()

	// act
	mux.HandleFunc("/", app.Callback())
	mux.ServeHTTP(rec, req)
	res = rec.Result()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	defer res.Body.Close()

	// assert
	assert.Equal(t, []byte(strconv.Itoa(http.StatusNotFound)), body)
}

func TestApplication_Use(t *testing.T) {
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux

	var calledTimes int

	// arrange
	calledTimes = 0
	app = NewApplication(nil)
	middleware := func(ctx *Context, next func() error) error {
		calledTimes++
		return next()
	}

	// act
	app.Use(middleware).Use(middleware).Use(middleware)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	mux = http.NewServeMux()
	mux.HandleFunc("/", app.Callback())
	mux.ServeHTTP(rec, req)

	// assert
	assert.Equal(t, 3, calledTimes)
}

func TestApplication_OnError(t *testing.T) {
	var called bool
	var errorMessage string
	var app *Application
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var mux *http.ServeMux
	var res *http.Response
	var err error

	// arrange
	called = false
	errorMessage = ""
	app = NewApplication(nil)
	app.middlewares = []Middleware {
		func(ctx *Context, next func() error) error {
			return errors.New("error message")
		},
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	mux = http.NewServeMux()

	// act
	app.OnError(func(err error) {
		called = true
		errorMessage = err.Error()
	})
	mux.HandleFunc("/", app.Callback())
	mux.ServeHTTP(rec, req)
	res = rec.Result()
	_, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	defer res.Body.Close()

	// assert
	assert.Equal(t, true, called)
	assert.Equal(t, "error message", errorMessage)
}
