package response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type privateLegacyDependencyFailure struct{ cause error }

func (e *privateLegacyDependencyFailure) Error() string { return "private dependency diagnostic" }
func (e *privateLegacyDependencyFailure) Unwrap() error { return e.cause }

func TestLegacyErrorPreservesPrivateCauseIdentity(t *testing.T) {
	private := &privateLegacyDependencyFailure{cause: context.Canceled}
	legacy := &Error{Code: 500, Message: "internal server error: InternalError", Err: fmt.Errorf("dependency: %w", private)}
	wrapped := fmt.Errorf("campaign: %w", legacy)
	if !errors.Is(wrapped, context.Canceled) {
		t.Fatal("legacy error lost cancellation cause")
	}
	var actual *privateLegacyDependencyFailure
	if !errors.As(wrapped, &actual) || actual != private {
		t.Fatal("legacy error lost typed dependency cause")
	}
	if legacy.Error() != "internal server error: InternalError" || legacy.StatusCode() != 500 {
		t.Fatal("private cause changed public error")
	}
	raw, err := json.Marshal(legacy)
	if err != nil || string(raw) != `{"statusCode":500,"message":"internal server error: InternalError"}` {
		t.Fatalf("private cause changed legacy JSON: %s %v", raw, err)
	}
}

func TestLegacyErrorUnwrapNil(t *testing.T) {
	var legacy *Error
	if legacy.Unwrap() != nil || (&Error{}).Unwrap() != nil {
		t.Fatal("empty legacy cause must unwrap to nil")
	}
}
