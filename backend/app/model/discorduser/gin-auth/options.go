package ginauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Option func(*options)

func WithErrRespHandler(h func(*gin.Context, error)) Option {
	return func(o *options) {
		o.errRespHandler = h
	}
}

type options struct {
	errRespHandler func(*gin.Context, error)
}

func newOptions(opts ...Option) *options {
	o := &options{
		errRespHandler: defaultRespHandler,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func defaultRespHandler(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, err)
}
