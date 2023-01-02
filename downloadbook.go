package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/html"
)

func getDownloadLink(mirrorLink string) (string, error) {
	resp, err := getRequest(mirrorLink)
	if err != nil {
		return "", err
	}

	body, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	body = body.FirstChild

	findHtmlBody(&body)

	downloadDiv := body.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.NextSibling

	// Find the actual download link
	linkNode := downloadDiv.FirstChild.NextSibling.FirstChild
	for _, attr := range linkNode.Attr {
		if attr.Key == "href" {
			downloadLink := attr.Val
			printlnWrapper(downloadLink, 5)
			return downloadLink, nil
		}
	}

	return "", nil
}

func requestDownload(link, filepath, title string) error {
	resp, err := getRequest(link)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath + title)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Download Successful! Your document is downloaded at " + path + title)

	return nil
}
