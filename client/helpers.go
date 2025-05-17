package client

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func GetRequest[T any](url string, mode string, params map[string]string, headers map[string]string) (T, error) {
	var results T

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return results, err
	}

	reqParams := req.URL.Query()
	for k, v := range params {
		reqParams.Add(k, v)
	}
	req.URL.RawQuery = reqParams.Encode()

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}

	if mode == "xml" {
		err = xml.NewDecoder(resp.Body).Decode(&results)
	} else if mode == "json" {
		err = json.NewDecoder(resp.Body).Decode(&results)
	}
	if err != nil {
		return results, err
	}

	return results, err
}
