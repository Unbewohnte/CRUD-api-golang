package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	dbhandle "github.com/Unbewohnte/crud-api/dbHandle"
	"github.com/Unbewohnte/crud-api/logs"
)

var (
	port   *uint  = flag.Uint("port", 8080, "Specifies a port on which the helping page will be served")
	dbname string = "database.db"
)

func init() {
	// set up logs, parse flags
	err := logs.SetUp()
	if err != nil {
		panic(err)
	}

	flag.Parse()
}

func main() {
	// create a local db file
	db, err := dbhandle.CreateLocalDB(dbname)
	if err != nil {
		log.Fatalf("error setting up a database: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", helpPage)
	mux.HandleFunc("/randomdata", db.HandleGlobalWeb)
	mux.HandleFunc("/randomdata/", db.HandleSpecificWeb)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Printf("Starting on port %d\n", *port)
	fmt.Printf("Starting on port %d\n", *port)

	log.Fatal(server.ListenAndServe())
}

func helpPage(w http.ResponseWriter, r *http.Request) {
	helpMessage := `
	<h1>CRUD api</h1>
	<ul>
		<li> (GET) <a href="/randomdata">/randomdata</a> - to get all info from database (obviously a bad idea in serious projects ᗜˬᗜ) </li>
		<li> (GET) /randomdata/{id} - to get specific random data under corresponding id </li>
		<li> (POST) /randomdata - to create random data </li>
		<li> (DELETE) /randomdata/{id} - to delete specified random data </li>
		<li> (PUT) /randomdata/{id} - to update random data with given id </li>
	</ul>
	`
	fmt.Fprint(w, helpMessage)
}

// func (dbHandler *randomDataHandler) updateSpecificRandomData(w http.ResponseWriter, r *http.Request) {
// 	if r.Header.Get("content-type") != "application/json" {
// 		w.WriteHeader(http.StatusUnsupportedMediaType)
// 		w.Write([]byte(fmt.Sprintf("Got `%s` instead of `application/json`", r.Header.Get("content-type"))))
// 		return
// 	}

// 	requestBody, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		log.Println("Error reading http request (create) : ", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		r.Body.Close()
// 		return
// 	}
// 	defer r.Body.Close()

// 	var givenUpdatedRandomData RandomData
// 	err = json.Unmarshal(requestBody, &givenUpdatedRandomData)
// 	if err != nil {
// 		log.Println("Error unmarshalling request body (updateSpecificRandomData) : ", err)
// 		return
// 	}

// 	givenID := strings.Split(r.URL.String(), "/")[2]

// 	dbBytes, err := dbHandler.readDatabase()
// 	if err != nil {
// 		log.Println("Error reading db (updateSpecificRandomData) : ", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	var db []RandomData

// 	err = json.Unmarshal(dbBytes, &db)
// 	if err != nil {
// 		log.Println("Error unmarshalling database (update) : ", err)
// 		return
// 	}

// 	int64GivenID, _ := strconv.ParseInt(givenID, 10, 64)
// 	var counter int64
// 	for _, randomData := range db {
// 		if int64GivenID == randomData.ID {
// 			var updatedRandomData RandomData

// 			updatedRandomData = givenUpdatedRandomData
// 			updatedRandomData.ID = randomData.ID
// 			updatedRandomData.DateCreated = randomData.DateCreated
// 			updatedRandomData.LastUpdated = time.Now().UTC()

// 			dbHandler.removeRandomData(int64GivenID)
// 			dbHandler.writeRandomData(updatedRandomData)

// 			log.Printf("Successfully updated RandomData with id %v \n", updatedRandomData.ID)
// 			return
// 		}
// 		counter++
// 	}

// 	w.WriteHeader(http.StatusNotFound)
// }

// func (dbHandler *randomDataHandler) deleteSpecificRandomData(w http.ResponseWriter, r *http.Request) {
// 	givenID := strings.Split(r.URL.String(), "/")[2]

// 	int64GivenID, _ := strconv.ParseInt(givenID, 10, 64)

// 	err := dbHandler.removeRandomData(int64GivenID)
// 	if err != nil {
// 		log.Println("Error removing RandomData (deleteSpecificRandomData) : ", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)

// 	log.Printf("Successfully deleted RandomData with id %v \n", int64GivenID)
// }

// func (dbHandler *randomDataHandler) handle(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		dbHandler.get(w, r)
// 	case "POST":
// 		dbHandler.create(w, r)
// 	default:
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}
// }

// func (dbHandler *randomDataHandler) handleSpecific(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		dbHandler.getSpecificRandomData(w, r)
// 	case "PUT":
// 		dbHandler.updateSpecificRandomData(w, r)
// 	case "DELETE":
// 		dbHandler.deleteSpecificRandomData(w, r)
// 	default:
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}
// }
