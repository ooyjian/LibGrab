package main

import "fmt"

// Global variables go here
var DEBUG_LEVEL int = 1

func printlnWrapper(content string, level int) {
	if level > DEBUG_LEVEL {
		fmt.Println(content)
	}
}

func main() {
	// set debug level
	DEBUG_LEVEL = 0
	MakeRequest("hello world")
}
