package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/pressly/goose"
	"log"
	"net/http"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", homepage).Methods("GET")
	router.HandleFunc("/flag", createFeatureFlag).Methods("POST")
	router.HandleFunc("/flags", getFeatureFlags).Methods("GET")

	log.Println("Listening for requests")
	log.Println("http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
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
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Status     string `json:"status"`
	Percentage int    `json:"percentage"`
}

func createFeatureFlag(w http.ResponseWriter, r *http.Request) {
	var ff FeatureFlag

	err := json.NewDecoder(r.Body).Decode(&ff)

	if err != nil {
		panic("Invalid Params")
	}
}

func getFeatureFlags(w http.ResponseWriter, r *http.Request) {

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
	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/lf")

	return db, err
}
