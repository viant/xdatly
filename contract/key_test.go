package contract

import (
	"reflect"
	"testing"
)

func TestFieldIsPK(t *testing.T) {
	for _, tc := range []struct {
		tag  reflect.StructTag
		want bool
	}{
		{`sqlx:"ID,primaryKey"`, true}, {`sqlx:"TENANT_ID,primaryKey=true"`, true},
		{`sqlx:"ID,PRIMARYKEY"`, true}, {`sqlx:"ID, primaryKey = true "`, true},
		{`sqlx:"ID,primaryKey=false"`, false}, {`sqlx:"ID,primaryKey=true,primaryKey=false"`, false},
		{`sqlx:"ID,autoincrement,primaryKey=false"`, true}, {`sqlx:"primaryKey=true"`, true},
		{`sqlx:"ID,primaryKey=\"true\""`, true}, {`sqlx:"ID,primaryKey=\"\\x74rue\""`, true},
		{`sqlx:"ID,description=\"text,primaryKey,more\""`, false},
		{`sqlx:"ID,description='text,primaryKey',primaryKey=false"`, false},
		{`sqlx:"ID,primaryKeySuffix"`, false}, {`sqlx:"ID"`, false}, {`pk:"true"`, false}, {``, false},
	} {
		field := &Field{Name: "ID", Tag: tc.tag}
		if got := field.IsPK(); got != tc.want {
			t.Errorf("%s: got %t want %t", tc.tag, got, tc.want)
		}
		if field.Tag != tc.tag {
			t.Fatal("predicate mutated metadata")
		}
	}
	if (*Field)(nil).IsPK() {
		t.Fatal("nil field is a key")
	}
}
func TestFieldIsPKCommonTagAllocations(t *testing.T) {
	field := &Field{Name: "ID", Tag: `sqlx:"ID,primaryKey,refTable=records" required:"true"`}
	if allocations := testing.AllocsPerRun(1000, func() {
		if !field.IsPK() {
			panic("missing key")
		}
	}); allocations != 0 {
		t.Fatalf("allocations: %v", allocations)
	}
}
func BenchmarkFieldIsPK(b *testing.B) {
	field := &Field{Name: "ID", Tag: `sqlx:"ID,primaryKey,refTable=records" required:"true"`}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !field.IsPK() {
			b.Fatal("missing key")
		}
	}
}
