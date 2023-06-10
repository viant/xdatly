package state

import "context"

type Stater interface {
	Into(ctx context.Context, state interface{}) error
}
