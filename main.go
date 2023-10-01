package main

import (
	"flag"
	"fmt"
	"os"
)

const DEBUG_LEVEL int = 10

var max_books int
var ext string
var path string

func printlnWrapper(content string, level int) {
	if level > DEBUG_LEVEL {
		fmt.Println(content)
	}
}

func init() {
	flag.IntVar(&max_books, "n", 5, "The maximum number of results that is shown.")
	flag.StringVar(&ext, "e", "", "Limit results to a certain file extension.")
	flag.StringVar(&path, "p", "./", "The path to store the downloaded document.")
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("You need to provide the name of a book.")
		return
	}
	book := args[1]
	// To make sure flag.Parse() works as intended
	os.Args = args[1:]

	flag.Parse()

	err := makeRequest(book)
	if err != nil {
		fmt.Println("Download failed :(")
		fmt.Println(err)
	}
}
