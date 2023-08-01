package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

var (
	URL string = "https://api.assemblyai.com/v2/upload"
	key string = "58ee1dccb3784af98e4d6047488b1624"
)

func handleErr(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, e.Error())
}

func Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handleEvent(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "route not found")
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	var fileBuff, fileErr = io.ReadAll(r.Body)

	if fileErr != nil {
		handleErr(w, fileErr)
	}

	var transcriptURL, URLErr = getURL(fileBuff)

	if URLErr != nil {
		handleErr(w, URLErr)
		return
	}

	var transcriptID, IDErr = getTranscriptID(transcriptURL)

	if IDErr != nil {
		handleErr(w, IDErr)
		return
	}

	var text, err = getText(transcriptID)

	if err != nil {
		handleErr(w, IDErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, text)
}

func getTranscriptID(transcriptURL map[string]string) (string, error) {
	var jsml, _ = json.Marshal(transcriptURL)
	var req, reqErr = http.NewRequest("POST", "https://api.assemblyai.com/v2/transcript", bytes.NewBuffer(jsml))

	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", key)

	if reqErr != nil {
		return "", reqErr
	}

	defer req.Body.Close()

	var client = http.Client{}
	var result, err = client.Do(req)

	if err != nil {
		return "", err
	}

	var idMap = map[string]string{}

	json.NewDecoder(result.Body).Decode(&idMap)
	return idMap["id"], nil
}

func getURL(fb []byte) (map[string]string, error) {
	var HTTPReq, HTTPErr = http.NewRequest("POST", URL, bytes.NewReader(fb))
	HTTPReq.Header.Set("authorization", key)

	if HTTPErr != nil {
		return nil, HTTPErr
	}

	defer HTTPReq.Body.Close()

	var client = &http.Client{}
	var result, err = client.Do(HTTPReq)

	if err != nil {
		return nil, err
	}

	var data = map[string]string{}

	json.NewDecoder(result.Body).Decode(&data)

	data["audio_url"] = data["upload_url"]
	delete(data, "upload_url")

	return data, nil
}

func getText(transcriptID string) (string, error) {
	time.Sleep(time.Second * 10)
	var req, reqErr = http.NewRequest("GET", "https://api.assemblyai.com/v2/transcript/"+transcriptID, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", key)

	if reqErr != nil {
		return "", reqErr
	}

	var client *http.Client = &http.Client{}
	var res, resErr = client.Do(req)
	if resErr != nil {
		return "", resErr
	}

	var data = map[string]string{}
	json.NewDecoder(res.Body).Decode(&data)

	if data["text"] != "" {
		return data["text"], nil
	}

	time.Sleep(time.Second * 10)
	return getText(transcriptID)
}
