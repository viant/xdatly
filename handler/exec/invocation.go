// Package exec provides dependency-light execution context contracts shared by
// xdatly handlers and runtimes.
package exec

import "context"

// InvocationPurpose identifies why a handler or predicate is executing.
type InvocationPurpose uint8

const (
	InvocationPurposeRequest InvocationPurpose = iota
	InvocationPurposeCacheWarmup
)

// WarmupPhase distinguishes warmup validation from the cache-filling phase.
type WarmupPhase uint8

const (
	WarmupPhaseNone WarmupPhase = iota
	WarmupPhasePrepare
	WarmupPhaseFill
)

// InvocationInfo is immutable execution metadata propagated through context.
// Private fields prevent consumers from changing the capability selected by
// the runtime after retrieving it.
type InvocationInfo struct {
	purpose                   InvocationPurpose
	warmupPhase               WarmupPhase
	mayBypassRowAuthorization bool
}

func (i *InvocationInfo) Purpose() InvocationPurpose {
	if i == nil {
		return InvocationPurposeRequest
	}
	return i.purpose
}

func (i *InvocationInfo) IsCacheWarmup() bool {
	return i != nil && i.purpose == InvocationPurposeCacheWarmup
}

func (i *InvocationInfo) WarmupPhase() WarmupPhase {
	if i == nil {
		return WarmupPhaseNone
	}
	return i.warmupPhase
}

// MayBypassRowAuthorization reports an explicit runtime capability. Consumers
// should not infer authorization policy from IsCacheWarmup alone.
func (i *InvocationInfo) MayBypassRowAuthorization() bool {
	return i != nil && i.mayBypassRowAuthorization
}

type invocationInfoKey struct{}

var requestInvocation = &InvocationInfo{purpose: InvocationPurposeRequest}

// WithCacheWarmup marks a trusted runtime-owned warmup phase. Transport values
// must never be converted into this context capability.
func WithCacheWarmup(ctx context.Context, phase WarmupPhase) context.Context {
	mayBypass := phase == WarmupPhasePrepare || phase == WarmupPhaseFill
	return context.WithValue(ctx, invocationInfoKey{}, &InvocationInfo{
		purpose:                   InvocationPurposeCacheWarmup,
		warmupPhase:               phase,
		mayBypassRowAuthorization: mayBypass,
	})
}

// InvocationFromContext returns request metadata when no explicit invocation
// descriptor has been installed.
func InvocationFromContext(ctx context.Context) *InvocationInfo {
	if ctx != nil {
		if result, ok := ctx.Value(invocationInfoKey{}).(*InvocationInfo); ok && result != nil {
			return result
		}
	}
	return requestInvocation
}

func IsCacheWarmup(ctx context.Context) bool {
	return InvocationFromContext(ctx).IsCacheWarmup()
}
