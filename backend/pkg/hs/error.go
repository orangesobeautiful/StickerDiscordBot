package hs

import (
	"bytes"
	"net/http"
	"strings"
)

var quoteEscaper = strings.NewReplacer(`'`, `\'`, `"`, `\"`)

var (
	// ErrInternalServerError internal server error
	ErrInternalServerError = NewErr(http.StatusInternalServerError, "internal server error")
)

// ErrResp is the error response
type ErrResp struct {
	Status  int
	Message string
	Detail  []string
}

// NewErr create a new ErrResp
func NewErr(status int, msg string, details ...string) *ErrResp {
	return &ErrResp{
		Status:  status,
		Message: msg,
		Detail:  details,
	}
}

// Error implement error interface
func (e ErrResp) Error() string {
	return e.Message
}

// ToJSONBytes convert to json byte slice
func (e ErrResp) ToJSONBytes() []byte {
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
	return bf.Bytes()
}
