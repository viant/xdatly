package response

// ParametrizedSQL helpers are convenience helpers over the shared response
// payloads. They remain public for current compatibility and inspection paths,
// but they are not the core payload contract itself.

func (m *Metric) ParametrizedSQL() []*ParametrizedSQL {
	var result []*ParametrizedSQL
	if m.Executions == nil {
		return result
	}
	for _, tmpl := range m.Executions {
		result = append(result, &ParametrizedSQL{Query: tmpl.SQL, Args: tmpl.Args})
	}
	return result
}

func (m Metrics) ParametrizedSQL() []*ParametrizedSQL {
	var result []*ParametrizedSQL
	for _, metric := range m {
		result = append(result, metric.ParametrizedSQL()...)
	}
	return result
}

func (m *Metric) SQL() string {
	if m.Executions != nil && len(m.Executions) > 0 {
		return ExpandSQL(m.Executions[0].SQL, m.Executions[0].Args)
	}
	return ""
}

func (m Metrics) SQL() string {
	if len(m) == 0 {
		return ""
	}
	return m[0].SQL()
}
