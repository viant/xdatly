package http

import "github.com/viant/xdatly/handler/response"

// NewError creates a new error
func NewError(code int, message string, object interface{}) *response.Error {
	return response.NewError(code, message, response.WithObject(object))
}
