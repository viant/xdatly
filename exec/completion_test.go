package exec

import (
	"github.com/viant/xdatly/tracing"
	"testing"
	"time"
)

func TestCompleteKeepsNativeTimingWithoutDefaultTrace(t *testing.T) {
	start := time.Now()
	c := New(WithStartTime(start))
	c.Complete(start.Add(time.Second), nil)
	if c.ElapsedMs != 1000 || c.Trace != nil {
		t.Fatal("completion created tracing")
	}
	c.Trace = &tracing.Trace{Spans: []*tracing.Span{{StartTime: start}}}
	end := start.Add(2 * time.Second)
	c.Complete(end, nil)
	if !c.Trace.Spans[0].EndTime.Equal(end) || !c.Trace.Spans[0].StartTime.Equal(start) || c.Trace.Spans[0].Status.Code != tracing.StatusOK {
		t.Fatal("native interval lost")
	}
}
