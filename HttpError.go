package logger

import (
	"encoding/json"
	"errors"
	"net/http"
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

func ResponseWithHttpError(w http.ResponseWriter, httpErr httpError) {
	w.WriteHeader(httpErr.Code())
	err := json.NewEncoder(w).Encode(httpErr)
	if err != nil {
		LogError(err.Error())
	}
}
