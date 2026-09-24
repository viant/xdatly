// Package client contains configuration shared by outbound client capabilities.
// HTTP and MCP contracts live in their respective subpackages.
package client

// Limits bounds one invocation. Timeout uses Go duration text. Implementations
// validate supported values and defaults; this package contains no networking.
type Limits struct {
	Timeout          string `json:"timeout,omitempty"`
	MaxResponseBytes int64  `json:"maxResponseBytes,omitempty"`
}
