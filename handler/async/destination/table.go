package destination

type Table struct {
	Connector         *string            `json:",omitempty"`
	TableName         *string            `json:",omitempty"`
	TableDataset      *string            `json:",omitempty"`
	Schema            *string            `json:",omitempty"`
	CreateDisposition *CreateDisposition `json:",omitempty"`
	Template          *string            `json:",omitempty"`
	WriteDisposition  *WriteDisposition  `json:",omitempty"`
}
