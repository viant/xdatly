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
