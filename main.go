package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type RandomData struct {
	// unexported for json
	ID          int64
	DateCreated time.Time
	LastUpdated time.Time
	// exported for json
	Title string `json:"title"`
	Text  string `json:"text"`
}

type randomDataHandler struct {
	dbFilepath string
}

func InitLogs() {
	var logsDir string = filepath.Join(".", "logs")

	err := os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logfile, err := os.Create(filepath.Join(logsDir, "logs.log"))
	if err != nil {
		panic(err)
	}
	log.SetOutput(logfile)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	helpMessage := `
	<h1>CRUD api in Go's standart library</h1>
	<ul>
		<li> (GET) <a href="/randomdata">/randomdata</a> - to get all database </li>
		<li> (GET) /randomdata/{id} - to get specific random data under corresponding id  </li>
		<li> (POST) /randomdata - to create random data</li>
		<li> (DELETE) /randomdata/{id} - to delete specified random data </li>
		<li> (PUT) /randomdata/{id} - to update random data with given id </li>
	</ul>
	`
	fmt.Fprint(w, helpMessage)
}

func newDatabaseHandler() *randomDataHandler {
	dbDirpath := filepath.Join(".", "database")
	err := os.MkdirAll(dbDirpath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	dbFilepath := filepath.Join(dbDirpath, "database.json")

	dbFile, err := os.OpenFile(dbFilepath, os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer dbFile.Close()

	log.Println("Successfully created new database handler")

	return &randomDataHandler{
		dbFilepath: dbFilepath,
	}
}

func (dbHandler *randomDataHandler) writeRandomData(newData RandomData) error {
	dbBytes, err := dbHandler.readDatabase()
	if err != nil {
		log.Println("Error reading db (writeRandomData) : ", err)
	}

	var db []RandomData

	err = json.Unmarshal(dbBytes, &db)
	if err != nil {
		log.Println("Error unmarshalling db (writeRandomData) : ", err)
	}

	db = append(db, newData)

	dbFile, err := os.OpenFile(dbHandler.dbFilepath, os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening db file (writeRandomData) : ", err)
	}
	defer dbFile.Close()

	jsonBytes, err := json.MarshalIndent(db, "", " ")
	if err != nil {
		log.Println("Error marshalling db (writeRandomData) : ", err)
	}

	dbFile.Write(jsonBytes)

	return nil
}

func (dbHandler *randomDataHandler) readDatabase() ([]byte, error) {
	dbBytes, err := os.ReadFile(dbHandler.dbFilepath)
	if err != nil {
		log.Println("Error reading db (readDatabase) : ", err)
		return nil, err
	}
	return dbBytes, nil
}

func (dbHandler *randomDataHandler) removeRandomData(id int64) error {
	dbBytes, err := dbHandler.readDatabase()
	if err != nil {
		return err
	}

	var db []RandomData
	err = json.Unmarshal(dbBytes, &db)
	if err != nil {
		return err
	}

	var counter int64 = 0
	for _, randomData := range db {
		if id == randomData.ID {
			db = append(db[:counter], db[counter+1:]...)
			err = dbHandler.writeDB(db)
			if err != nil {
				return err
			}
		}
		counter++
	}
	return nil
}

func (dbHandler *randomDataHandler) writeDB(db []RandomData) error {
	jsonEncodedDB, err := json.MarshalIndent(db, "", " ")
	if err != nil {
		return err
	}

	dbFile, err := os.OpenFile(dbHandler.dbFilepath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dbFile.Close()

	dbFile.Write(jsonEncodedDB)

	return nil
}

func (dbHandler *randomDataHandler) get(w http.ResponseWriter, r *http.Request) {
	dbBytes, err := dbHandler.readDatabase()
	if err != nil {
		log.Println("Error reading db (get) : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(dbBytes)
}

func (dbHandler *randomDataHandler) create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Got `%s` instead of `application/json`", r.Header.Get("content-type"))))
		return
	}
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading http request (create) : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var newRandomData RandomData
	err = json.Unmarshal(requestBody, &newRandomData)
	if err != nil {
		log.Printf("Error unmarshalling http request (create) : %q \n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newRandomData.DateCreated = time.Now().UTC()
	newRandomData.LastUpdated = newRandomData.DateCreated
	newRandomData.ID = time.Now().UTC().UnixNano()

	err = dbHandler.writeRandomData(newRandomData)
	if err != nil {
		log.Println("Error writing RandomData (create): ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Println("Successfuly added to db : ", newRandomData)
}

func (dbHandler *randomDataHandler) getSpecificRandomData(w http.ResponseWriter, r *http.Request) {
	givenID := strings.Split(r.URL.String(), "/")[2]

	dbBytes, err := dbHandler.readDatabase()
	if err != nil {
		log.Println("Error reading db (getSpecificRandomData) : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var db []RandomData

	err = json.Unmarshal(dbBytes, &db)
	if err != nil {
		log.Println("Error unmarshalling database (getSpecificRandomData) : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	int64GivenID, _ := strconv.ParseInt(givenID, 10, 64)
	for _, randomData := range db {
		if int64GivenID == randomData.ID {
			response, err := json.MarshalIndent(randomData, "", " ")
			if err != nil {
				log.Println("Error marshaling response(getSpecificRandomData) : ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(response)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (dbHandler *randomDataHandler) updateSpecificRandomData(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Got `%s` instead of `application/json`", r.Header.Get("content-type"))))
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading http request (create) : ", err)
		w.WriteHeader(http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()

	var givenUpdatedRandomData RandomData
	err = json.Unmarshal(requestBody, &givenUpdatedRandomData)
	if err != nil {
		log.Println("Error unmarshalling request body (updateSpecificRandomData) : ", err)
		return
	}

	givenID := strings.Split(r.URL.String(), "/")[2]

	dbBytes, err := dbHandler.readDatabase()
	if err != nil {
		log.Println("Error reading db (updateSpecificRandomData) : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var db []RandomData

	err = json.Unmarshal(dbBytes, &db)
	if err != nil {
		log.Println("Error unmarshalling database (update) : ", err)
		return
	}

	int64GivenID, _ := strconv.ParseInt(givenID, 10, 64)
	var counter int64
	for _, randomData := range db {
		if int64GivenID == randomData.ID {
			var updatedRandomData RandomData

			updatedRandomData = givenUpdatedRandomData
			updatedRandomData.ID = randomData.ID
			updatedRandomData.DateCreated = randomData.DateCreated
			updatedRandomData.LastUpdated = time.Now().UTC()

			dbHandler.removeRandomData(int64GivenID)
			dbHandler.writeRandomData(updatedRandomData)

			log.Printf("Successfully updated RandomData with id %v \n", updatedRandomData.ID)
			return
		}
		counter++
	}

	w.WriteHeader(http.StatusNotFound)
}

func (dbHandler *randomDataHandler) deleteSpecificRandomData(w http.ResponseWriter, r *http.Request) {
	givenID := strings.Split(r.URL.String(), "/")[2]

	int64GivenID, _ := strconv.ParseInt(givenID, 10, 64)

	err := dbHandler.removeRandomData(int64GivenID)
	if err != nil {
		log.Println("Error removing RandomData (deleteSpecificRandomData) : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Printf("Successfully deleted RandomData with id %v \n", int64GivenID)
}

func (dbHandler *randomDataHandler) handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		dbHandler.get(w, r)
	case "POST":
		dbHandler.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (dbHandler *randomDataHandler) handleSpecific(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		dbHandler.getSpecificRandomData(w, r)
	case "PUT":
		dbHandler.updateSpecificRandomData(w, r)
	case "DELETE":
		dbHandler.deleteSpecificRandomData(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func init() {
	InitLogs()
}

func main() {

	databaseHandler := newDatabaseHandler()

	servemux := http.NewServeMux()
	servemux.HandleFunc("/", homepage)
	servemux.HandleFunc("/randomdata", databaseHandler.handle)
	servemux.HandleFunc("/randomdata/", databaseHandler.handleSpecific)

	server := &http.Server{
		Addr:         ":8000",
		Handler:      servemux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
