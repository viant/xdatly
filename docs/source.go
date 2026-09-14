package docs

import (
	"sort"
	"strings"
)

// Source declares an ordered documentation overlay in the canonical resource
// store. DocURLs takes precedence over DocURL. BaseURL is a resource directory,
// optionally qualified by its registered filesystem namespace.
type Source struct {
	GlobalURLs  []string          `json:"globalURLs,omitempty" yaml:"GlobalURLs,omitempty"`
	Substitutes map[string]string `json:"substitutes,omitempty" yaml:"Substitutes,omitempty"`
	BaseURL     string            `json:"baseURL,omitempty" yaml:"BaseURL,omitempty"`
	DocURL      string            `json:"docURL,omitempty" yaml:"DocURL,omitempty"`
	DocURLs     []string          `json:"docURLs,omitempty" yaml:"DocURLs,omitempty"`
}

func (s Source) Clone() Source {
	if s.Substitutes != nil {
		copy := make(map[string]string, len(s.Substitutes))
		for key, value := range s.Substitutes {
			copy[key] = value
		}
		s.Substitutes = copy
	}
	s.GlobalURLs = append([]string(nil), s.GlobalURLs...)
	s.DocURLs = append([]string(nil), s.DocURLs...)
	return s
}

// Expand follows original Datly's longest-name-first substitution order, with
// lexical ties to make equal-length overlapping substitutions deterministic.
func (s Source) Expand(text string) string {
	keys := make([]string, 0, len(s.Substitutes))
	for key := range s.Substitutes {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) != len(keys[j]) {
			return len(keys[i]) > len(keys[j])
		}
		return keys[i] < keys[j]
	})
	for _, key := range keys {
		text = strings.ReplaceAll(text, "${"+key+"}", s.Substitutes[key])
		text = strings.ReplaceAll(text, "$"+key, s.Substitutes[key])
	}
	return text
}

func (s Source) IsZero() bool {
	return len(s.GlobalURLs) == 0 && s.BaseURL == "" && s.DocURL == "" && len(s.DocURLs) == 0 && len(s.Substitutes) == 0
}

// Overlay preserves package source declarations unless authored DQL supplies
// that group. Dictionary content merging remains the documentation loader's job.
func (s Source) Overlay(authored Source) Source {
	result := s.Clone()
	if authored.BaseURL != "" {
		result.BaseURL = authored.BaseURL
	}
	if authored.GlobalURLs != nil {
		result.GlobalURLs = append([]string(nil), authored.GlobalURLs...)
	}
	if authored.DocURL != "" {
		result.DocURL = authored.DocURL
		result.DocURLs = nil
	}
	if len(authored.DocURLs) > 0 {
		result.DocURLs = append([]string(nil), authored.DocURLs...)
	}
	if authored.Substitutes != nil {
		if result.Substitutes == nil {
			result.Substitutes = map[string]string{}
		}
		for name, value := range authored.Substitutes {
			result.Substitutes[name] = value
		}
	}
	return result
}
