package handler

import "context"

// Func is the typed pure-function handler contract. Use Contract when a custom
// handler needs scoped DI or response control.
type Func[I any, O any] func(ctx context.Context, input *I) (*O, error)
