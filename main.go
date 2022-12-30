package main

import (
	"fmt"
	"os"
	"strconv"
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
	book := args[1]
	var err error
	MAX_BOOKS, err = strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Wrong format of input.")
		return
	}
	makeRequest(book)
}
