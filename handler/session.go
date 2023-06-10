package handler

import (
	"github.com/viant/xdatly/handler/differ"
	"github.com/viant/xdatly/handler/mbus"
	"github.com/viant/xdatly/handler/sqlx"
	"github.com/viant/xdatly/handler/state"
	"github.com/viant/xdatly/handler/validator"
)

type Session interface {
	Validator() validator.Service
	Differ() differ.Service
	MessageBus() mbus.Service
	Db(opts ...sqlx.Option) sqlx.Service
	Stater() state.Stater
}
