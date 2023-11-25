package response

import (
	"strings"
	"time"
)

type (
	Metric struct {
		View      string     `json:",omitempty"`
		Elapsed   string     `json:",omitempty"`
		ElapsedMs int        `json:",omitempty"`
		Rows      int        `json:",omitempty"`
		Execution *Execution `json:",omitempty"`
	}

	Metrics    []*Metric
	CacheStats struct {
		Type           string
		RecordsCounter int
		Key            string
		Dataset        string
		Namespace      string
		FoundWarmup    bool   `json:",omitempty"`
		FoundLazy      bool   `json:",omitempty"`
		ErrorType      string `json:",omitempty"`
		ErrorCode      int    `json:",omitempty"`
		ExpiryTime     *time.Time
	}

	Execution struct {
		View    []*SQLExecution `json:",omitempty"`
		Summary []*SQLExecution `json:",omitempty"`
		Elapsed string          `json:",omitempty"`
	}

	ParametrizedSQL struct {
		Query string
		Args  []interface{}
	}

	SQLExecution struct {
		SQL        string        `json:",omitempty"`
		Args       []interface{} `json:",omitempty"`
		CacheStats *CacheStats   `json:",omitempty"`
		Error      string        `json:",omitempty"`
		CacheError string        `json:",omitempty"`
	}
)

func (m *Metric) ParametrizedSQL() []*ParametrizedSQL {
	var result = make([]*ParametrizedSQL, 0)
	if m.Execution == nil {
		return result
	}
	for _, tmpl := range m.Execution.View {
		result = append(result, &ParametrizedSQL{Query: tmpl.SQL, Args: tmpl.Args})
	}
	return result
}

func (m *Metric) Name() string {
	return strings.Title(m.View)
}

func (m *Metric) SQL() string {
	if m.Execution != nil && len(m.Execution.View) > 0 {
		tmpl := m.Execution.View[0]
		SQL := ExpandSQL(tmpl.SQL, tmpl.Args)
		return SQL
	}
	return ""
}

func (m *Metric) HideSQL() *Metric {
	ret := *m
	if m.Execution == nil {
		return &ret
	}
	ret.Execution = &Execution{
		View:    make([]*SQLExecution, len(m.Execution.View)),
		Summary: make([]*SQLExecution, len(m.Execution.Summary)),
	}
	copy(ret.Execution.View, m.Execution.View)
	copy(ret.Execution.Summary, m.Execution.Summary)
	for _, elem := range m.Execution.View {
		elem.SQL = ""
		elem.Args = nil
	}
	for _, elem := range m.Execution.Summary {
		elem.SQL = ""
		elem.Args = nil
	}
	return &ret
}

// Basic returns n
func (m Metrics) Basic() Metrics {
	var result = make(Metrics, len(m))
	copy(result, m)
	for _, item := range m {
		item.Execution = nil
	}
	return result
}

// SQL returns main view SQL
func (m Metrics) SQL() string {
	if m == nil || len(m) == 0 {
		return ""
	}
	return (m)[0].SQL()
}

// Lookup looks up view metric by name
func (m Metrics) Lookup(viewName string) *Metric {
	for _, candidate := range m {
		if candidate.View == viewName {
			return candidate
		}
	}
	return nil
}

func (m Metrics) ParametrizedSQL() []*ParametrizedSQL {
	var result []*ParametrizedSQL
	for _, metric := range m {
		result = append(result, metric.ParametrizedSQL()...)
	}
	return result
}
