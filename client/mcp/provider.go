package mcp

import "context"

// Provider resolves or reuses an MCP client from explicit deployment options.
// Runtime composition supplies its implementation through native DI.
type Provider interface {
	Client(ctx context.Context, options Options) (Client, error)
}
