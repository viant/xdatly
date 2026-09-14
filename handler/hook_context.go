package handler

import "context"

type hookAppliedContextKey struct{}

// WithHookApplied marks ctx so a pre-invoke hook can recognize that it has
// already run for the current invocation.
func WithHookApplied(ctx context.Context) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, hookAppliedContextKey{}, true)
}

// HookApplied reports whether ctx already carries the pre-invoke hook-applied
// marker.
func HookApplied(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	applied, _ := ctx.Value(hookAppliedContextKey{}).(bool)
	return applied
}

// RunPreInvokeHook runs hook unless the context already carries the
// hook-applied marker.
func RunPreInvokeHook(ctx context.Context, input any, hook func(context.Context, any) error) error {
	if hook == nil {
		return nil
	}
	if HookApplied(ctx) {
		return nil
	}
	return hook(ctx, input)
}

// InvokeWithSession marks the hook as already applied, attaches sess to the
// invocation context, and then delegates to invoke.
func InvokeWithSession(
	ctx context.Context,
	sess Session,
	input any,
	invoke func(context.Context, any) (any, error),
) (any, error) {
	hookApplied := WithHookApplied(ctx)
	return invoke(WithSession(hookApplied, sess), input)
}
