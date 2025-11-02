package apperror

import "net/http"

type HTTPError struct {
	Err        error
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

func NewHTTPError(err error, statusCode int) *HTTPError {
	return &HTTPError{
		Err:        err,
		StatusCode: statusCode,
	}
}

func NewHTTPErrorWithMessage(err error, statusCode int, message string) *HTTPError {
	return &HTTPError{
		Err:        err,
		StatusCode: statusCode,
		Message:    message,
	}
}

func BadRequest(err error) *HTTPError {
	return NewHTTPError(err, http.StatusBadRequest)
}

func Unauthorized(err error) *HTTPError {
	return NewHTTPError(err, http.StatusUnauthorized)
}

func Forbidden(err error) *HTTPError {
	return NewHTTPError(err, http.StatusForbidden)
}

func NotFound(err error) *HTTPError {
	return NewHTTPError(err, http.StatusNotFound)
}

func InternalServerError(err error) *HTTPError {
	return NewHTTPError(err, http.StatusInternalServerError)
}
