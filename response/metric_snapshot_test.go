package response

import "testing"

func TestHideMetricsPreservesNativeExecution(t *testing.T) {
	source := Metrics{{Executions: SQLExecutions{{SQL: "SELECT ?", Args: []any{"private"}}, nil}}}
	hidden := source.HideMetrics()
	if source[0].Executions[0].SQL != "SELECT ?" || source[0].Executions[0].Args[0] != "private" {
		t.Fatal("native capture mutated")
	}
	if hidden[0].Executions[0].SQL != "" || hidden[0].Executions[0].Args != nil || hidden[0].Executions[1] != nil {
		t.Fatal("projection not hidden")
	}
}
