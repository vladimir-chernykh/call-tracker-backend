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
	"net/url"
	"io/ioutil"
)

type AudioService struct {
	Storage calltracker.CallStorage
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

func (c *AudioService) Process(call *calltracker.Call) (error) {

	aacFile, err := c.Storage.Dump(call)
	if err != nil {
		panic(err)
	}

	wavFile, err := c.Convert(*aacFile)
	if err != nil {
		panic(err)
	}

	remoteId, err := c.Send(*wavFile)
	if err != nil {
		panic(err)
	}
	call.RemoteId = *remoteId

	cmd := exec.Command("rm", *aacFile, *wavFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	rErr := cmd.Run()
	if rErr != nil {
		panic(rErr)
	}

	go c.GetDuration(*call)
	go c.GetSTT(*call)

	return nil
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

func (c *AudioService) GetDuration(call calltracker.Call) (error) {
	err := c.getMetric("duration", call.RemoteId, call)
	if err != nil {
		panic(err)
	}
	return nil
}

func (c *AudioService) GetSTT(call calltracker.Call) (error) {
	err := c.getMetric("stt", call.RemoteId, call)
	if err != nil {
		panic(err)
	}
	return nil
}

func (c *AudioService) getMetric(name string, remoteId string, call calltracker.Call) (error) {
	res, err := http.PostForm("http://64.58.125.108:3000/"+name, url.Values{"content_id": {remoteId}})
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	m := calltracker.Metric{Name: name, Call: call, Data: body}
	c.Storage.SaveMetric(&m)

	return nil
}
