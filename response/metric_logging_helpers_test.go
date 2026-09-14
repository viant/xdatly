package response

import "testing"

func TestHideMetricsPreservesNilEntries(t *testing.T) {
	metrics := Metrics{nil, &Metric{View: "records", Executions: SQLExecutions{nil, &SQLExecution{SQL: "SELECT secret", Args: []interface{}{1}}}}}
	hidden := metrics.HideMetrics()
	if spans := hidden.ToSpans(nil); len(spans) != 2 {
		t.Fatalf("span projection=%+v", spans)
	}
	if len(hidden) != 2 || hidden[0] != nil || hidden[1].Executions[0] != nil || hidden[1].Executions[1].SQL != "" || hidden[1].Executions[1].Args != nil {
		t.Fatalf("hidden=%+v", hidden)
	}
	if metrics[1].Executions[1].SQL != "SELECT secret" || len(metrics[1].Executions[1].Args) != 1 {
		t.Fatal("native data mutated")
	}
}
