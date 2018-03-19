package calltrackersvc

import (
	"net/http"
	"bytes"
	"strings"
	"fmt"
	"io"

	"github.com/gorilla/mux"
	"github.com/vladimir-chernykh/call-tracker-backend/calltracker"
	"database/sql"
	"github.com/vladimir-chernykh/call-tracker-backend/postgres"
)

func ReceiveFileHandler(DB *sql.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rr *http.Request) {
		vars := mux.Vars(rr)

		fmt.Println("AAAAAAAA", vars["phone"])

		var Buf bytes.Buffer
		// in your case file would be fileupload
		file, header, err := rr.FormFile("audio")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		name := strings.Split(header.Filename, ".")
		fmt.Printf("File name %s\n", name[0])
		// Copy the file data to my buffer
		io.Copy(&Buf, file)

		// do something with the contents...
		// I normally have a struct defined and unmarshal into a struct, but this will
		// work as an example
		//contents := Buf.String()
		//fmt.Println(contents)
		// I reset the buffer in case I want to use it again
		// reduces memory allocations in more intense projects
		p := calltracker.Phone{Number: vars["phone"]}
		a := calltracker.Audio{Buffer: Buf.Bytes()}
		c := calltracker.Call{Phone: p, Audio: a}

		s := postgres.CallService{DB}
		id, err := s.Save(&c)
		if err != nil {
			panic(err)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		fmt.Println("ID: ", id)
		Buf.Reset()
		// do something else
		// etc write header
		rw.WriteHeader(http.StatusCreated)
		return
	})
}
