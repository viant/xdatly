package contract

import (
	"reflect"
	"testing"
)

func TestFieldColumn(t *testing.T) {
	for _, tc := range []struct {
		tag  reflect.StructTag
		want string
	}{
		{``, ""}, {`json:"id"`, ""}, {`sqlx:""`, ""}, {`sqlx:"-"`, ""},
		{`sqlx:"ORDER_ID"`, "ORDER_ID"},
		{`sqlx:"-,primaryKey"`, ""},
		{`sqlx:"ID,column=-"`, ""}, {`sqlx:"ID,name=-"`, ""},
		{`sqlx:"ID,column=\"-\""`, ""},
		{`sqlx:"ORDER_ID,primaryKey,required=true"`, "ORDER_ID"},
		{`sqlx:"column=ORDER_ID"`, "ORDER_ID"},
		{`sqlx:"name=ORDER_ID"`, "ORDER_ID"},
		{`sqlx:",column=ORDER_ID"`, "ORDER_ID"},
		{`sqlx:"OLD,name=FIRST,column=LAST"`, "LAST"},
		{`sqlx:"OLD, CoLuMn = LAST "`, "LAST"},
		{`sqlx:"OLD,column="`, ""}, {`sqlx:"OLD,column"`, ""},
		{`sqlx:"ID|LEGACY_ID,primaryKey"`, "ID|LEGACY_ID"},
		{`sqlx:"ID,column=\"Order ID\""`, "Order ID"},
		// Preserve the pinned SQLX/Tagly parser's handling of quoted commas.
		{`sqlx:"ID,column=\"A,B\""`, `"A`},
		{`sqlx:"ID,column=\"\\x49D\""`, "ID"},
		{`sqlx:"ID,description='text,column=WRONG',required=true"`, "ID"},
		{`sqlx:"ID,description=\"text,column=WRONG\",name=RIGHT"`, "RIGHT"},
		{`sqlx:"ID,description=(text,column=WRONG),name=RIGHT"`, "RIGHT"},
		{`sqlx:"ID,refcolumn=OTHER,namespace=X"`, "ID"},
		{`sqlx:"ID,ns=PREFIX,caseFormat=lc"`, "ID"},
		{`sqlx:"ID" on:"ID:ID"`, ""},
		// SQLX retains explicit column metadata even for transient/presence options.
		{`sqlx:"ID,transient"`, "ID"}, {`sqlx:"ID,presence"`, "ID"},
	} {
		t.Run(string(tc.tag), func(t *testing.T) {
			field := &Field{Name: "Unrelated", Tag: tc.tag}
			if got := field.Column(); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
			if field.Name != "Unrelated" || field.Tag != tc.tag || field.Err() != nil {
				t.Fatal("column lookup changed field metadata")
			}
		})
	}
	if (*Field)(nil).Column() != "" {
		t.Fatal("nil field has a column")
	}
}

func TestFieldColumnAfterTagEdit(t *testing.T) {
	field := &Field{Name: "ID", Tag: `sqlx:"OLD,primaryKey"`}
	if field.Column() != "OLD" {
		t.Fatal("missing original column")
	}
	updated := field.SetTag("sqlx", "NEW,primaryKey")
	if updated.Column() != "NEW" || updated.RemoveTag("sqlx").Column() != "" {
		t.Fatal("lookup retained stale metadata")
	}
}

func TestFieldColumnCommonTagsAllocateNothing(t *testing.T) {
	for _, tag := range []reflect.StructTag{`sqlx:"ID"`, `sqlx:"ID,primaryKey,required=true"`, `sqlx:"column=ID"`, `sqlx:"OLD,name=ID"`} {
		field := &Field{Tag: tag}
		if count := testing.AllocsPerRun(1000, func() {
			if field.Column() != "ID" {
				t.Fatal("missing column")
			}
		}); count != 0 {
			t.Fatalf("%s: %v allocations", tag, count)
		}
	}
}

var benchmarkColumn string

func BenchmarkFieldColumn(b *testing.B) {
	for _, tag := range []reflect.StructTag{`sqlx:"ID"`, `sqlx:"ID,primaryKey,required=true"`, `sqlx:"column=ID"`, `sqlx:"ID,column=\"Order ID\""`} {
		b.Run(string(tag), func(b *testing.B) {
			field := &Field{Tag: tag}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchmarkColumn = field.Column()
			}
		})
	}
}
