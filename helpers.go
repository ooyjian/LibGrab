package main

import (
	"net/http"
)

type BadCodeErr struct {
	respStatus string
}

func (e BadCodeErr) Error() string {
	return e.respStatus
}

func getRequest(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		printlnWrapper(resp.Status, 100)
		return resp, BadCodeErr{respStatus: resp.Status}
	}
	return resp, nil
}
