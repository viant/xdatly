package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestExplicitErrorPreservesBodyAndCause(t *testing.T) {
	secret := errors.New("private database detail")
	for _, payload := range []any{nil, map[string]any{"message": "", "error": nil, "violations": []any{}}, map[string]any{"message": "denied", "error": map[string]any{"reason": "policy"}}} {
		public := &Error{Code: 401, Payload: payload, Cause: secret}
		wrapped := fmt.Errorf("handler failed: %w", public)
		body, ok := ErrorBody(wrapped)
		if !ok || ErrorStatusCode(wrapped, 500) != 401 || !errors.Is(wrapped, secret) {
			t.Fatal("wrapped error lost control or cause")
		}
		got, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		want, _ := json.Marshal(payload)
		if string(got) != string(want) {
			t.Fatalf("body=%s", got)
		}
		encoded, err := json.Marshal(public)
		if err != nil || string(encoded) != "{}" {
			t.Fatalf("internal fields serialized: %s %v", encoded, err)
		}
	}
	if _, ok := ErrorBody(secret); ok {
		t.Fatal("ordinary internal error became public payload")
	}
}
