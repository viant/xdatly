package predicate

type (
	Template struct {
		Name   string
		Source string
		Args   []*NamedArgument
	}

	NamedArgument struct {
		Name     string
		Position int
	}
)
