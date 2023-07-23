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

// BindHandler
func BindHandler[reqType any, respType any](
	h func(*gin.Context, reqType) (respType, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bindDeal(ctx, false, h)
	}
}

// BindJSONHandler
func BindJSONHandler[reqType any, respType any](
	h func(*gin.Context, reqType) (respType, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bindDeal(ctx, true, h)
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
	ctx *gin.Context, isBindJSON bool,
	h func(*gin.Context, reqType) (respType, error),
) {
	var req reqType
	var err error
	if isBindJSON {
		err = ctx.ShouldBindJSON(&req)
	} else {
		err = ctx.ShouldBind(&req)
	}
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
