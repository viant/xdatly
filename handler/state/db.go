package state

import "database/sql"

// DBProvider represents a database provider
type DBProvider func(name string) (*sql.DB, error)
