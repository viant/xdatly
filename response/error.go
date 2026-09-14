package response

import (
	"errors"
	"net/http"
)

// BodyError explicitly supplies a client-visible error body. The body is not
// merged with internal diagnostics; nil deliberately means JSON null.
type BodyError interface {
	error
	StatusCoder
	ResponseBody() any
}

// Error separates the public status/payload from a private, wrappable cause.
// Payload may be a typed struct or map; its own JSON tags control empty/null
// fields. Never place confidential diagnostics in Payload.
type Error struct {
	Code    int   `json:"-"`
	Payload any   `json:"-"`
	Cause   error `json:"-"`
}

func (e *Error) Error() string {
	if e != nil && e.Cause != nil {
		return e.Cause.Error()
	}
	return http.StatusText(e.StatusCode())
}
func (e *Error) StatusCode() int {
	if e == nil {
		return http.StatusInternalServerError
	}
	if e.Code == 0 {
		return ErrorStatusCode(e.Cause, http.StatusInternalServerError)
	}
	return e.Code
}
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
func (e *Error) ResponseBody() any {
	if e == nil {
		return nil
	}
	return e.Payload
}

// ErrorBody follows ordinary wrapping (including errors.Join). A true result
// is explicit even when the returned body is nil.
func ErrorBody(err error) (any, bool) {
	var public BodyError
	if err == nil || !errors.As(err, &public) {
		return nil, false
	}
	return public.ResponseBody(), true
}
