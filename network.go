package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func networkRequest(url string, data interface{}, responsePayload interface{}) error {

	var requestType bool

	if data != nil {
		requestType = false // POST
	} else {
		requestType = true // GET
	}

	var resp *http.Response
	var err error
	var postData []byte

	if requestType {
		resp, err = http.Get(url)
	} else {
		postData, err = json.Marshal(data)
		if err != nil {
			return err
		}
		resp, err = http.Post(url, "application/json", bytes.NewBuffer(postData))
	}

	if err != nil {
		return err
	}

	// Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responsePayload)

	if err != nil {
		return err
	}

	return nil

}
