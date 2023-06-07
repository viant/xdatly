package sqlx

type Inserter interface {
	Flusher
	Insert(tableName string, data interface{}) error
}
