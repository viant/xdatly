package parameter

type (
	Predicate interface {
		Expand(value interface{}) (*Criteria, error)
	}

	PredicateFactory interface {
		NewPredicate(args []interface{}, options ...interface{}) (Predicate, error)
	}
)
