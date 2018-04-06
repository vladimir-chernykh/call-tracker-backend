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
	"encoding/json"
)

func ReceiveFileHandler(DB *sql.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rr *http.Request) {
		vars := mux.Vars(rr)

		var Buf bytes.Buffer
		file, header, err := rr.FormFile("audio")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		name := strings.Split(header.Filename, ".")
		fmt.Printf("File name %s\n", name[0])

		io.Copy(&Buf, file)

		p := calltracker.Phone{Number: vars["phone"]}
		a := calltracker.Audio{Buffer: Buf.Bytes()}
		c := calltracker.Call{Phone: p, Audio: a}

		s := postgres.New(DB)
		id, err := s.Save(&c)
		if err != nil {
			panic(err)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		fmt.Println("ID: ", id)
		Buf.Reset()
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		json.NewEncoder(rw).Encode(map[string]int64{"id": *id})
		return
	})
}

func GetCallResultsStub(DB *sql.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rr *http.Request) {
		vars := mux.Vars(rr)

		if vars["call"] == "10" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("{\"stt\":{\"text\":\"Пусть просто такую информацию легче записать голосом чем писать текстом\"},\"duration\":{\"duration\":4.56}}"))
			return
		}

		if vars["call"] == "11" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("{\"stt\":{\"text\":\"Пусть просто такую информацию легче записать голосом чем писать текстом\"}}"))
			return
		}

		if vars["call"] == "12" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("{\"duration\":{\"duration\":4.56}}"))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte("{}"))
		rw.WriteHeader(http.StatusNotFound)

		return
	})
}
