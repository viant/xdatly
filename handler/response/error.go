package response

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Error represents an error
type Error struct {
	View      string      `json:"view,omitempty" `
	Parameter string      `json:"parameter,omitempty" `
	Code      int         `json:"statusCode,omitempty" `
	Err       error       `json:"-"`
	Message   string      `json:"message,omitempty" `
	Object    interface{} `json:"object,omitempty" `
}

func (e *Error) StatusCode() int {
	return e.Code
}

func (e *Error) SetStatusCode(code int) {
	e.Code = code
}

func (e *Error) Reader() io.Reader {
	return strings.NewReader(e.Message)
}

func (e *Error) Error() string {
	return e.Message
}

// NewParameterError creates a new parameter error
func NewParameterError(view, parameter string, err error, opts ...ErrorOption) *Error {
	ret := &Error{
		View:      view,
		Parameter: parameter,
		Err:       err,
	}
	for _, opt := range opts {
		opt(ret)
	}
	if coder, ok := err.(StatusCoder); ok && ret.Code == 0 {
		ret.Code = coder.StatusCode()
	}
	return ret
}

// NewError creates a new error
func NewError(code int, message string, opts ...ErrorOption) *Error {
	ret := &Error{
		Code:    code,
		Message: message,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			opt(ret)
		}
	}
	if ret.Err == nil {
		ret.Err = fmt.Errorf(ret.Message)
	}
	return ret
}

// BuildErrorResponse builds an error response
func BuildErrorResponse(err error) Response {
	statusCode := 400
	aResponse, ok := err.(Response)
	if !ok {
		aResponse = NewBuffered(WithStatusCode(statusCode))
	}
	return aResponse
}

// ErrorOption represents an error option
type ErrorOption func(e *Error)

// WithObject sets the object
func WithObject(object interface{}) ErrorOption {
	return func(e *Error) {
		if object == nil {
			return
		}
		objectType := reflect.TypeOf(object)
		if objectType.Kind() == reflect.Ptr {
			objectType = objectType.Elem()
		}
		switch objectType.Kind() {
		case reflect.Struct, reflect.Slice:
			e.Object = object
		}
	}
}

// WithError sets the error
func WithError(err error) ErrorOption {
	return func(e *Error) {
		e.Err = err
	}
}

// WithErrorStatusCode sets the error status code
func WithErrorStatusCode(code int) ErrorOption {
	return func(e *Error) {
		e.Code = code
	}
}

// WithErrorMessage sets the error message
func WithErrorMessage(message string) ErrorOption {
	return func(e *Error) {
		e.Message = message
	}
}

// WithView sets the error view
func WithView(view string) ErrorOption {
	return func(e *Error) {
		e.View = view
	}
}

// WithParameter sets the error parameter
func WithParameter(parameter string) ErrorOption {
	return func(e *Error) {
		e.Parameter = parameter
	}
}
