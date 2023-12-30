package hserr

import (
	"errors"

	"golang.org/x/xerrors"
)

type newOption struct {
	details         []string
	isInternal      bool
	extraCallerSkip uint
	err             error
}

type NewOption func(*newOption)

func WithDetails(details ...string) NewOption {
	return func(o *newOption) {
		o.details = details
	}
}

func WithIsInternal() NewOption {
	return func(o *newOption) {
		o.isInternal = true
	}
}

func WithExtraCallerSkip(extraCallerSkip uint) NewOption {
	return func(o *newOption) {
		o.extraCallerSkip = extraCallerSkip
	}
}

func WithWrapErr(err error) NewOption {
	if errors.Is(err, (*noWrapErrResp)(nil)) || errors.Is(err, (*wrapErrResp)(nil)) {
		panic("err is already a hserr.ErrResp, don't wrap it again")
	}

	return func(o *newOption) {
		o.err = err
	}
}

func New(code int, msg string, opts ...NewOption) error {
	o := &newOption{}
	for _, opt := range opts {
		opt(o)
	}

	callerSkip := 1 + int(o.extraCallerSkip)

	result := noWrapErrResp{
		message:    msg,
		details:    o.details,
		isInternal: o.isInternal,
		httpStatus: code,
		err:        o.err,
		frame:      xerrors.Caller(callerSkip),
	}

	if o.err != nil {
		return &wrapErrResp{
			noWrapErrResp: result,
		}
	}

	return &result
}
