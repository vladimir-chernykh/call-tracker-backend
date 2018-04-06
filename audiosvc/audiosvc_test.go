package audiosvc_test

import (
	"testing"
	"os/exec"
	"bytes"
	"github.com/vladimir-chernykh/call-tracker-backend/audiosvc"
	"github.com/stretchr/testify/assert"
	"strings"
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
