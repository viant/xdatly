package sqlx

type Deleter interface {
	Flusher
	Delete(tableName string, data interface{}) error
}
