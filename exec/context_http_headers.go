package exec

import (
	"net/http"
	"os"
	"strings"
)

const trackingHeaderEnvKey = "XDATLY_TRACING_HEADER"

// HTTP request-header parsing is compatibility surface for current transport
// integration. It remains isolated here so the generic exec Context contract
// does not take ownership of header parsing or incoming trace propagation.
func (c *Context) setHeader(header http.Header) {
	c.Header = make(map[string]string)
	trackingHeaderKey := os.Getenv(trackingHeaderEnvKey)
	if trackingHeaderKey == "" {
		trackingHeaderKey = "xtraceid"
	}
	trackingHeaderKey = strings.ReplaceAll(strings.ToLower(trackingHeaderKey), "-", "")
	for k := range header {
		lowerKey := strings.ToLower(k)
		if strings.Contains(lowerKey, "auth") {
			continue
		}
		if lowerKey == "traceparent" {
			if traceID, parentSpanID, ok := parseTraceparent(header.Get(k)); ok {
				c.TraceID = traceID
				c.parentSpanID = parentSpanID
				continue
			}
		}
		if trackingHeaderKey == strings.ReplaceAll(lowerKey, "-", "") {
			c.TraceID = header.Get(k)
			continue
		}
		c.Header[k] = header.Get(k)
	}
}

func parseTraceparent(value string) (string, *string, bool) {
	parts := strings.Split(strings.TrimSpace(value), "-")
	if len(parts) != 4 {
		return "", nil, false
	}
	traceID := strings.TrimSpace(parts[1])
	spanID := strings.TrimSpace(parts[2])
	if len(traceID) != 32 || len(spanID) != 16 {
		return "", nil, false
	}
	parentSpanID := spanID
	return traceID, &parentSpanID, true
}
