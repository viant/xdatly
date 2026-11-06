package contract

import (
	"reflect"
	"testing"
)

func TestSetMCPNamePreservesOtherMetadata(t *testing.T) {
	for _, tc := range []struct {
		tag  reflect.StructTag
		want string
	}{
		{`json:"pricingModel" choice:"0:fixed" choice:"10:usage"`, `name=BillingModel`},
		{`mcp:"OldName,aliases=old|legacy,custom=value" json:"pricingModel"`, `name=BillingModel,aliases=old|legacy,custom=value`},
		{`mcp:"name=OldName,aliases=old;legacy"`, `name=BillingModel,aliases=old;legacy`},
		{`mcp:"OldName,name=OtherName,aliases=legacy"`, `name=BillingModel,aliases=legacy`},
		{`mcp:"name=OldName,name=OtherName,custom=value"`, `name=BillingModel,custom=value`},
		{`mcp:"aliases=legacy"`, `name=BillingModel,aliases=legacy`},
	} {
		t.Run(string(tc.tag), func(t *testing.T) {
			original := &Field{Name: "PricingModel", Tag: tc.tag}
			result := original.SetMCPName("BillingModel")
			if result.Err() != nil || result.Tag.Get("mcp") != tc.want {
				t.Fatalf("%s: %v", result.Tag, result.Err())
			}
			if original.Tag != tc.tag || result.Name != "PricingModel" || result.Tag.Get("json") != original.Tag.Get("json") {
				t.Fatal("unrelated metadata changed")
			}
			if result == original {
				t.Fatal("first edit did not copy")
			}
			if err := result.Validate(reflect.TypeFor[int]()); err != nil {
				t.Fatal(err)
			}
			if result.SetDescription("Pricing choice") != result {
				t.Fatal("chain copied twice")
			}
		})
	}
}
func TestSetMCPNameNoOpAndErrors(t *testing.T) {
	for _, tag := range []reflect.StructTag{`mcp:"BillingModel,aliases=legacy"`, `mcp:"aliases=legacy,name=BillingModel"`} {
		field := &Field{Name: "PricingModel", Tag: tag}
		if field.SetMCPName("BillingModel") != field {
			t.Fatal("no-op copied field")
		}
	}
	for _, name := range []string{"", " ", "-", "name,aliases=other", "bad\nname"} {
		original := &Field{Name: "PricingModel", Tag: `json:"pricingModel"`}
		result := original.SetMCPName(name)
		if result.Err() == nil || original.Err() != nil || original.Tag != `json:"pricingModel"` {
			t.Fatalf("name %q: %+v", name, result)
		}
	}
	hidden := (&Field{Name: "PricingModel"}).Hide()
	original := hidden.Clone()
	if hidden.SetMCPName("BillingModel").Err() == nil {
		t.Fatal("MCP exclusion silently cleared")
	}
	if hidden.Tag != original.Tag {
		t.Fatal("hidden metadata changed")
	}
}
