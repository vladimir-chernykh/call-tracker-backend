package postgres_test

import (
	"database/sql"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/vladimir-chernykh/call-tracker-backend/postgres"
	"github.com/vladimir-chernykh/call-tracker-backend/calltracker"
)

func TestSave(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(0, 0)
	var DB *sql.DB
	DB, err := sql.Open("postgres", "user=backend dbname=call_tracker sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	s := postgres.CallService{DB}
	p := calltracker.Phone{Number: "+70000000000"}
	a := calltracker.Audio{Buffer: []byte{}}
	c := calltracker.Call{Audio: a, Phone: p}
	id, err := s.Save(&c)
	if err != nil {
		panic(err)
	}
	assert.True(*id >= 0)
}
