package main

import (
	"flag"
	"fmt"
	"os"
)

// Global variables go here
var DEBUG_LEVEL int = 50
var MAX_BOOKS int = 3

func printlnWrapper(content string, level int) {
	if level > DEBUG_LEVEL {
		fmt.Println(content)
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("You need to provide the name of a book.")
		return
	}
	book := args[1]

	numBooksPtr := flag.Int("num", 3, "The maximum number of results that is shown.")
	// extPtr := flag.String("ext", "", "Limit results to a certain file extension.")
	MAX_BOOKS = *numBooksPtr

	makeRequest(book)
}
