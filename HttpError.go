package logger

import (
	"errors"
)

type HttpError interface {
	Error() string
	Code() int
}

type httpError struct {
	error
	code int
}

func (e httpError) Code() int {
	return e.code
}

func NewHttpError(code int, message string) HttpError {
	if code < 400 {
		panic("code below 400 does not indicate an error")
	}
	return httpError{
		error: errors.New(message),
		code:  code,
	}
}
