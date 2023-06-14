package sqlx

type Loader interface {
	Flusher
	Load(tableName string, data interface{}) error
}
