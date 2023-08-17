package codec

type (
	Predicate struct {
		Parameter string //name of parameter
		Template  string //predicate template
		Criteria
	}
)
