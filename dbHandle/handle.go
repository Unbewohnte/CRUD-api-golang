package dbhandle

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		var dateCreated int64
		var lastUpdated int64

		rows.Scan(&id, &title, &text, &dateCreated, &lastUpdated)

		var randomData = randomdata.RandomData{
			ID:          id,
			Title:       title,
			Text:        text,
			DateCreated: dateCreated,
			LastUpdated: lastUpdated,
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

	// there should be only one row, because we looked for a specific ID, which is a primary key
	for row.Next() {
		var id uint
		var title string
		var text string
		var dateCreated int64
		var lastUpdated int64

		row.Scan(&id, &title, &text, &dateCreated, &lastUpdated)

		var randomData = randomdata.RandomData{
			ID:          id,
			Title:       title,
			Text:        text,
			DateCreated: dateCreated,
			LastUpdated: lastUpdated,
		}

		return &randomData, nil
	}

	return nil, nil
}

// Delete `RandomData` from db with given id
func (db *DB) DeleteSpecific(id uint) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id=%d", tableName, id))
	if err != nil {
		return err
	}

	return nil
}

// Replace `Randomdata` in db with given id with a new one.
// Does not affect `DateCreated`
func (db *DB) UpdateSpecific(id uint, newRD randomdata.RandomData) error {
	// `DateCreated` won`t be changed because we`re patching an already existing thing
	_, err := db.Exec(fmt.Sprintf("UPDATE %s SET title='%s', text='%s', last_updated=%d WHERE id=%d", tableName, newRD.Title, newRD.Text, newRD.LastUpdated, id))
	if err != nil {
		return err
	}

	return nil
}

// Create a new `RandomData` row in db
func (db *DB) Create(rd randomdata.RandomData) error {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (title, text, date_created, last_updated) VALUES (?, ?, ?, ?)", tableName), rd.Title, rd.Text, rd.DateCreated, rd.LastUpdated)
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

		randomData, err := randomdata.FromJson(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// create date created, last updated
		randomData.DateCreated = time.Now().UTC().Unix()
		randomData.LastUpdated = time.Now().UTC().Unix()

		err = db.Create(*randomData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Could not create a row: %s", err)
			return
		}

		w.WriteHeader(http.StatusAccepted)

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

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
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
			return
		}

		w.Write(rdJsonBytes)

	case http.MethodDelete:
		err := db.DeleteSpecific(uint(providedID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Could not delete from db: %s\n", err)
			return
		}

	case http.MethodPatch:
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

		randomData, err := randomdata.FromJson(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// create date created, last updated
		randomData.LastUpdated = time.Now().UTC().Unix()

		err = db.UpdateSpecific(uint(providedID), *randomData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Could not update a RandomData: %s\n", err)
			return
		}

		w.WriteHeader(http.StatusAccepted)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
