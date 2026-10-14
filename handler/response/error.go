package response

import (
	"io"
	"strings"
)

// Error preserves the original public legacy error value. Err retains the
// private cause; presentation owns Message and Object.
type Error struct {
	View      string      `json:"view,omitempty" `
	Parameter string      `json:"parameter,omitempty" `
	Code      int         `json:"statusCode,omitempty" `
	Err       error       `json:"-"`
	Message   string      `json:"message,omitempty" `
	Object    interface{} `json:"object,omitempty" `
}

func (e *Error) StatusCode() int        { return e.Code }
func (e *Error) SetStatusCode(code int) { e.Code = code }
func (e *Error) Reader() io.Reader      { return strings.NewReader(e.Message) }
func (e *Error) Error() string          { return e.Message }
