package audiosvc

import (
	"github.com/vladimir-chernykh/call-tracker-backend/calltracker"
	"strings"
	"os/exec"
	"bytes"
	"net/http"
	"mime/multipart"
	"os"
	"io"
	"encoding/json"
)

type AudioService struct {
	storage calltracker.CallStorage
}

type sendResponse struct {
	Result result `json:"result"`
}

type result struct {
	ContentId string `json:"content_id"`
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

	remoteId, sErr := c.Send(*wavFile)
	if sErr != nil {
		return nil, sErr
	}
	aac.RemoteId = *remoteId

	cmd := exec.Command("rm", *aacFile, *wavFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	rErr := cmd.Run()
	if rErr != nil {
		panic(rErr)
	}

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

func (c *AudioService) Send(wav string) (*string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("audio", wav)
	if err != nil {
		panic(err)
	}

	fh, err := os.Open(wav)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		panic(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	res, err := http.Post("http://64.58.125.108:3000/content", contentType, bodyBuf)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	data := sendResponse{}
	json.NewDecoder(res.Body).Decode(&data)

	return &data.Result.ContentId, nil
}

func (c *AudioService) GetDuration() (string, error) {
	return "", nil
}

func (c *AudioService) GetSTT() (string, error) {
	return "", nil
}
