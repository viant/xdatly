package differ

type (
	//ChangeLog represents a change log
	ChangeLog struct {
		Changes []*Change
	}
)

// ToChangeRecords converts changeLog to change records
func (l *ChangeLog) ToChangeRecords(source, id, userID string) []*ChangeRecord {
	var result []*ChangeRecord
	for _, change := range l.Changes {
		result = append(result, &ChangeRecord{
			Source:   source,
			SourceID: id,
			UserID:   userID,
			Path:     change.Path.String(),
			Change:   string(change.Type),
			From:     change.From,
			To:       change.To,
		})
	}
	return result
}
