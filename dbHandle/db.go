package dbhandle

import "database/sql"

type DB struct {
	*sql.DB
}

const drivername string = "sqlite3"
