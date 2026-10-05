package handler

import (
	"context"
	"errors"
	"testing"
)

type eligibilitySDKEntity struct{ ID int }
type eligibilitySDKOutput struct{}
type eligibilitySDKHook[T, P, O any] struct {
	eligible bool
	err      error
}

func (h *eligibilitySDKHook[T, P, O]) WriteEligible(context.Context, *T, LifecycleContext[T, P, O], WriteAction) (bool, error) {
	return h.eligible, h.err
}

type eligibilitySDKPromotedHook struct {
	eligibilitySDKHook[eligibilitySDKEntity, NoParent, eligibilitySDKOutput]
}

var _ WriteEligibilityHook[eligibilitySDKEntity, NoParent, eligibilitySDKOutput] = (*eligibilitySDKHook[eligibilitySDKEntity, NoParent, eligibilitySDKOutput])(nil)
var _ WriteEligibilityHook[eligibilitySDKEntity, NoParent, eligibilitySDKOutput] = (*eligibilitySDKPromotedHook)(nil)

func TestWriteEligibilityHookTypedOptionalContract(t *testing.T) {
	if _, ok := any(&entityHook{}).(WriteEligibilityHook[entityRecord, entityParent, entityOutput]); ok {
		t.Fatal("base lifecycle unexpectedly requires eligibility")
	}
	cause := errors.New("eligibility rejected")
	for _, tc := range []struct {
		eligible bool
		err      error
	}{{true, nil}, {false, nil}, {false, cause}} {
		var hook WriteEligibilityHook[eligibilitySDKEntity, NoParent, eligibilitySDKOutput] = &eligibilitySDKPromotedHook{eligibilitySDKHook: eligibilitySDKHook[eligibilitySDKEntity, NoParent, eligibilitySDKOutput]{eligible: tc.eligible, err: tc.err}}
		got, err := hook.WriteEligible(context.Background(), &eligibilitySDKEntity{}, LifecycleContext[eligibilitySDKEntity, NoParent, eligibilitySDKOutput]{Output: &eligibilitySDKOutput{}}, WriteInsert)
		if got != tc.eligible || !errors.Is(err, tc.err) {
			t.Fatalf("WriteEligible = %v, %v", got, err)
		}
	}
}
