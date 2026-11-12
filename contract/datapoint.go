package contract

// Datapoint returns the explicit SQLX column when nonempty, or the Go field
// name otherwise. A nil field returns an empty string.
func (f *Field) Datapoint() string {
	if f == nil {
		return ""
	}
	if column := f.Column(); column != "" {
		return column
	}
	return f.Name
}
