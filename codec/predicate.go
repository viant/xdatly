package codec

import "context"

type (
	Predicate struct {
		Parameter string //name of parameter
		Template  string //predicate template
		Criteria
	}

	PredicateHandler interface {
		Compute(ctx context.Context, value interface{}) (*Criteria, error)
	}
)
