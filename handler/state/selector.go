package state

type QuerySelector struct {
	Columns      []string      `json:",omitempty"`
	Fields       []string      `json:",omitempty"`
	OrderBy      string        `json:",omitempty"`
	Offset       int           `json:",omitempty"`
	Limit        int           `json:",omitempty"`
	Page         int           `json:",omitempty"`
	Criteria     string        `json:",omitempty"`
	Placeholders []interface{} `json:",omitempty"`
}

func (s *QuerySelector) CurrentLimit() int {
	return s.Limit
}

func (s *QuerySelector) CurrentOffset() int {
	return s.Offset
}

func (s *QuerySelector) CurrentPage() int {
	return s.Page
}

func (s *QuerySelector) SetCriteria(expanded string, placeholders []interface{}) {
	s.Criteria = expanded
	s.Placeholders = placeholders
}

type NamedQuerySelector struct {
	Name string `json:",omitempty"`
	QuerySelector
}

type QuerySelectors []*NamedQuerySelector

func (s QuerySelectors) Find(name string) *NamedQuerySelector {
	for _, n := range s {
		if n.Name == name {
			return n
		}
	}
	return nil
}
