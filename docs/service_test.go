package docs

import (
	"context"
	"database/sql"
	"testing"
)

type testConnector struct{}

func (testConnector) DB() (*sql.DB, error) { return nil, nil }

type testService struct{}

func (testService) Lookup(_ context.Context, key string) (string, bool, error) {
	return key, true, nil
}

type testProvider struct{}

func (testProvider) Service(_ context.Context, _ ...Option) (Service, error) {
	return testService{}, nil
}

func TestDocsContractsCompile(t *testing.T) {
	var _ Connector = testConnector{}
	var _ Service = testService{}
	var _ Provider = testProvider{}
}

func TestDocsOptionsHelpers(t *testing.T) {
	var opts Options
	WithURL("docs/orders.md")(&opts)
	WithConnector(testConnector{})(&opts)

	if opts.URL != "docs/orders.md" {
		t.Fatalf("expected URL option to be preserved")
	}
	if opts.Connector == nil {
		t.Fatalf("expected connector option to be preserved")
	}
}
