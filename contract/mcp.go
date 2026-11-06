package contract

import (
	"fmt"
	"strings"
)

// SetMCPName sets the public MCP field name without changing its Go/JSON name
// or binding location. Existing aliases and other MCP options are preserved.
// As with other helpers, use the returned copy-on-write field.
func (f *Field) SetMCPName(name string) *Field {
	if f.err != nil {
		return f
	}
	name = strings.TrimSpace(name)
	if name == "" || name == "-" || strings.ContainsAny(name, ",\r\n\t") {
		return f.withError(fmt.Errorf("invalid MCP field name %q", name))
	}
	raw := f.Tag.Get("mcp")
	if strings.TrimSpace(raw) == "-" {
		return f.withError(fmt.Errorf("cannot name an MCP-hidden field; remove the mcp exclusion explicitly first"))
	}
	parts := strings.Split(raw, ",")
	same, found := true, false
	options := make([]string, 0, len(parts))
	for i, part := range parts {
		token := strings.TrimSpace(part)
		if token == "" {
			continue
		}
		key, value, pair := strings.Cut(token, "=")
		if (i == 0 && !pair) || (pair && strings.TrimSpace(key) == "name") {
			if !pair {
				value = token
			}
			found = true
			same = same && strings.TrimSpace(value) == name
			continue
		}
		options = append(options, part)
	}
	if found && same {
		return f.SetTag("mcp", raw)
	}
	value := "name=" + name
	if len(options) > 0 {
		value += "," + strings.Join(options, ",")
	}
	return f.SetTag("mcp", value)
}
