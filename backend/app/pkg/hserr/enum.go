package hserr

import (
	"net/http"
)

var (
	ErrBadRequest   = New(http.StatusBadRequest, "bad request")
	ErrUnauthorized = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden    = New(http.StatusForbidden, "forbidden")
)

func NewInternalError(err error, message string) error {
	return New(
		http.StatusInternalServerError,
		"internal server error",
		WithWrapErr(err),
		WithIsInternal(),
		WithExtraCallerSkip(1),
	)
}
