package docs

import "testing"

func TestSourceCloneAndOrderedSubstitution(t *testing.T) {
	source := Source{DocURLs: []string{"a"}, Substitutes: map[string]string{"AB": "long", "A": "short", "ZZ": "$AA", "AA": "first"}}
	for i := 0; i < 100; i++ {
		if got := source.Expand("$AB ${A} $ZZ"); got != "long short shortA" {
			t.Fatal(got)
		}
	}
	copy := source.Clone()
	copy.DocURLs[0] = "b"
	copy.Substitutes["A"] = "changed"
	if source.DocURLs[0] != "a" || source.Substitutes["A"] != "short" {
		t.Fatal("clone aliases source")
	}
}

func TestSourceOverlayPreservesRuleWithAuthoredGlobals(t *testing.T) {
	base := Source{GlobalURLs: []string{"old.yaml"}, DocURL: "rule.yaml", Substitutes: map[string]string{"A": "old"}}
	result := base.Overlay(Source{GlobalURLs: []string{"new.yaml"}, Substitutes: map[string]string{"A": "new"}})
	if result.DocURL != "rule.yaml" || result.GlobalURLs[0] != "new.yaml" || result.Substitutes["A"] != "new" {
		t.Fatal(result)
	}
	if base.GlobalURLs[0] != "old.yaml" || base.Substitutes["A"] != "old" {
		t.Fatal("source overlay mutated base")
	}
}
