package mcp

import "github.com/viant/xdatly/client"

// Options describes an MCP server, protocol, tool and session bounds. Per-call
// credentials are supplied explicitly through CallOptions rather than inferred.
type Options struct {
	URL             string            `json:"url"`
	Transport       string            `json:"transport,omitempty"`
	ProtocolVersion string            `json:"protocolVersion,omitempty"`
	Tool            string            `json:"tool"`
	Header          map[string]string `json:"header,omitempty"`
	MaxSessions     int               `json:"maxSessions,omitempty"`
	client.Limits
}
