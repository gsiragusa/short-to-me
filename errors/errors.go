package errors

import "net/http"

type Error struct {
	Message    string `json:"message"`
	HttpStatus int    `json:"http_status"`
}

func NewErrorNotFound() Error {
	return Error{
		Message:    "The resource was not found",
		HttpStatus: http.StatusNotFound,
	}
}

func NewErrorBadRequest() Error {
	return Error{
		Message:    "There was something wrong with your request",
		HttpStatus: http.StatusBadRequest,
	}
}

func NewInternalServerError() Error {
	return Error{
		Message:    "An unexpected error occurred. Please try again later",
		HttpStatus: http.StatusInternalServerError,
	}
}

func ErrResponse(err Error) *Error {
	return &Error{
		Message:    err.Message,
		HttpStatus: err.HttpStatus,
	}
}
