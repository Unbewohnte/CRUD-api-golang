package dbhandle

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	randomdata "github.com/Unbewohnte/crud-api/randomData"
	_ "github.com/mattn/go-sqlite3"
)

func (db *DB) GetEverything() ([]*randomdata.RandomData, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}

	var contents []*randomdata.RandomData
	for rows.Next() {
		var id uint
		var title string
		var text string
		rows.Scan(&id, &title, &text)

		var randomData = randomdata.RandomData{
			ID:    id,
			Title: title,
			Text:  text,
		}

		fmt.Println(id, title, text)

		contents = append(contents, &randomData)
	}

	return contents, nil
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
	_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (title, text) VALUES (?, ?)", tableName), rd.Title, rd.Text)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) HandleGlobalWeb(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		if r.Header.Get("content-type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data randomdata.RandomData
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = db.Create(data)
		if err != nil {
			log.Printf("Could not create a row: %s", err)
		}
		w.WriteHeader(http.StatusAccepted)

	case http.MethodGet:
		data, err := db.GetEverything()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Could not retrieve db contents: %s\n", err)
			return
		}
		w.Write()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (db *DB) HandleSpecificWeb(w http.ResponseWriter, r *http.Request) {

}
