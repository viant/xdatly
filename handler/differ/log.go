package differ

type (
	//ChangeLog represents a change log
	ChangeLog struct {
		Changes []*Change
	}
)

// ToChangeRecords converts changeLog to change records
func (l *ChangeLog) ToChangeRecords(options ...LogOption) []*ChangeRecord {
	var result []*ChangeRecord

	opts := &logOptions{}
	for _, opt := range options {
		opt(opts)
	}

	for _, change := range l.Changes {
		result = append(result, &ChangeRecord{
			Source:   opts.source,
			SourceID: opts.id,
			UserID:   opts.userID,
			Path:     change.Path.String(),
			Change:   string(change.Type),
			From:     change.From,
			To:       change.To,
		})
	}
	return result
}
