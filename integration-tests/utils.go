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

func Request(url string, method string, payload interface{}, headers ...header) (responseBody []byte, err error) {
	var req *http.Request
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

	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//whenever we're redirected we take the Location and return it
	if resp.StatusCode == 302 || resp.StatusCode == 301 {
		redirectLocation, err := resp.Location()
		if err != nil {
			return nil, err
		}
		return []byte(redirectLocation.String()), nil
	}

	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
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

// isNil returns true if the error is nil
func isNil(err error) bool {
	return err == nil
}
