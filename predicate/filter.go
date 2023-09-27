package predicate

type NamedFilters []*NamedFilter

type NamedFilter struct {
	Name    string   `xls:"name=Filter"`
	Include []string `json:",omitempty"`
	Exclude []string `json:",omitempty"`
}

type StringsFilter struct {
	Include []string `json:",omitempty"`
	Exclude []string `json:",omitempty"`
}

type IntFilter struct {
	Include []int `json:",omitempty"`
	Exclude []int `json:",omitempty"`
}

type BoolFilter struct {
	Include []bool `json:",omitempty"`
	Exclude []bool `json:",omitempty"`
}
