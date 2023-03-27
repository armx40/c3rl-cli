package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func network_request(url string, params map[string]string, headers map[string]string, data interface{}) (body []byte, err error) {

	var request_method string

	if data != nil {
		request_method = "POST" // POST
	} else {
		request_method = "GET" // GET
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	postData, err := json.Marshal(data)
	if err != nil {
		return
	}

	req, err := http.NewRequest(request_method, url, bytes.NewReader(postData))
	if err != nil {
		err = fmt.Errorf("Got error %s", err.Error())
		return
	}

	/* set headers */
	for i := range headers {
		req.Header.Set(i, headers[i])
	}
	/* */

	/* prepare url */
	q := req.URL.Query()
	for i := range params {
		q.Add(i, params[i])
	}
	req.URL.RawQuery = q.Encode()
	/**/

	response, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("Got error %s", err.Error())
		return
	}
	defer response.Body.Close()

	// Read the response body on the line below.
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return

}
