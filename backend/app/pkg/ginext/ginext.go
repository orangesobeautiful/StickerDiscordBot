package ginext

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var bindErrHandler, respErrHandler func(*gin.Context, error)

func SetBindErrHandler(h func(*gin.Context, error)) {
	bindErrHandler = h
}

func SetRespErrHandler(h func(*gin.Context, error)) {
	respErrHandler = h
}

type bindType int

const (
	bindTypeShouldBind bindType = iota + 1
	bindTypeShouldBindURI
)

func (b bindType) GetShouldBindFunc(ctx *gin.Context) func(obj any) error {
	switch b {
	case bindTypeShouldBind:
		return ctx.ShouldBind
	case bindTypeShouldBindURI:
		return ctx.ShouldBindUri
	}

	return nil
}

// BindHandler
func BindHandler[reqType any, respType any](
	h func(*gin.Context, reqType) (respType, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bindDeal(ctx, bindTypeShouldBind, h)
	}
}

// BindURIHandler
func BindURIHandler[reqType any, respType any](
	h func(*gin.Context, reqType) (respType, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bindDeal(ctx, bindTypeShouldBindURI, h)
	}
}

// Handler
func Handler[respType any](h func(*gin.Context) (respType, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := h(ctx)
		respDeal(ctx, resp, err)
	}
}

func bindDeal[reqType any, respType any](
	ctx *gin.Context, bindType bindType,
	h func(*gin.Context, reqType) (respType, error),
) {
	var req reqType
	var err error
	err = bindType.GetShouldBindFunc(ctx)(&req)
	if err != nil {
		if bindErrHandler != nil {
			bindErrHandler(ctx, err)
		}
		return
	}

	resp, err := h(ctx, req)
	respDeal(ctx, resp, err)
}

func respDeal(ctx *gin.Context, resp any, err error) {
	if err != nil {
		if respErrHandler != nil {
			respErrHandler(ctx, err)
		}
		return
	}

	if isEmptyResp(resp) {
		ctx.Status(http.StatusOK)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
