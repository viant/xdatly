package response

import (
	"net/http"
	"sync"
)

// Errors represents a list of errors
type Errors struct {
	Message string `json:",omitempty" `
	Errors  []*Error
	mutex   sync.Mutex
	status  int
}

// Error represents an error
func (e *Errors) Error() string {
	if e.Message == "" {
		if len(e.Errors) > 0 {
			return e.Errors[0].Error()
		}
	}
	return e.Message
}

// Append appends an error
func (e *Errors) Append(err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	switch actual := err.(type) {
	case *Error:
		e.Errors = append(e.Errors, actual)
		e.updateStatusCode(actual.Code)
	case *Errors:
		e.Errors = append(e.Errors, actual.Errors...)
		code := actual.status
		e.updateStatusCode(code)
	default:
		e.Errors = append(e.Errors, &Error{Message: err.Error(), Err: err})
	}
}

func (e *Errors) StatusCode() int {
	return e.status
}

// AddError adds an error
func (e *Errors) AddError(view, param string, err error, opts ...ErrorOption) {
	if err == nil {
		return
	}
	e.mutex.Lock()
	e.Errors = append(e.Errors, NewParameterError(view, param, err, opts...))
	e.Message = err.Error()
	e.mutex.Unlock()
}

// SetStatusCode sets status code
func (e *Errors) SetStatusCode(status int) {
	e.updateStatusCode(status)
}

func (e *Errors) updateStatusCode(code int) {
	if statusCodePriority(code) > statusCodePriority(e.status) {
		e.status = code
	}
}

// HasError	returns true if has error
func (e *Errors) HasError() bool {
	return len(e.Errors) > 0
}

// NewErrors creates a new errors
func NewErrors() *Errors {
	return &Errors{Errors: make([]*Error, 0)}
}

const (
	priorityDefault = iota
	priority400
	priority404
	priority403
	priority401
)

func statusCodePriority(status int) int {
	switch status {
	case http.StatusUnauthorized:
		return priority401
	case http.StatusForbidden:
		return priority403
	case http.StatusNotFound:
		return priority404
	case http.StatusBadRequest:
		return priority400
	case 0:
		return -1
	default:
		return priorityDefault
	}
}
