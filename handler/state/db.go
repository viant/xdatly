package state

import "database/sql"

type dBProviderKey string

var DBProviderKey dBProviderKey = "dbProvider"

// DBProvider represents a database provider
type DBProvider func(name string) (*sql.DB, error)
