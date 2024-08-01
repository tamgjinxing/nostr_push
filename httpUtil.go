package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func DoPost(url string, message *MessageBean) (string, error) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String(), nil
}
