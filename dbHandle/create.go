package dbhandle

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// full path to the local database that has been created by `CreateLocalDB`
var dbpath string

// Creates a local database file in the same directory
// as executable if does not exist already. Creates a table "randomdata" with
// such structure: (id PRIMARY KEY, title TEXT, text TEXT), which
// represents underlying fields of `RandomData`
func CreateLocalDB(dbName string) (*DB, error) {
	// double check if dbName is actually just a name, not a path
	dbName = filepath.Base(dbName)

	executablePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(executablePath)

	dbpath = filepath.Join(exeDir, dbName)

	// create db if does not exist
	dbfile, err := os.OpenFile(dbpath, os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	dbfile.Close()

	// create table that suits `RandomData` struct
	db, err := sql.Open(drivername, dbpath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS randomdata (id INTEGER PRIMARY KEY, title TEXT, text TEXT)")
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
