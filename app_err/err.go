package app_err

import (
	"net/http"
)

type Error struct {
	Code    string
	Status  int
	Message string
	Err     error
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return e.Message + ":" + e.Err.Error()
}

func BadRequest(code string, message string, err error) *Error {
	return &Error{Code: code,
		Status:  400,
		Message: message,
		Err:     err,
	}
}

func Unauthorized(code string, message string, err error) *Error {
	return &Error{Code: code,
		Status:  401,
		Message: message,
		Err:     err,
	}
}

func Forbidden(code string, message string, err error) *Error {
	return &Error{Code: code,
		Status:  http.StatusForbidden,
		Message: message,
		Err:     err}
}

func NotFound(message string, err error) *Error {
	return &Error{
		Code:    "not_found",
		Status:  404,
		Message: message,
		Err:     err,
	}
}

func Conflict(code string, message string, err error) *Error {
	return &Error{Code: code,
		Status:  http.StatusConflict,
		Message: message,
		Err:     err}
}

func Internal(err error) *Error {
	return &Error{Code: "internal_error",
		Status:  500,
		Message: "Internal Server Error",
		Err:     err}
}
