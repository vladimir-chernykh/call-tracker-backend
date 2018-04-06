package audiosvc_test

import (
	"testing"
	"os/exec"
	"bytes"
	"github.com/vladimir-chernykh/call-tracker-backend/audiosvc"
	"github.com/stretchr/testify/assert"
	"strings"
	"github.com/vladimir-chernykh/call-tracker-backend/calltracker"
	"database/sql"
	"github.com/vladimir-chernykh/call-tracker-backend/postgres"
)

func TestExec(t *testing.T) {
	cmd := exec.Command("ls", "-la")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func TestConvert(t *testing.T) {
	converter := audiosvc.AudioService{}

	wavFile, err := converter.Convert("fixture.aac")
	if err != nil {
		panic(err)
	}

	assert.True(t, len(*wavFile) >= 0)
	assert.True(t, strings.Contains(*wavFile, ".wav"))

	cmd := exec.Command("rm", *wavFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	rErr := cmd.Run()
	if rErr != nil {
		panic(rErr)
	}

}

func TestSend(t *testing.T) {
	converter := audiosvc.AudioService{}
	res, err := converter.Send("f.wav")
	if err != nil {
		panic(err)
	}

	assert.True(t, len(*res) >= 0)
}

func TestGetStt(t *testing.T) {
	var DB *sql.DB
	DB, err := sql.Open("postgres", "user=backend dbname=call_tracker sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	s := postgres.New(DB)
	p := calltracker.Phone{Number: "+70000000000"}
	a := calltracker.Audio{Buffer: []byte{}}
	c := calltracker.Call{Audio: a, Phone: p, RemoteId: "6f16efee6ef54945b6ee44c3dbc4b44b"}

	_, sErr := s.Save(&c)
	if sErr != nil {
		panic(sErr)
	}

	meter := audiosvc.AudioService{Storage: s}
	meter.GetSTT(c)
}

func TestGetDuration(t *testing.T) {

}
