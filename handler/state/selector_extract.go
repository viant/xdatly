package state

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ExtractQuerySelector converts a struct spec into a QuerySelector.
//
// Supported input:
//   - struct / *struct (exported fields; `selector` tag takes precedence over `json` tag)
//
// Supported keys/fields (case-insensitive; underscores ignored):
//
//	Columns, Fields, OrderBy, Offset, Limit, Page, Criteria, Placeholders
func ExtractQuerySelector(spec any) (QuerySelector, error) {
	var dest QuerySelector
	if spec == nil {
		return dest, nil
	}
	switch spec.(type) {
	case QuerySelector, *QuerySelector, NamedQuerySelector, *NamedQuerySelector:
		return QuerySelector{}, fmt.Errorf("unsupported selector spec: %T (expected struct or *struct)", spec)
	}

	value := reflect.ValueOf(spec)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return dest, nil
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return QuerySelector{}, fmt.Errorf("unsupported selector spec: %T (expected struct or *struct)", spec)
	}
	applyQuerySelectorFromStruct(&dest, value)
	return dest, nil
}

// ExtractNamedQuerySelector converts a spec into a NamedQuerySelector with the supplied view name.
func ExtractNamedQuerySelector(viewName string, spec any) (*NamedQuerySelector, error) {
	if viewName == "" {
		return nil, fmt.Errorf("viewName was empty")
	}
	qs, err := ExtractQuerySelector(spec)
	if err != nil {
		return nil, err
	}
	return &NamedQuerySelector{Name: viewName, QuerySelector: qs}, nil
}

func applyQuerySelectorFromStruct(dest *QuerySelector, value reflect.Value) {
	if dest == nil {
		return
	}
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}
		name := selectorFieldName(field)
		if name == "" {
			continue
		}
		setQuerySelectorField(dest, name, value.Field(i))
	}
}

func selectorFieldName(field reflect.StructField) string {
	if tag := strings.TrimSpace(field.Tag.Get("selector")); tag != "" {
		return strings.Split(tag, ",")[0]
	}
	if tag := strings.TrimSpace(field.Tag.Get("parameter")); tag != "" && tag != "-" {
		// datly-style parameter tag: `parameter:"<name>,kind=query,in=_page"`
		if name := strings.Split(tag, ",")[0]; name != "" {
			return name
		}
	}
	if tag := strings.TrimSpace(field.Tag.Get("json")); tag != "" && tag != "-" {
		if name := strings.Split(tag, ",")[0]; name != "" {
			return name
		}
	}
	return field.Name
}

func setQuerySelectorField(dest *QuerySelector, name string, value reflect.Value) {
	if dest == nil {
		return
	}
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}

	switch normalizeSelectorKey(name) {
	case "columns":
		if v, ok := asStringSlice(value); ok {
			dest.Columns = v
		}
	case "fields":
		if v, ok := asStringSlice(value); ok {
			dest.Fields = v
		}
	case "orderby":
		if v, ok := asOrderBy(value); ok {
			dest.OrderBy = v
		}
	case "offset":
		if v, ok := asInt(value); ok {
			dest.Offset = v
		}
	case "limit":
		if v, ok := asInt(value); ok {
			dest.Limit = v
		}
	case "page":
		if v, ok := asInt(value); ok {
			dest.Page = v
		}
	case "criteria":
		if v, ok := asString(value); ok {
			dest.Criteria = v
		}
	case "placeholders":
		if v, ok := asInterfaceSlice(value); ok {
			dest.Placeholders = v
		}
	}
}

func normalizeSelectorKey(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, "-", "")
	return strings.ToLower(name)
}

func asString(v reflect.Value) (string, bool) {
	if !v.IsValid() {
		return "", false
	}
	if v.Kind() == reflect.String {
		return v.String(), true
	}
	return "", false
}

func asInt(v reflect.Value) (int, bool) {
	if !v.IsValid() {
		return 0, false
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return int(v.Float()), true
	case reflect.String:
		i, err := strconv.Atoi(strings.TrimSpace(v.String()))
		return i, err == nil
	default:
		return 0, false
	}
}

func asStringSlice(v reflect.Value) ([]string, bool) {
	if !v.IsValid() {
		return nil, false
	}
	if v.Kind() == reflect.String {
		raw := strings.TrimSpace(v.String())
		if raw == "" {
			return []string{}, true
		}
		parts := strings.Split(raw, ",")
		out := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
		return out, true
	}
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, false
	}
	out := make([]string, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Interface && !elem.IsNil() {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.String {
			continue
		}
		out = append(out, elem.String())
	}
	return out, true
}

func asOrderBy(v reflect.Value) (string, bool) {
	if s, ok := asString(v); ok {
		return s, true
	}
	if fields, ok := asStringSlice(v); ok {
		return strings.Join(fields, ","), true
	}
	return "", false
}

func asInterfaceSlice(v reflect.Value) ([]interface{}, bool) {
	if !v.IsValid() {
		return nil, false
	}
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, false
	}
	out := make([]interface{}, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		out = append(out, v.Index(i).Interface())
	}
	return out, true
}
