package exec

import (
	"context"
	"testing"
)

func TestCacheWarmupInvocationContext(t *testing.T) {
	ordinary := InvocationFromContext(context.Background())
	if ordinary.Purpose() != InvocationPurposeRequest || ordinary.IsCacheWarmup() || ordinary.MayBypassRowAuthorization() || ordinary.WarmupPhase() != WarmupPhaseNone {
		t.Fatalf("ordinary invocation = %+v", ordinary)
	}

	for _, phase := range []WarmupPhase{WarmupPhasePrepare, WarmupPhaseFill} {
		ctx := WithCacheWarmup(context.Background(), phase)
		actual := InvocationFromContext(ctx)
		if actual.Purpose() != InvocationPurposeCacheWarmup || !IsCacheWarmup(ctx) || !actual.IsCacheWarmup() || !actual.MayBypassRowAuthorization() || actual.WarmupPhase() != phase {
			t.Fatalf("warmup invocation for phase %v = %+v", phase, actual)
		}
	}
}

func TestInvalidWarmupPhaseDoesNotGrantAuthorizationBypass(t *testing.T) {
	actual := InvocationFromContext(WithCacheWarmup(context.Background(), WarmupPhaseNone))
	if !actual.IsCacheWarmup() || actual.MayBypassRowAuthorization() {
		t.Fatalf("invalid warmup phase = %+v", actual)
	}
}
