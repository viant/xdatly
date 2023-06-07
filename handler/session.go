package handler

import (
	"github.com/viant/xdatly/handler/differ"
	"github.com/viant/xdatly/handler/mbus"
	"github.com/viant/xdatly/handler/response"
	"github.com/viant/xdatly/handler/sqlx"
	"github.com/viant/xdatly/handler/validator"
)

type Session interface {
	Validator() validator.Service
	Differ() differ.Service
	MessageBus() mbus.Service
	Db(opts ...sqlx.Option) sqlx.Service
	Response() response.Response
	StateInto(dest interface{}) error
}

/*
	{"ViewHandler":"Foo"}


*/
