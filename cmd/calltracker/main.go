package main

import (
	"net/http"
	"fmt"
	"database/sql"
	"github.com/gorilla/mux"

	"github.com/vladimir-chernykh/call-tracker-backend/calltrackersvc"
)

func main() {
	run()
}

func run() {
	var DB *sql.DB
	DB, err := sql.Open("postgres", "host=172.17.0.1 user=backend dbname=call_tracker sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	r := mux.NewRouter()

	r.Handle("/api/v1/phones/{phone}", calltrackersvc.ReceiveFileHandler(DB))
	r.Handle("/api/v1/calls/{call}", calltrackersvc.GetCallResultsStub(DB))
	// Get port to listen from env
	srv := &http.Server{
		Addr:    ":80",
		Handler: r,
	}

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
