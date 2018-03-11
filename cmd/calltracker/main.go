package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

func main() {
	run()
}

func run() {
	r := mux.NewRouter()

	r.Handle("/api/v1/phones/{phone}", stubUploadAudioHandler())
	// Get port to listen from env
	srv := &http.Server{
		Addr:    ":80",
		Handler: r,
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func stubUploadAudioHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rr *http.Request) {
		vars := mux.Vars(rr)
		if vars["phone"] == "+79161298967" {
			rw.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			rw.WriteHeader(http.StatusCreated)
		}

	})
}
