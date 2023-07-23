package hserr

import (
	"bytes"
	"net/http"
	"strings"
)

var quoteEscaper = strings.NewReplacer(`'`, `\'`, `"`, `\"`)

// ErrInternalServerError internal server error
var (
	ErrInternalServerError = NewErr(http.StatusInternalServerError, "internal server error")
	ErrBadRequest          = NewErr(http.StatusBadRequest, "bad request")
	ErrUnauthorized        = NewErr(http.StatusUnauthorized, "unauthorized")
)

// ErrResp is the error response
type ErrResp struct {
	HttpStatus int
	Code       int
	Message    string
	Detail     []string
}

// NewErr create a new ErrResp
func NewErr(httpStatus int, msg string, details ...string) *ErrResp {
	return &ErrResp{
		HttpStatus: httpStatus,
		Message:    msg,
		Detail:     details,
	}
}

func NewErrWithCode(httpStatus, code int, msg string, details ...string) *ErrResp {
	return &ErrResp{
		HttpStatus: httpStatus,
		Code:       code,
		Message:    msg,
		Detail:     details,
	}
}

// Error implement error interface
func (e *ErrResp) Error() string {
	return e.Message
}

// ToJSONBytes convert to json byte slice
func (e *ErrResp) MarshalJSON() ([]byte, error) {
	bf := bytes.NewBuffer(nil)
	bf.WriteString(`{"Message":"`)
	bf.WriteString(quoteEscaper.Replace(e.Message))
	bf.WriteString(`","Detail":[`)
	if len(e.Detail) > 0 {
		bf.WriteByte('"')
		bf.WriteString(quoteEscaper.Replace(e.Detail[0]))
		bf.WriteByte('"')
	}
	for i := 1; i < len(e.Detail); i++ {
		bf.WriteString(`,"`)
		bf.WriteString(quoteEscaper.Replace(e.Detail[i]))
		bf.WriteByte('"')
	}
	bf.WriteString(`]}`)
	return bf.Bytes(), nil
}
