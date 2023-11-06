package tjson

type (
	//Column represent column header
	Column struct {
		Name string `json:",omitempty"`
		Type string `json:",omitempty"`
	}

	//Value represents tabular value
	Value string

	//Record represents tabular record
	Record []Value

	//Records represents tabular records
	Records []Record

	//Columns represents columns
	Columns []*Column

	//Tabular represents a tabular holder
	Tabular struct {
		Columns Columns `json:",omitempty"`
		Records Records `json:",omitempty"`
	}
)
