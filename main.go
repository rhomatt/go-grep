package main

import (
	"os"
	"fmt"
	_ "regexp"
)

/*
grep [flags] patterns [file]
goals: 
support multiple patterns
support -r, -i
support stdin and file
*/
func main() {
	f, e := os.Stdin.Stat()
	if e != nil {
		fmt.Println("error getting stdin info:", e)
		os.Exit(1)
	}
	if f.Size() > 0 {
		// grab from stdin
		fmt.Println("From stdin")
	} else {
		// grab from file
		fmt.Println("From file")
	}
}
