package predicate

import "strings"

var TagName = "predicate"

type (

	//Tag represents predicate tag
	Tag struct {
		Exclusion bool
		Name      string
		Args      []string
	}
)

func (t *Tag) NormalizedName(name string) string {
	if t.Name != "" {
		return t.Name
	}
	toLower := strings.ToLower(name)
	if strings.HasSuffix(toLower, "excl") {
		t.Name = name[:len(name)-4]
	} else if strings.HasSuffix(toLower, "exclusion") {
		t.Name = name[:len(name)-9]
	} else if strings.HasSuffix(toLower, "incl") {
		t.Name = name[:len(name)-4]
	} else if strings.HasSuffix(toLower, "inclusion") {
		t.Name = name[:len(name)-9]
	}
	if t.Name == "" {
		t.Name = name
	}
	return t.Name
}

// ParseTag parses predicate tag
func ParseTag(tag string) *Tag {
	ret := &Tag{}
	if tag == "" {
		return ret
	}
	elements := strings.Split(tag, ";")
	for _, element := range elements {
		pair := strings.Split(element, "=")
		switch strings.ToLower(pair[0]) {
		case "name":
			if len(pair) == 2 {
				ret.Name = pair[1]
			}
		case "exclusion":
			ret.Exclusion = true
		case "args":
			if len(pair) == 2 {
				ret.Args = strings.Split(pair[1], ",")
			}
		default:
			if len(pair) == 1 {
				ret.Name = pair[0]
			}
		}
	}
	return ret
}
