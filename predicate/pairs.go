package predicate

import (
	"fmt"
	"strings"
)

// Pairs extracts optional alternating name/value arguments. Names must be
// nonempty, unique strings. Values retain their Go types. Handlers explicitly
// allowlist names, validate types and seed their own configuration attributes;
// this helper never writes fields or overrides bound identity/dependencies.
// Pass only trusted authored configuration, never request-supplied options.
func Pairs(args ...any) (map[string]any, error) {
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("predicate options require name/value pairs, got %d arguments", len(args))
	}
	result := make(map[string]any, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		name, ok := args[i].(string)
		if !ok || strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("predicate option %d requires a nonempty string name", i/2)
		}
		if _, exists := result[name]; exists {
			return nil, fmt.Errorf("duplicate predicate option %q", name)
		}
		result[name] = args[i+1]
	}
	return result, nil
}
