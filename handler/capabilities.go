package handler

import (
	"context"

	"github.com/viant/xdatly/connector"
	"github.com/viant/xdatly/differ"
	xlogger "github.com/viant/xdatly/logger"
	xmbus "github.com/viant/xdatly/mbus"
)

const (
	// DifferKey identifies the typed comparison capability.
	DifferKey ValueKey = "differ"
	// LoggerKey identifies the invocation-scoped logger capability.
	LoggerKey ValueKey = "logger"
	// ValidatorKey identifies the invocation-scoped validation capability.
	ValidatorKey ValueKey = "validator"
	// MessageBusKey identifies the invocation-scoped message-bus capability.
	MessageBusKey ValueKey = "mbus"
	// ConnectorKey is an opt-in exact-name DB provider for specialized handlers.
	ConnectorKey ValueKey = "connector"
)

type Logger = xlogger.Logger

// Validator is the injectable validation service available to custom handlers.
type Validator interface {
	Validate(ctx context.Context, value any, opts ...any) (*Validation, error)
}

type MessageBus = xmbus.Service

// Capabilities groups fixed handler primitives supplied by an invocation.
// Database Data remains invocation-owned and is exposed independently by
// DataKey.
type Capabilities struct {
	Differ     differ.Differ
	Logger     Logger
	Validator  Validator
	MessageBus MessageBus
	// Connector is absent unless the application explicitly grants DB access.
	// It is not part of managed Data and never provides a transaction implicitly.
	Connector connector.Provider
}
