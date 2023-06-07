package sqlx

type Updater interface {
	Flusher
	Update(tableName string, data interface{}) error
}
