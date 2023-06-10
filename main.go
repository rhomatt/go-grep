package main

import (
	"os"
	"io/fs"
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
		panic(e)
	}
	if f.Mode() & fs.ModeCharDevice == 0 {
		// grab from stdin
		fmt.Println("From stdin")
	} else {
		// grab from file
		fmt.Println("From file")
	}
}
