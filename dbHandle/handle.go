package dbhandle

import (
	"net/http"

	randomdata "github.com/Unbewohnte/crud-api/randomData"
	_ "github.com/mattn/go-sqlite3"
)

func (db *DB) GetEverything() ([]*randomdata.RandomData, error) {
	return nil, nil
}

func (db *DB) GetSpecific() (*randomdata.RandomData, error) {
	return nil, nil
}

func (db *DB) DeleteSpecific() error {
	return nil
}

func (db *DB) PatchSpecific() error {
	return nil
}

func (db *DB) Create(rd randomdata.RandomData) error {
	_, err := db.Exec("INSERT INTO randomdata (title, text) VALUES (?, ?)", rd.Title, rd.Text)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) HandleSpecificWeb(w http.ResponseWriter, r *http.Request) {

}

func (db *DB) HandleGlobalWeb(w http.ResponseWriter, r *http.Request) {

}
