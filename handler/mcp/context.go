// Package mcp defines the MCP capabilities and lifecycle contracts exposed to
// xdatly handlers without coupling them to a Datly runtime implementation.
package mcp

import (
	"context"

	"github.com/viant/mcp-protocol/schema"
)

// Client exposes the MCP client capabilities available to a handler hook.
type Client interface {
	CanElicit() bool
	CanGenerateContent() bool
	Elicit(ctx context.Context, params *schema.ElicitRequestParams) (*schema.ElicitResult, error)
	GenerateContent(ctx context.Context, params *schema.CreateMessageRequestParams) (*schema.CreateMessageResult, error)
}

// Context carries request-scoped MCP capabilities.
type Context interface {
	Client() Client
}

type contextKey struct{}

// WithContext attaches MCP capabilities to an invocation context.
func WithContext(ctx context.Context, value Context) context.Context {
	return context.WithValue(ctx, contextKey{}, value)
}

// LookupContext returns the MCP capabilities attached to ctx.
func LookupContext(ctx context.Context) (Context, bool) {
	if ctx == nil {
		return nil, false
	}
	value, ok := ctx.Value(contextKey{}).(Context)
	return value, ok && value != nil
}
