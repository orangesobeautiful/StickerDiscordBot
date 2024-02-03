package domain

import "errors"

var _ error = (*NotFoundError)(nil)

type NotFoundError struct {
	label string
}

func NewNotFoundError(label string) error {
	return &NotFoundError{label: label}
}

func (e *NotFoundError) Error() string {
	return e.label + " not found"
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	var e *NotFoundError
	return errors.As(err, &e)
}

var _ error = (*AlreadyExistsError)(nil)

type AlreadyExistsError struct {
	label string
}

func NewAlreadyExistsError(label string) error {
	return &AlreadyExistsError{label: label}
}

func (e *AlreadyExistsError) Error() string {
	return e.label + " already exists"
}

func IsAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}

	var e *AlreadyExistsError
	return errors.As(err, &e)
}
