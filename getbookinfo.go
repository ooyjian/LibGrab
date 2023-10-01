package main

import (
	"errors"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type bookCounter struct {
	mu      sync.Mutex
	counter int
}

func (n *bookCounter) inc() {
	n.mu.Lock()
	n.counter++
	n.mu.Unlock()
}

func (n *bookCounter) count() int {
	n.mu.Lock()
	c := n.counter
	n.mu.Unlock()
	return c
}

func makeRequest(name string) error {
	name = strings.Replace(name, " ", "+", -1)
	url := "http://libgen.rs/search.php?req=" + name + "&lg_topic=libgen&open=0&view=simple&res=25&phrase=1&column=def"
	// alturl := "http://libgen.rs/search.php?&req=" + name + "&phrase=1&view=simple&column=def&sort=def&sortmode=ASC&page=1"
	resp, err := getRequest(url)
	if err != nil {
		return err
	}

	body, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}
	body = body.FirstChild

	bookTable, err := getBookInfo(body)
	if err != nil {
		return err
	}
	if bookTable != nil {
		err = displayBooks(bookTable)
		if err != nil {
			return err
		}
	}

	return nil
}

func getBookInfo(n *html.Node) ([]map[string]string, error) {
	if err := findHtmlBody(&n); err != nil {
		return nil, err
	}
	n = n.FirstChild

	bookTable := make([]map[string]string, 30, 30)
	for n != nil {
		printlnWrapper(n.Data, 1)
		for _, attr := range n.Attr {
			printlnWrapper("Key: "+attr.Key, 1)
			printlnWrapper("Val: "+attr.Val, 1)
			if attr.Key == "class" && attr.Val == "c" {
				n = n.FirstChild.FirstChild.NextSibling
				if n == nil {
					printlnWrapper("There are no books with this name available.", 100)
					return nil, errors.New("There are no books with this name available.")
				}
				if n.DataAtom == 0 {
					printlnWrapper("There are no books with this name available.", 100)
					return nil, errors.New("There are no books with this name available.")
				}
				// parse to the actual book rows (starting from the second <tr>)
				var wg sync.WaitGroup
				ct := bookCounter{counter: 0}
				bookIndex := 0
				for row := n; row != nil; row = row.NextSibling {
					if ct.count() == max_books {
						break
					}
					rowElems := row.FirstChild
					if rowElems == nil {
						continue
					}
					bookMap := make(map[string]string)
					bookTable[bookIndex] = bookMap
					bookIndex++
					wg.Add(1)
					go func() {
						if err := parseEntry(rowElems, bookMap, &ct); err != nil {
							printlnWrapper(err.Error(), 100)
						}
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

				return bookTable, nil
			}
		}
		n = n.NextSibling
	}
	printlnWrapper("Can't find the right <table>", 1)
	return nil, errors.New("Can't find the right <table>")
}

func findHtmlBody(n **html.Node) error {
	for (*n).DataAtom != atom.Body {
		if *n == nil {
			printlnWrapper("Can't find <body>", 100)
			return errors.New("Can't find <body>")
		}
		if (*n).DataAtom == atom.Html {
			*n = (*n).FirstChild
		}
		*n = (*n).NextSibling
	}
	return nil
}

func parseEntry(n *html.Node, m map[string]string, b *bookCounter) error {
	if n == nil {
		printlnWrapper("The input node is nil", 100)
		return errors.New("The input node is nil")
	}

	if n.DataAtom != atom.Td {
		printlnWrapper("The element needs to be a <td>", 10)
		return errors.New("The element needs to be a <td>")
	}

	checkElemExist := func(elem string) bool {
		if n.FirstChild == nil {
			printlnWrapper("There are no "+elem+" provided", 3)
			return false
		}
		return true
	}

	for i := 1; i < 10; i++ {
		n = n.NextSibling.NextSibling
		switch i {
		case 1: // author, grab name
			var author string
			if n.FirstChild.FirstChild != nil {
				author = n.FirstChild.FirstChild.Data
			}
			m["author"] = author
			printlnWrapper(author, 3)
		case 2: // book title and link
			var title string
			if checkElemExist("TITLE") {
				for temp := n.FirstChild; temp != nil; temp = temp.NextSibling {
					for _, attr := range temp.Attr {
						if attr.Key == "id" {
							title = temp.FirstChild.Data
							break
						}
					}
				}
			}
			m["title"] = title
			printlnWrapper(title, 3)
		case 3:
			if checkElemExist("PUBLISHER") {
				publisher := n.FirstChild.Data
				m["publisher"] = publisher
				printlnWrapper(publisher, 3)
			}
		case 4:
			if checkElemExist("YEAR") {
				year := n.FirstChild.Data
				m["year"] = year
				printlnWrapper(year, 3)
			}
		case 5:
			if checkElemExist("PAGES") {
				pages := n.FirstChild.Data
				m["pages"] = pages
				printlnWrapper(pages, 3)
			}
		case 7:
			if checkElemExist("SIZE") {
				size := n.FirstChild.Data
				m["size"] = size
				printlnWrapper(size, 3)
			}
		case 8:
			var extension string
			if checkElemExist("EXTENSION") {
				extension = n.FirstChild.Data
			}
			if len(ext) > 0 && extension != ext {
				m["title"] = ""
				return nil
			}
			m["extension"] = extension
			printlnWrapper(extension, 3)
		case 9:
			if checkElemExist("MIRROR1") {
				for _, attr := range n.FirstChild.Attr {
					if attr.Key == "href" {
						m["mirror1"] = attr.Val
						printlnWrapper(attr.Val, 3)
						b.inc()
						break
					}
				}
			} else {
				n = n.NextSibling
				if checkElemExist("MIRROR2") {
					for _, attr := range n.FirstChild.Attr {
						if attr.Key == "href" {
							m["mirror2"] = attr.Val
							printlnWrapper(attr.Val, 3)
							b.inc()
							break
						}
					}
				} else {
					m["title"] = "" // so it does not get displayed
				}
			}
		}
	}
	return nil
}
