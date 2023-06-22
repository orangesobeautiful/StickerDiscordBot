package hs

import (
	"encoding/json"
	"net/http"
	"sync"
)

var ctxPool = sync.Pool{
	New: func() any {
		return new(Context)
	},
}

// Context is the context of a http request
type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// newContext return a new Context
func newContext(r *http.Request, w http.ResponseWriter) *Context {
	newCtx := ctxPool.Get().(*Context)
	newCtx.Request = r
	newCtx.ResponseWriter = w

	return newCtx
}

// putContext put a Context back to pool
func putContext(ctx *Context) {
	ctxPool.Put(ctx)
}

// StatusCode set the status code of response
func (ctx *Context) StatusCode(statusCode int) {
	ctx.ResponseWriter.WriteHeader(statusCode)
}

// setJSONHeader set the header of response to json
func (ctx *Context) setJSONHeader() {
	w := ctx.ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// writeJSON write data to response body as json
func (ctx *Context) writeJSON(data any) error {
	ctx.setJSONHeader()
	w := ctx.ResponseWriter
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}
