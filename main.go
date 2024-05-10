package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pressly/goose"
	"log"
	"net/http"
)

func main() {

	dbMigrate()

	router := mux.NewRouter()
	router.HandleFunc("/flag", createFeatureFlag).Methods("POST")
	router.HandleFunc("/flag/{slug}", updateFeatureFlag).Methods("POST")
	router.HandleFunc("/flags", getFeatureFlags).Methods("GET")
	router.HandleFunc("/", homepage).Methods("GET")

	log.Println("Listening for requests")
	log.Println("http://localhost:8001")
	log.Fatal(http.ListenAndServe(":8001", router))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, World!")
	respondWithJSON(w, 200, "Little Flagger")
}

type FeatureFlag struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Percentage  int    `json:"percentage"`
}

func createFeatureFlag(w http.ResponseWriter, r *http.Request) {
	var ff FeatureFlag

	err := json.NewDecoder(r.Body).Decode(&ff)

	if err != nil {
		panic("Invalid Params")
	}

	ff.Id = uuid.New().String()

	db, _ := dbOpen()
	defer db.Close()

	fmt.Println(ff)

	insert, _ := db.Query("REPLACE INTO feature_flag VALUES (?, ?, ?, ?, ?, ?)", ff.Id, ff.Name, ff.Slug, ff.Description, ff.Status, ff.Percentage)

	defer func(insert *sql.Rows) {
		err := insert.Close()
		if err != nil {

		}
	}(insert)

	respondWithJSON(w, 201, ff)
}

func getFeatureFlags(w http.ResponseWriter, r *http.Request) {
	db, _ := dbOpen()
	defer db.Close()

	results, _ := db.Query("SELECT * FROM feature_flag")

	var flags []FeatureFlag

	for results.Next() {
		var ff FeatureFlag
		err := results.Scan(&ff.Id, &ff.Name, &ff.Slug, &ff.Description, &ff.Status, &ff.Percentage)

		if err != nil {
			panic(err.Error())
		}

		flags = append(flags, ff)
	}

	respondWithJSON(w, 200, flags)
}

func updateFeatureFlag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	slug := vars["slug"]

	var ff FeatureFlag

	err := json.NewDecoder(r.Body).Decode(&ff)

	if err != nil {
		panic("Invalid Params")
	}

	ff.Id = uuid.New().String()

	respondWithJSON(w, 200, slug)
}

func dbMigrate() {
	db, err := dbOpen()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = goose.SetDialect("mysql")
	if err != nil {
		panic(err)
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		panic(err)
	}
}

func dbOpen() (*sql.DB, error) {
	fmt.Println("Connect to Database")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/lf")

	return db, err
}
