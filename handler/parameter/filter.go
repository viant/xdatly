package parameter

// Filter represents predicate filter
type Filter struct {
	Name    string
	Tag     string
	Include interface{}
	Exclude interface{}
}

// Filters represents a filter collection
type Filters []*Filter
