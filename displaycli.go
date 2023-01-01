package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/table"
)

func displayBooks(bt []map[string]string) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedColumnLengths([]int{3, 30, 20, 10, 10, 9})

	cleanbt := make([]map[string]string, 0, len(bt))

	t.AppendHeader(table.Row{"#", "Title", "Author", "Size", "Ext", "Publisher", "Year", "Pages"})

	counter := 0
	for i, r := range bt {
		row := make([]interface{}, 8, 8)
		row[0] = i + 1
		row[1] = r["title"]
		row[2] = r["author"]
		row[3] = r["size"]
		row[4] = r["extension"]
		row[5] = r["publisher"]
		row[6] = r["year"]
		row[7] = r["pages"]
		if row[1] == "" {
			continue
		}
		counter++
		if counter > max_books {
			break
		}
		cleanbt = append(cleanbt, r)
		t.AppendRow(row)
	}

	style := table.StyleColoredDark
	style.Options = table.Options{true, true, true, true, true}
	t.SetStyle(style)
	t.Render()

	var id int
	fmt.Print("Enter the book you want to download (1-" + strconv.Itoa(max_books) + "): ")
	_, err := fmt.Scanf("%d", &id)
	if err != nil {
		return err
	}
	id -= 1 // to match the right index of the book
	if id < 0 || id >= max_books {
		printlnWrapper("The input is not within range.", 100)
		return nil
	}

	for i := 1; i < 3; i++ {
		link, err := getDownloadLink(cleanbt[id]["mirror"+strconv.Itoa(i)])
		if len(link) == 0 {
			printlnWrapper("Could not find a link", 100)
			continue
		}
		err = requestDownload(link, "./", cleanbt[id]["title"]+"."+cleanbt[id]["extension"])
		if err == nil {
			return nil
		}
	}

	return err
}
