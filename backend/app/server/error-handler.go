package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend/app/pkg/hserr"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type errHandler struct {
	uniTranslator *ut.UniversalTranslator
}

func newErrHandler(uniTranslator *ut.UniversalTranslator) *errHandler {
	return &errHandler{
		uniTranslator: uniTranslator,
	}
}

func (eh *errHandler) getBindErrConvert(locale string) func(err error) error {
	return func(err error) error {
		var reWriteErr error

		switch realErr := err.(type) {
		case validator.ValidationErrors:
			trans, _ := eh.uniTranslator.GetTranslator(locale)

			detailList := make([]string, 0, len(realErr))
			for _, valErr := range realErr {
				detailList = append(detailList, valErr.Translate(trans))
			}

			reWriteErr = hserr.New(
				http.StatusBadRequest,
				"param of request validate failed",
				hserr.WithDetails(detailList...),
			)
		case *validator.InvalidValidationError:
			reWriteErr = hserr.NewInternalError(err, "validator.InvalidValidationError")
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

		return reWriteErr
	}
}
