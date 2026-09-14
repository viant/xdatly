package destination

type Table struct {
	Connector         *string            `json:"connector,omitempty"`
	TableName         *string            `json:"tableName,omitempty"`
	TableDataset      *string            `json:"tableDataset,omitempty"`
	TableSchema       *string            `json:"tableSchema,omitempty"`
	CreateDisposition *CreateDisposition `json:"createDisposition,omitempty"`
	Template          *string            `json:"template,omitempty"`
	WriteDisposition  *WriteDisposition  `json:"writeDisposition,omitempty"`
}
