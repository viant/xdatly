package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/viant/mcp-protocol/schema"
	xhttp "github.com/viant/xdatly/client/http"
	xmcp "github.com/viant/xdatly/client/mcp"
)

// Standard HTTP clients satisfy the contract without a runtime-specific shim.
var _ xhttp.Client = (*http.Client)(nil)

type httpProvider struct{}

func (httpProvider) Client(context.Context, xhttp.Options) (xhttp.Client, error) {
	return &http.Client{}, nil
}

type mcpClient struct{}

func (mcpClient) CallTool(context.Context, *schema.CallToolRequest, xmcp.CallOptions) (*schema.CallToolResult, error) {
	return &schema.CallToolResult{}, nil
}

type mcpProvider struct{}

func (mcpProvider) Client(context.Context, xmcp.Options) (xmcp.Client, error) {
	return mcpClient{}, nil
}

var _ xhttp.Provider = httpProvider{}
var _ xmcp.Provider = mcpProvider{}

func TestProtocolContractsRemainDistinct(t *testing.T) {
	if _, ok := any(mcpClient{}).(xhttp.Client); ok {
		t.Fatal("MCP client unexpectedly satisfies HTTP contract")
	}
	if _, ok := any(&http.Client{}).(xmcp.Client); ok {
		t.Fatal("HTTP client unexpectedly satisfies MCP contract")
	}
	client, err := (mcpProvider{}).Client(context.Background(), xmcp.Options{Tool: "context"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.CallTool(context.Background(), &schema.CallToolRequest{}, xmcp.CallOptions{}); err != nil {
		t.Fatal(err)
	}
}
