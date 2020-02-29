package gokoa

// A Context contains information related to a single HTTP request.
type Context struct {
	Request *Request
	Response *Response
	app *Application

	// State is the recommended namespace for passing information
	// through different middlewares.
	State map[string]interface{}
}

// NewContext returns a new empty Context.
//
// NewContext allocates memory for State, which means that a key-value
// pair can be directly added to State, without calling make() by
// yourself.
func NewContext() *Context {
	return &Context{
		State: make(map[string]interface{}),
	}
}

func (ctx *Context) GetStatus() int {
	return ctx.Response.GetStatus()
}

func (ctx *Context) SetStatus(statusCode int) {
	ctx.Response.SetStatus(statusCode)
}

func (ctx *Context) GetBody() []byte {
	return ctx.Response.GetBody()
}

func (ctx *Context) SetBody(body interface{}) {
	ctx.Response.SetBody(body)
}
