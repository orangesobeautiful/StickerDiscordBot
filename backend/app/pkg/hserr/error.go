package hserr

import (
	"bytes"
	"fmt"

	"golang.org/x/xerrors"
)

var _ ErrResp = (*noWrapErrResp)(nil)

type ErrResp interface {
	error
	xerrors.Formatter
	MarshalJSON() ([]byte, error)

	IsInternal() bool
	HttpStatus() int
	Message() string
	Detail() []string
}

type noWrapErrResp struct {
	message    string
	details    []string
	isInternal bool
	httpStatus int

	err error

	frame xerrors.Frame
}

func (e *noWrapErrResp) IsInternal() bool {
	return e.isInternal
}

func (e *noWrapErrResp) HttpStatus() int {
	return e.httpStatus
}

func (e *noWrapErrResp) Message() string {
	return e.message
}

func (e *noWrapErrResp) Detail() []string {
	return e.details
}

func (e *noWrapErrResp) MarshalJSON() ([]byte, error) {
	bf := bytes.NewBuffer(nil)
	bf.WriteString(`{"Message":"`)
	bf.WriteString(e.message)
	bf.WriteString(`","Detail":[`)
	if len(e.details) > 0 {
		bf.WriteByte('"')
		bf.WriteString(e.details[0])
		bf.WriteByte('"')
	}
	for i := 1; i < len(e.details); i++ {
		bf.WriteString(`,"`)
		bf.WriteString(e.details[i])
		bf.WriteByte('"')
	}
	bf.WriteString(`]}`)
	return bf.Bytes(), nil
}

func (e *noWrapErrResp) Error() string {
	return fmt.Sprint(e)
}

func (e *noWrapErrResp) Format(s fmt.State, v rune) { xerrors.FormatError(e, s, v) }

func (e *noWrapErrResp) FormatError(p xerrors.Printer) error {
	p.Print(e.message)
	if len(e.details) > 0 {
		p.Printf(": %v", e.details[0])
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
