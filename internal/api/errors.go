package api

import (
	"fmt"
)

// StatusError can be used by internal code to suggest to the API controllers what status and message to return to the client instead of the
// default 'Internal error - something went wrong' error response.
type StatusError interface {
	error
	Status() Status
	Message() string
}

type Status string

const (
	StatusBadRequest    Status = "bad request"
	StatusNotFound      Status = "not found"
	StatusUnauthorized  Status = "unauthorized"
	StatusInternalError Status = "internal error"
)

var _ StatusError = (*statusError)(nil) // Interface guard

func (status Status) Error(messageFormat string, args ...any) StatusError {
	return Error(status, messageFormat, args...)
}

func (status Status) ErrorWithCause(cause error, messageFormat string, args ...any) StatusError {
	return ErrorWithCause(status, cause, messageFormat, args...)
}

func Error(status Status, messageFormat string, args ...any) StatusError {
	return statusError{
		status:  status,
		message: fmt.Sprintf(messageFormat, args...),
	}
}

func ErrorWithCause(status Status, cause error, messageFormat string, args ...any) StatusError {
	return statusError{
		status:  status,
		message: fmt.Sprintf(messageFormat, args...),
		cause:   cause,
	}
}

type statusError struct {
	status  Status
	message string
	cause   error
}

func (statusErr statusError) Error() string {
	if statusErr.cause == nil {
		return fmt.Sprintf("%s: %s", statusErr.status, statusErr.message)
	}
	return fmt.Sprintf("%s: %s: %v", statusErr.status, statusErr.message, statusErr.cause)
}

func (statusErr statusError) Status() Status {
	return statusErr.status
}

func (statusErr statusError) Message() string {
	return statusErr.message
}

func (statusErr statusError) Unwrap() error {
	return statusErr.cause
}
