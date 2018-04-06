package audiosvc

import (
	"github.com/vladimir-chernykh/call-tracker-backend/calltracker"
	"io/ioutil"
	"strings"
	"os/exec"
	"bytes"
)

type AudioService struct {
	storage calltracker.CallStorage
}

func New(storage calltracker.CallStorage) calltracker.AudioService {
	return &AudioService{storage}
}

func (c *AudioService) Process(aac *calltracker.Call) (*calltracker.Call, error) {

	aacFile, err := c.storage.Dump(aac)
	if err != nil {
		panic(err)
	}

	wavFile, err := c.Convert(*aacFile)
	if err != nil {
		return nil, err
	}

	wav, err := ioutil.ReadFile(*wavFile)
	if err != nil {
		panic(err)
	}

	remoteId, sErr := c.Send(wav)
	if sErr != nil {
		return nil, sErr
	}
	aac.RemoteId = remoteId

	go c.GetDuration()
	go c.GetSTT()

	return nil, nil
}

func (c *AudioService) Convert(aacFile string) (*string, error) {
	wavFile := strings.Replace(aacFile, ".aac", ".wav", 1)

	cmd := exec.Command("ffmpeg", "-i", aacFile, wavFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return &wavFile, nil
}

func (c *AudioService) Send(aac []byte) (string, error) {
	return "", nil
}

func (c *AudioService) GetDuration() (string, error) {
	return "", nil
}

func (c *AudioService) GetSTT() (string, error) {
	return "", nil
}
