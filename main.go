package main

import "fmt"

// Global variables go here
var DEBUG_LEVEL int = 4

func printlnWrapper(content string, level int) {
	if level > DEBUG_LEVEL {
		fmt.Println(content)
	}
}

func main() {
	makeRequest("hello world")
}
