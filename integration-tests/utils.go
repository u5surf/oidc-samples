package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	mathRand "math/rand"
	"net/http"
	"strconv"
	"time"
)

func newRequest(url string, method string, payload interface{}, headers ...header) (req *http.Request, err error) {
	if method == "GET" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		reqPayloadJSON, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(method, url, bytes.NewReader(reqPayloadJSON))
	}

	if err != nil {
		return nil, err
	}

	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	return req, err
}

func getResponse(req *http.Request, httpClient *http.Client) (responseBody []byte, cookies []*http.Cookie, err error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Whenever we're redirected we take the Location and return it.
	if resp.StatusCode == 302 || resp.StatusCode == 301 {
		redirectLocation, err := resp.Location()
		if err != nil {
			return nil, nil, err
		}
		return []byte(redirectLocation.String()), resp.Cookies(), nil
	}

	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return responseBody, resp.Cookies(), nil
}

func makeRequest(httpClient *http.Client, url string, method string, payload interface{}, headers ...header) (responseBody []byte, err error) {
	req, err := newRequest(url, method, payload, headers...)
	if err != nil {
		return nil, err
	}

	res, _, err := getResponse(req, httpClient)
	return res, err
}

func hex2bytes(s string) []byte {
	lgh := len(s)
	data := make([]byte, lgh/2)

	for i := 0; i < lgh; i += 2 {
		a, _ := strconv.ParseInt(s[i:i+2], 16, 32)
		data[i/2] = byte(a)
	}
	return data
}

func randPIN() int {
	mathRand.Seed(time.Now().UnixNano())
	return mathRand.Intn(9000) + 1000
}
