package contract

import (
	"strconv"
	"strings"
)

// IsPK reports explicit sqlx primaryKey/autoincrement metadata. It does not
// infer keys from Go field names. Ordinary tags require no allocation.
func (f *Field) IsPK() bool {
	if f == nil {
		return false
	}
	raw := f.Tag.Get("sqlx")
	comma, equal := strings.IndexByte(raw, ','), strings.IndexByte(raw, '=')
	if equal < 0 || (comma >= 0 && equal > comma) {
		if comma < 0 {
			return false
		}
		raw = raw[comma+1:] // sqlx's positional column name
	}
	primary := false
	for raw != "" {
		token, rest := keyTagToken(raw)
		raw = rest
		key, value, assigned := strings.Cut(token, "=")
		switch {
		case strings.EqualFold(strings.TrimSpace(key), "autoincrement"):
			return true // sqlx treats the annotation's presence as a primary key.
		case strings.EqualFold(strings.TrimSpace(key), "primaryKey"):
			if !assigned {
				primary = true
				continue
			}
			value = strings.TrimSpace(value)
			if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '`' && value[len(value)-1] == '`')) {
				if strings.IndexByte(value, '\\') < 0 {
					value = value[1 : len(value)-1]
				} else if decoded, err := strconv.Unquote(value); err == nil {
					value = decoded
				}
			}
			value = strings.TrimSpace(value)
			primary = value == "" || value == "true"
		}
	}
	return primary
}

// Split sqlx options without mistaking a comma inside a quoted value for a flag.
func keyTagToken(raw string) (string, string) {
	var quote byte
	for i := 0; i < len(raw); i++ {
		ch := raw[i]
		if quote != 0 {
			if ch == '\\' && quote != '`' && i+1 < len(raw) {
				i++
				continue
			}
			if ch == quote {
				quote = 0
			}
			continue
		}
		switch ch {
		case '\'', '"', '`':
			quote = ch
		case ',':
			return raw[:i], raw[i+1:]
		}
	}
	return raw, ""
}
