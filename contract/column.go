package contract

import (
	"strings"

	"github.com/viant/tagly/tags"
)

// Column returns the explicit SQLX column metadata, including any | aliases.
// Missing SQLX tags, a selected "-" column, and relation-only fields return "".
// It does not infer a column from Name or apply namespace/case conversion.
// Ordinary unquoted tags require no allocation; complex options use Tagly's
// native tokenizer. Results always reflect the current Tag, without caching.
func (f *Field) Column() string {
	if f == nil {
		return ""
	}
	raw := f.Tag.Get("sqlx")
	if raw == "" || raw == "-" {
		return ""
	}
	if _, relation := f.Tag.Lookup("on"); relation {
		return ""
	}
	column, options := tags.Values(raw).Name()
	if options == "" {
		return explicitColumn(column)
	}
	// Quoted and scoped values can contain commas or apparent column options.
	// Let the owning tokenizer interpret those rather than duplicating its grammar.
	if strings.ContainsAny(string(options), "\"'`(){}[]") {
		_ = options.MatchPairs(func(key, value string) error {
			if columnOption(key) {
				column = strings.TrimSpace(value)
			}
			return nil
		})
		return explicitColumn(column)
	}
	remaining := string(options)
	for remaining != "" {
		var option string
		option, remaining, _ = strings.Cut(remaining, ",")
		key, value, _ := strings.Cut(option, "=")
		if columnOption(key) {
			column = strings.TrimSpace(value)
		}
	}
	return explicitColumn(column)
}

func explicitColumn(column string) string {
	if column == "-" {
		return ""
	}
	return column
}

func columnOption(key string) bool {
	key = strings.TrimSpace(key)
	return strings.EqualFold(key, "column") || strings.EqualFold(key, "name")
}
