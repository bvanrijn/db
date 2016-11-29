package db

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Record represents a database record
type Record struct {
	ID   int
	URL  string
	Tags []string
}

// Database holds Records
type Database struct {
	Records          []Record
	SearchCacheCount int
	ZeroResultsCount int

	searchResultCache map[string][]Record
	zeroResultsCache  []string
}

// Add a Record to a Database
func (database *Database) Add(record Record) {
	database.Records = append(database.Records, record)
}

// Search a Database
func (database *Database) Search(searchTerm string) []Record {
	var result []Record

	var putInZeroCache = true
	var putInResultCache = true

	var timeStart = time.Now()

	// first, search the zero cache
	for _, cachedZeroSearch := range database.zeroResultsCache {
		if cachedZeroSearch == searchTerm {
			var timeEnd = time.Now()
			var msg = "cached search for '%s' completed in %v with 0 results"

			putInZeroCache = false

			msg = fmt.Sprintf(msg, searchTerm, timeEnd.Sub(timeStart))

			log.Printf(msg)

			return nil
		}
	}

	// then, search the normal cache
	if _, ok := database.searchResultCache[searchTerm]; ok {
		var timeEnd = time.Now()
		var msg = "search for '%s' completed in %v with %d result"

		putInResultCache = false
		result := database.searchResultCache[searchTerm]

		msg = fmt.Sprintf(msg, searchTerm, timeEnd.Sub(timeStart), len(result))

		log.Printf(msg)

		return result
	}

	// lastly, search the database
	for _, record := range database.Records {
		for _, tag := range record.Tags {
			if tag == searchTerm {
				result = append(result, record)
				break
			}
		}
	}

	var timeEnd = time.Now()

	var msg = "search for '%s' completed in %v with %d result"

	if len(result) == 0 && putInZeroCache {
		database.zeroResultsCache = append(database.zeroResultsCache, searchTerm)
		database.ZeroResultsCount++
	}

	if putInResultCache {
		// TODO
	}

	if len(result) != 1 {
		msg = msg + "s"
	}

	msg = fmt.Sprintf(msg, searchTerm, timeEnd.Sub(timeStart), len(result))

	log.Printf(msg)

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
	log.Printf("DB server running on port %d...\n", port)
	http.HandleFunc("/", index)
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		api(w, r, database)
	})
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
