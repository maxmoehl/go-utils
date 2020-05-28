package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HttpError interface {
	Error() string
	Code() int
	Response(w http.ResponseWriter)
}

type httpError struct {
	error
	code int
}

func (e httpError) Code() int {
	return e.code
}

// Creates a new HttpError that can be sent back
func NewHttpError(code int, message string) HttpError {
	if code < 400 {
		panic("code below 400 does not indicate an error")
	}
	return httpError{
		error: errors.New(message),
		code:  code,
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
		},
	})
	if err != nil {
		LogError(err.Error())
	} else {
		LogWarning(fmt.Sprintf("status %d occured with message %s", e.Code(), e.Error()))
	}
}
