package state

// Selector is the public invocation-scoped query-selector contract.
type Selector struct {
	Columns      []string      `json:",omitempty"`
	Fields       []string      `json:",omitempty"`
	OrderBy      string        `json:",omitempty"`
	Offset       int           `json:",omitempty"`
	Limit        int           `json:",omitempty"`
	Page         int           `json:",omitempty"`
	Criteria     string        `json:",omitempty"`
	Placeholders []interface{} `json:",omitempty"`
}

func (s *Selector) Clone() *Selector {
	if s == nil {
		return nil
	}
	result := *s
	result.Columns = append([]string(nil), s.Columns...)
	result.Fields = append([]string(nil), s.Fields...)
	result.Placeholders = append([]interface{}(nil), s.Placeholders...)
	return &result
}

func (s *Selector) CurrentLimit() int {
	return s.Limit
}

func (s *Selector) CurrentOffset() int {
	return s.Offset
}

func (s *Selector) CurrentPage() int {
	return s.Page
}

func (s *Selector) SetCriteria(expanded string, placeholders []interface{}) {
	s.Criteria = expanded
	s.Placeholders = placeholders
}

// NamedSelector binds a selector to one named view/query scope.
type NamedSelector struct {
	Name string `json:",omitempty"`
	Selector
}

func (s *NamedSelector) Clone() *NamedSelector {
	if s == nil {
		return nil
	}
	result := &NamedSelector{Name: s.Name}
	if selector := s.Selector.Clone(); selector != nil {
		result.Selector = *selector
	}
	return result
}

// Selectors is the reduced named-selector collection contract.
type Selectors []*NamedSelector

func (s Selectors) Find(name string) *NamedSelector {
	for _, n := range s {
		if n != nil && n.Name == name {
			return n
		}
	}
	return nil
}

func (s Selectors) Clone() Selectors {
	if len(s) == 0 {
		return nil
	}
	result := make(Selectors, len(s))
	for i, selector := range s {
		result[i] = selector.Clone()
	}
	return result
}
