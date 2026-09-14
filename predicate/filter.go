package predicate

// NamedFilters is the reduced public filter payload contract.
type NamedFilters []*NamedFilter

// NamedFilter groups include/exclude values under one named filter.
type NamedFilter struct {
	Name    string   `xls:"name=Filter"`
	Include []string `json:",omitempty"`
	Exclude []string `json:",omitempty"`
}

// StringsFilter is the reduced string filter payload shape.
type StringsFilter struct {
	Include []string `json:",omitempty"`
	Exclude []string `json:",omitempty"`
}

// IntFilter is the reduced int filter payload shape.
type IntFilter struct {
	Include []int `json:",omitempty"`
	Exclude []int `json:",omitempty"`
}

// BoolFilter is the reduced bool filter payload shape.
type BoolFilter struct {
	Include []bool `json:",omitempty"`
	Exclude []bool `json:",omitempty"`
}
