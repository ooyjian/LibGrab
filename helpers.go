package main

import (
	"net/http"
)

func getRequest(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		printlnWrapper(resp.Status, 100)
		return nil
	}
}
