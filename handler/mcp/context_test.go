package mcp

import (
	"context"
	"testing"
)

type testContext struct{}

func (testContext) Client() Client { return nil }

func TestContextRoundTrip(t *testing.T) {
	if _, ok := LookupContext(context.Background()); ok {
		t.Fatal("plain context must not expose MCP capabilities")
	}
	want := testContext{}
	actual, ok := LookupContext(WithContext(context.Background(), want))
	if !ok || actual != want {
		t.Fatalf("unexpected MCP context: actual=%#v ok=%v", actual, ok)
	}
}

func TestLookupContextRejectsNil(t *testing.T) {
	if _, ok := LookupContext(nil); ok {
		t.Fatal("nil context must not expose MCP capabilities")
	}
}
