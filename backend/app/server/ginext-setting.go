package server

import (
	"errors"
	"net/http"

	"backend/app/pkg/ginext"
	"backend/app/pkg/hserr"
	"backend/app/pkg/log"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func setGinextErrorHandler(eh *errHandler) {
	ginext.SetBindErrHandler(func(ctx *gin.Context, err error) {
		reqLang := ctx.MustGet("lang").(language.Tag).String()

		bindErrConverter := eh.getBindErrConvert(reqLang)

		reWriteErr := bindErrConverter(err)

		ginHSERROutput(ctx, reWriteErr)
	})

	ginext.SetRespErrHandler(ginHSERROutput)
}

func ginHSERROutput(ctx *gin.Context, err error) {
	var respErr hserr.ErrResp
	if errors.As(err, &respErr) {
		statusCode := respErr.HTTPStatus()
		if statusCode >= http.StatusInternalServerError && statusCode < 600 {
			logErrorMessage(err)
		}

		ctx.JSON(statusCode, respErr)
		return
	}

	err = hserr.New(
		http.StatusInternalServerError,
		"unknown error",
		hserr.WithDetails(err.Error()),
		hserr.WithIsInternal(),
	)
	logErrorMessage(err)

	ctx.JSON(http.StatusInternalServerError, err)
}

func logErrorMessage(err error) {
	// TODO: use slog
	log.Errorf("respErr=%+v", err)
}
