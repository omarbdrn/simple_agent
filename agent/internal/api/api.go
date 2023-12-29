package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func PerformRequest(request HTTPRequest) (*http.Response, error) {
	body := []byte{}

	if request.IsJson {
		if requestBodyBytes, ok := request.Body.([]byte); ok {
			body = requestBodyBytes
		} else {
			body, _ = json.Marshal(request.Body)
		}
	} else {
		if requestBodyBytes, ok := request.Body.([]byte); ok {
			body = requestBodyBytes
		} else {
			body = []byte(request.Body.(string))
		}
	}

	req, err := http.NewRequest(request.Method, BaseURL+request.Endpoint, bytes.NewBuffer(body))

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("User-Agent", "SimpleAgent v0.1")

	if request.IsJson {
		req.Header.Set("Content-Type", "application/json")
	}

	if len(request.Headers) > 0 {
		for key, value := range request.Headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	if res.StatusCode != http.StatusOK {
		return &http.Response{}, errors.New("404 Not Found")
	}

	return res, nil
}
