package logger

import "context"

type contextKey struct{}
type contextValue struct{ logger Logger }

// WithContext attaches the invocation's configured logger. A nil logger shadows
// an inherited logger; it never substitutes a global or default logger.
func WithContext(ctx context.Context, value Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, contextValue{logger: value})
}

// FromContext returns the same logger attached by WithContext, or nil if absent.
func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return nil
	}
	value, _ := ctx.Value(contextKey{}).(contextValue)
	return value.logger
}
