package sqlx

type Loader interface {
	Flusher
	Insert(tableName string, data interface{}) error
}
