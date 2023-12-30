package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"backend/app/pkg/ginext"
	"backend/app/pkg/hserr"
	"backend/app/pkg/log"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/language"
)

func setGinextErrorHandler(uniTranslator *ut.UniversalTranslator) {
	ginext.SetBindErrHandler(func(ctx *gin.Context, err error) {
		var reWriteErr error

		switch realErr := err.(type) {
		case validator.ValidationErrors:
			trans, _ := uniTranslator.GetTranslator(
				ctx.MustGet("lang").(language.Tag).String())

			detailList := make([]string, 0, len(realErr))
			for _, valErr := range realErr {
				detailList = append(detailList, valErr.Translate(trans))
			}

			reWriteErr = hserr.New(
				http.StatusBadRequest,
				"param of request validate failed",
				hserr.WithDetails(detailList...),
			)
		case *json.UnmarshalTypeError:
			reWriteErr = hserr.New(
				http.StatusBadRequest,
				"param of request validate failed",
				hserr.WithDetails(
					realErr.Field+
						" should be "+
						realErr.Type.Name()+
						" not "+realErr.Value,
				),
			)
		case *json.SyntaxError:
			reWriteErr = hserr.New(
				http.StatusBadRequest,
				"decode json body failed",
				hserr.WithDetails(fmt.Sprint(err)),
			)
		default:
			reWriteErr = hserr.New(
				http.StatusBadRequest,
				"bad request format",
				hserr.WithDetails(fmt.Sprint(err)),
			)
		}

		ginHSERROutput(ctx, reWriteErr)
	})

	ginext.SetRespErrHandler(ginHSERROutput)
}

func ginHSERROutput(ctx *gin.Context, err error) {
	var respErr hserr.ErrResp
	if errors.As(err, &respErr) {
		statusCode := respErr.HttpStatus()
		if statusCode >= http.StatusInternalServerError && statusCode < 600 {
			// TODO: use slog
			log.Errorf("respErr=%+v", err)
		}

		ctx.JSON(statusCode, respErr)
		return
	}

	ctx.JSON(http.StatusInternalServerError, hserr.New(
		http.StatusInternalServerError,
		"unknown error",
		hserr.WithDetails(err.Error()),
		hserr.WithIsInternal(),
	))
}
