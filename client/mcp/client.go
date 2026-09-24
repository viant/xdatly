// Package mcp defines outbound MCP capabilities using protocol contracts rather
// than any concrete MCP implementation. Clients are borrowed from providers.
package mcp

import (
	"context"
	"net/http"

	"github.com/viant/mcp-protocol/schema"
)

// CallOptions contains explicit per-call transport metadata. Header values are
// forwarded unchanged during initialization, discovery, and tool invocation;
// implementations must isolate sessions with different header sets.
type CallOptions struct {
	Header http.Header
}

// Client invokes an MCP tool and preserves its native result semantics.
// A tool's isError result is distinct from a transport or protocol error.
// Lifecycle and session cleanup belong to the provider's owner.
type Client interface {
	CallTool(ctx context.Context, request *schema.CallToolRequest, options CallOptions) (*schema.CallToolResult, error)
}
