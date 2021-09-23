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
	port      *uint  = flag.Uint("port", 8000, "Specifies a port on which the helping page will be served")
	dbname    string = "database.db"
	tableName string = "randomdata"
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
	db, err := dbhandle.CreateLocalDB(dbname, tableName)
	if err != nil {
		log.Fatalf("error setting up a database: %s", err)
	}
	log.Printf("Created %s db\n", tableName)

	// set up patterns and handlers
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
