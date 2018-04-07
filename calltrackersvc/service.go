package calltrackersvc

import (
	"net/http"
	"bytes"
	"io"

	"github.com/gorilla/mux"
	"github.com/vladimir-chernykh/call-tracker-backend/calltracker"
	"database/sql"
	"github.com/vladimir-chernykh/call-tracker-backend/postgres"
	"github.com/vladimir-chernykh/call-tracker-backend/audiosvc"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func ReceiveFileHandler(DB *sql.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rr *http.Request) {
		vars := mux.Vars(rr)

		var Buf bytes.Buffer
		file, _, err := rr.FormFile("audio")
		if err != nil {
			panic(err)
		}
		defer file.Close()

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

		asvc := audiosvc.New(s)
		go asvc.Process(&c)

		Buf.Reset()
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		json.NewEncoder(rw).Encode(map[string]int64{"id": *id})
		return
	})
}

func GetCallResultsHandler(DB *sql.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rr *http.Request) {
		vars := mux.Vars(rr)

		s := postgres.New(DB)
		out, err := s.GetMetrics(vars["call"])
		if err != nil {
			log.Error(err)
			panic(err)
		}
		if out == nil {
			out = []byte("{}")
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(out)

		return
	})
}
