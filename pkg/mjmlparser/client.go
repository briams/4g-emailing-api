package mjmlparser

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// HTTPClient connects to apis
func HTTPClient(method string, URL string, reqBody []byte, headers map[string]string) (int, http.Header, []byte, error) {
	req, err := http.NewRequest(method, URL, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	req.SetBasicAuth("2f95d3a0-c3b0-4eb9-beec-6049441d5cb1", "077a988e-03a0-41cb-ab53-2592681366d4")
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, nil, nil, err
	}

	return res.StatusCode, res.Header, body, nil
}
