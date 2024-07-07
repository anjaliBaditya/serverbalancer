package main

import (
	"fmt"
	"os"
)

// handleError logs the error and exits the program
func handleError(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
