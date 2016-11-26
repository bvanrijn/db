package db

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Record represents a database record
type Record struct {
	ID   int
	URL  string
	Tags []string
}

// Database holds Records
type Database struct {
	Records []Record
}

// Add a Record to a Database
func (database *Database) Add(record Record) {
	database.Records = append(database.Records, record)
}

// Search a Database
func (database *Database) Search(searchTerm string) []Record {
	var result []Record

	for _, record := range database.Records {
		for _, tag := range record.Tags {
			if tag == searchTerm {
				result = append(result, record)
				break
			}
		}
	}

	return result
}

// Save a Database to disk
func (database *Database) Save(path string) {
	dump, err := json.Marshal(database)
	if err != nil {
		log.Printf("%s: %s\n", path, err)
	}

	err = ioutil.WriteFile(path, dump, 0644)
	if err != nil {
		log.Printf("%s: %s\n", path, err)
	}
}

// Load a Database from file
func (database *Database) Load(path string) Database {
	var dump Database

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("%s: %s\n", path, err)
	}

	err = json.Unmarshal(dat, &dump)
	if err != nil {
		log.Printf("%s: %s\n", path, err)
	}

	return dump
}

func api(w http.ResponseWriter, r *http.Request, database *Database) {
	query := r.URL.Query()
	action := query.Get("action")

	switch action {
	case "search":
		q := query.Get("q")
		results := database.Search(q)
		resultsJSON, err := json.Marshal(results)
		if err != nil {
			log.Println(err)
		}
		io.WriteString(w, string(resultsJSON))
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad request")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "The DB sever appears to be working correctly.")
}

// Serve starts an HTTP server to query the Database
func (database *Database) Serve(port int) {
	http.HandleFunc("/", index)
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		api(w, r, database)
	})
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
	log.Printf("DB server running on port %d...\n", port)
}
