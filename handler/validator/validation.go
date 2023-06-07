package validator

type (
	Violation struct {
		Location string
		Field    string
		Value    interface{}
		Message  string
		Check    string
	}

	Validation struct {
		Violations []*Violation
		Failed     bool
	}
)
