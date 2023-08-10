package http

import (
	"net/http"
)

type (
	Http interface {
		RawRequest() *http.Request
		RequestOf(state interface{}) *http.Request
		RouteRequest() *http.Request
		FailWithCode(statusCode int, err error) error
	}
)
