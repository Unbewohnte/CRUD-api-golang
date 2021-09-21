package dbhandle

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	randomdata "github.com/Unbewohnte/crud-api/randomData"
	_ "github.com/mattn/go-sqlite3"
)

// Collect and return all rows in db
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

		contents = append(contents, &randomData)
	}

	return contents, nil
}

// Get `RandomData` from db with given id
func (db *DB) GetSpecific(id uint) (*randomdata.RandomData, error) {
	row, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE id=%d", tableName, id))
	if err != nil {
		return nil, err
	}

	// there should be only one row, because we looked for a specific ID
	for row.Next() {
		var id uint
		var title string
		var text string

		row.Scan(&id, &title, &text)

		var randomData = randomdata.RandomData{
			ID:    id,
			Title: title,
			Text:  text,
		}

		return &randomData, nil
	}

	return nil, nil
}

// Delete `RandomData` from db with given id
func (db *DB) DeleteSpecific(id uint) error {
	return nil
}

// Edit `Randomdata` from db with given id
func (db *DB) PatchSpecific(id uint) error {
	return nil
}

// Create a new `RandomData` row in db
func (db *DB) Create(rd randomdata.RandomData) error {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (title, text) VALUES (?, ?)", tableName), rd.Title, rd.Text)
	if err != nil {
		return err
	}

	return nil
}

// Handler function for all `RandomData`s in database
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
		return

	case http.MethodGet:
		randomDatas, err := db.GetEverything()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Could not retrieve db contents: %s\n", err)
			return
		}
		randomDatasJsonBytes, err := randomdata.ToJsonAll(randomDatas, true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Could not convert to json: %s\n", err)
			return
		}

		w.Write(randomDatasJsonBytes)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

// Handler function for a specific `RandomData` in database
func (db *DB) HandleSpecificWeb(w http.ResponseWriter, r *http.Request) {
	providedIDstr := strings.Split(r.RequestURI, "/")[2]

	providedID, err := strconv.ParseUint(providedIDstr, 10, 32)
	if err != nil {
		// most likely a bad id
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		randomData, err := db.GetSpecific(uint(providedID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Could not retrieve a specific RandomData: %s\n", err)
			return
		}

		rdJsonBytes, err := randomData.ToJson()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Could not convert RandomData to Json: %s\n", err)
		}

		w.Write(rdJsonBytes)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}
