package exec

import "time"

// Complete records the canonical invocation boundary. Existing native trace
// spans are completed; no default trace or span is created by this operation.
func (c *Context) Complete(end time.Time, err error) {
	if c == nil {
		return
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	c.ElapsedMs = int(end.Sub(c.StartTime).Milliseconds())
	if c.Trace != nil && len(c.Trace.Spans) > 0 && c.Trace.Spans[0] != nil {
		root := c.Trace.Spans[0]
		root.EndTime = end
		root.SetStatus(err)
	}
}
