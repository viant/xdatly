package handler

import (
	"context"
	"github.com/viant/xdatly/handler/differ"
	"github.com/viant/xdatly/handler/http"
	"github.com/viant/xdatly/handler/mbus"
	"github.com/viant/xdatly/handler/sqlx"
	"github.com/viant/xdatly/handler/state"
	"github.com/viant/xdatly/handler/validator"
)

type Session interface {
	Validator() *validator.Service
	Differ() *differ.Service
	MessageBus() *mbus.Service
	Db(opts ...sqlx.Option) (*sqlx.Service, error)
	Stater() *state.Service
	FlushTemplate(ctx context.Context) error
	Route(ctx context.Context, route *http.Route) (Session, error)
	Http() http.Http
}
