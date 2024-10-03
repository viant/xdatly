package http

import (
	"context"
	"github.com/viant/xdatly/handler/state"
	"net/http"
)

type (
	Http interface {
		RequestOf(ctx context.Context, any interface{}) (*http.Request, error)
		NewRequest(ctx context.Context, opts ...state.Option) (*http.Request, error)
		Redirect(ctx context.Context, route *Route, request *http.Request) error
		FailWithCode(statusCode int, err error) error
	}
)
