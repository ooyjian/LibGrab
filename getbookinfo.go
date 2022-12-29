package main

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func parseEntry(n *html.Node, m map[string]string) {
	if n == nil {
		printlnWrapper("The input node is nil", 1)
		return
	}

	if n.DataAtom != atom.Td {
		printlnWrapper("The element needs to be a <td>", 1)
		return
	}

	checkElemExist := func(elem string) bool {
		if n.FirstChild == nil {
			printlnWrapper("There are no "+elem+" provided", 1)
			return false
		}
		return true
	}

	for i := 1; i < 11; i++ {
		n = n.NextSibling.NextSibling
		switch i {
		case 1: // author, grab name
			author := n.FirstChild.FirstChild.Data
			m["author"] = author
			printlnWrapper(author, 1)
		case 2: // book title and link
			title := n.FirstChild.FirstChild.Data
			m["title"] = title
			printlnWrapper(title, 1)
		case 3:
			if checkElemExist("PUBLISHER") {
				publisher := n.FirstChild.Data
				m["publisher"] = publisher
				printlnWrapper(publisher, 1)
			}
		case 4:
			if checkElemExist("YEAR") {
				year := n.FirstChild.Data
				m["year"] = year
				printlnWrapper(year, 1)
			}
		case 5:
			if checkElemExist("PAGES") {
				pages := n.FirstChild.Data
				m["pages"] = pages
				printlnWrapper(pages, 1)
			}
		case 7:
			size := n.FirstChild.Data
			m["size"] = "size"
			printlnWrapper(size, 1)
		case 8:
			extension := n.FirstChild.Data
			m["extension"] = extension
			printlnWrapper(extension, 1)
		case 9:
			if checkElemExist("MIRROR1") {
				for _, attr := range n.FirstChild.Attr {
					if attr.Key == "href" {
						m["mirror1"] = attr.Val
						printlnWrapper(attr.Val, 1)
					}
				}
			}
		case 10:
			n = n.PrevSibling
			if checkElemExist("MIRROR2") {
				for _, attr := range n.FirstChild.Attr {
					if attr.Key == "href" {
						m["mirror2"] = attr.Val
						printlnWrapper(attr.Val, 1)
					}
				}
			}
		}
	}
}

func getBookInfo(n *html.Node) bool {
	for n.DataAtom != atom.Body {
		if n == nil {
			printlnWrapper("Can't find <body>", 1)
			return false
		}
		if n.DataAtom == atom.Html {
			n = n.FirstChild
		}
		n = n.NextSibling
	}
	n = n.FirstChild

	bookTable := make([]map[string]string, 30, 30)
	for n != nil {
		// printlnWrapper(n.Data)
		for _, attr := range n.Attr {
			// printlnWrapper("Key: " + attr.Key)
			// printlnWrapper("Val: " + attr.Val)
			if attr.Key == "class" && attr.Val == "c" {
				n = n.FirstChild.FirstChild.NextSibling
				if n.DataAtom == 0 {
					printlnWrapper("There are no books with this name available.", 1)
					return false
				}
				// parse to the actual book rows (starting from the second <tr>)
				var wg sync.WaitGroup
				bookIndex := 0
				for row := n; row != nil; row = row.NextSibling {
					rowElems := row.FirstChild
					if rowElems == nil {
						continue
					}
					bookMap := make(map[string]string)
					bookTable[bookIndex] = bookMap
					bookIndex++
					wg.Add(1)
					go func() {
						parseEntry(rowElems, bookMap)
						wg.Done()
					}()
				}
				wg.Wait()
				for j := 0; j < len(bookTable); j++ {
					m := bookTable[j]
					for k, v := range m {
						printlnWrapper(k+": "+v, 2)
					}
				}
				return true
			}
		}
		n = n.NextSibling
	}
	printlnWrapper("Can't find the right <table>", 1)
	return false
}

func MakeRequest(name string) error {
	name = strings.Replace(name, " ", "+", -1)
	url := "http://libgen.rs/search.php?req=" + name + "&lg_topic=libgen&open=0&view=simple&res=25&phrase=1&column=def"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		printlnWrapper(resp.Status, 1)
		return nil
	}

	body, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}
	body = body.FirstChild

	if !getBookInfo(body) {
		return nil
	}

	return nil
}
