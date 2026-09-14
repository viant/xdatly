package mcp

import "context"

// Initializer is an MCP-aware input hook. The unified handler engine invokes
// it after the ordinary input initializer when MCP context is present.
type Initializer interface {
	InitMCP(ctx context.Context, mcp Context) error
}

// Finalizer is an MCP-aware output hook. The unified handler engine invokes it
// after ordinary successful finalization when MCP context is present.
type Finalizer interface {
	FinalizeMCP(ctx context.Context, mcp Context) error
}
