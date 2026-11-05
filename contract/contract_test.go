package contract_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/viant/bindly"
	"github.com/viant/bindly/provider/request"
	"github.com/viant/xdatly/contract"
)

type exampleContractProvider struct{}

func (*exampleContractProvider) Contract(ctx context.Context, binder bindly.Binder) (*contract.Policy, error) {
	// The arbitrary state belongs to the consumer, not the contract interface.
	var state struct {
		Debug bool `bind:"kind=header,in=X-Debug"`
	}
	if err := binder.Bind(ctx, &state); err != nil {
		return nil, err
	}
	if !state.Debug {
		return nil, nil
	}
	return &contract.Policy{Name: "debug"}, nil
}

var _ contract.Provider = (*exampleContractProvider)(nil)

func TestProviderAcceptsNativeBindlyWithoutSDKLookup(t *testing.T) {
	root, err := bindly.NewInjector()
	if err != nil {
		t.Fatal(err)
	}
	for _, debug := range []bool{false, true} {
		req := httptest.NewRequest("GET", "/", nil)
		if debug {
			req.Header.Set("X-Debug", "true")
		}
		transport, err := request.New(req)
		if err != nil {
			t.Fatal(err)
		}
		scope, err := root.ForScope(transport.Providers()...)
		if err != nil {
			t.Fatal(err)
		}
		policy, err := (&exampleContractProvider{}).Contract(context.Background(), scope)
		if err != nil {
			t.Fatal(err)
		}
		if debug && (policy == nil || policy.Name != "debug") || !debug && policy != nil {
			t.Fatalf("debug=%v policy=%+v", debug, policy)
		}
		transport.Close()
	}
}
