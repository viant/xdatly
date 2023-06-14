package handler

import (
	"github.com/viant/xdatly/handler/differ"
	"github.com/viant/xdatly/handler/mbus"
	"github.com/viant/xdatly/handler/sqlx"
	"github.com/viant/xdatly/handler/state"
	"github.com/viant/xdatly/handler/validator"
)

type ISession interface {
	Validator() *validator.Service
	Differ() *differ.Service
	MessageBus() *mbus.Service
	Db(opts ...sqlx.Option) *sqlx.Service
	Stater() *state.Service
}

type Session struct {
	session ISession
}

func (s *Session) Validator() *validator.Service {
	return s.session.Validator()
}
func (s *Session) Differ() *differ.Service {
	return s.session.Differ()
}

func (s *Session) MessageBus() *mbus.Service {
	return s.session.MessageBus()
}

func (s *Session) Db(opts ...sqlx.Option) *sqlx.Service {
	return s.session.Db(opts...)
}
func (s *Session) Stater() *state.Service {
	return s.session.Stater()
}

func NewSession(session ISession) *Session {
	return &Session{session: session}
}
