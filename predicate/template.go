package predicate

// Template is the public predicate-template contract.
type Template struct {
	Name   string
	Source string
	Args   []*NamedArgument
}

// NamedArgument describes one named template argument.
type NamedArgument struct {
	Name     string
	Position int
}

// Lookup defines the minimal read-only predicate-template registry contract.
type Lookup interface {
	Lookup(name string) (*Template, error)
}
