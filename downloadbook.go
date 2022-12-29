package main

import (
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func getDownloadLink(mirrorLink string) error {
	resp, err := http.Get(mirrorLink)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if resp.StatusCode != 200 {
		printlnWrapper(resp.Status, 100)
	}

	body, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}
	body = body.FirstChild

	findHtmlBody(&body)

	return nil
}
