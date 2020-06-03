package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpError interface {
	Error() string
	Code() int
	Response(w http.ResponseWriter)
	Cause() error
}

type httpError struct {
	err
	code int
}

func (e httpError) Code() int {
	return e.code
}

type Error interface {
	Error() string
	Cause() error
}

type err struct {
	message string
	cause   error
}

func (e err) Error() string {
	return e.message
}

func (e err) Cause() error {
	return e.cause
}

// Creates a new HttpError that can be sent back
func NewHttpError(code int, message string, cause error) HttpError {
	if code < 400 {
		panic("code below 400 does not indicate an error")
	}
	return httpError{
		err: err{
			message: message,
			cause:   cause,
		},
		code: code,
	}
}

// Sends back the current error by using the ResponseWriter that is passed in
func (e httpError) Response(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(e.Code())
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    e.Code(),
			"message": e.Error(),
			"cause":   getCause(e.cause),
		},
	})
	if err != nil {
		LogError(err.Error())
	} else {
		LogWarning(fmt.Sprintf("status %d occured with message %s", e.Code(), e.Error()))
	}
}

func getCause(e error) interface{} {
	if e == nil {
		return nil
	}
	if he, ok := e.(HttpError); ok {
		return map[string]interface{}{
			"code":    he.Code(),
			"message": he.Error(),
			"cause":   getCause(he.Cause()),
		}
	} else if err, ok := e.(Error); ok {
		return map[string]interface{}{
			"message": e.Error(),
			"cause":   getCause(err),
		}
	} else {
		return e.Error()
	}
}
