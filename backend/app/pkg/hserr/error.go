package hserr

import (
	"fmt"

	"golang.org/x/xerrors"
)

var _ ErrResp = (*noWrapErrResp)(nil)

type ErrResp interface {
	error
	xerrors.Formatter

	IsInternal() bool
	HTTPStatus() int
	GetMessage() string
	GetDetails() []string
}

type noWrapErrResp struct {
	Message    string
	Details    []string
	isInternal bool
	httpStatus int

	err error

	frame xerrors.Frame
}

func (e *noWrapErrResp) IsInternal() bool {
	return e.isInternal
}

func (e *noWrapErrResp) HTTPStatus() int {
	return e.httpStatus
}

func (e *noWrapErrResp) GetMessage() string {
	return e.Message
}

func (e *noWrapErrResp) GetDetails() []string {
	return e.Details
}

func (e *noWrapErrResp) Error() string {
	return fmt.Sprint(e)
}

func (e *noWrapErrResp) Format(s fmt.State, v rune) { xerrors.FormatError(e, s, v) }

func (e *noWrapErrResp) FormatError(p xerrors.Printer) error {
	p.Print(e.Message)
	if len(e.Details) > 0 {
		p.Printf(": %v", e.Details[0])
	}

	e.frame.Format(p)

	return e.err
}

var (
	_ ErrResp         = (*noWrapErrResp)(nil)
	_ xerrors.Wrapper = (*wrapErrResp)(nil)
)

type wrapErrResp struct {
	noWrapErrResp
}

func (e *wrapErrResp) Unwrap() error {
	return e.err
}
