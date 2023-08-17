package http

import (
	"context"
	"net/http"
)

type (
	Http interface {
		RawRequest() *http.Request
		RequestOf(ctx context.Context, state interface{}) (*http.Request, error)
		RouteRequest(ctx context.Context) (*http.Request, error)
		FailWithCode(statusCode int, err error) error
	}
)
