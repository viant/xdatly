package handler

import "context"

type sessionContextKey struct{}

// WithSession attaches one public handler Session to ctx.
func WithSession(ctx context.Context, sess Session) context.Context {
	if sess == nil {
		return ctx
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, sessionContextKey{}, sess)
}

// SessionFromContext resolves the public handler Session from ctx when one was
// previously attached with WithSession.
func SessionFromContext(ctx context.Context) (Session, bool) {
	if ctx == nil {
		return nil, false
	}
	value := ctx.Value(sessionContextKey{})
	if value == nil {
		return nil, false
	}
	sess, ok := value.(Session)
	return sess, ok
}
