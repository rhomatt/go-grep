package main

import (
	"bufio"
	"os"
	"io/fs"
	"fmt"
	"regexp"
)

// return a match if found
func processLine(line , pattern string) (string, bool) {
	re := regexp.MustCompilePOSIX(pattern)
	match := re.FindString(line)

	if match != "" {
		return match, true
	}

	return "", false
}

// we use fs.File, not os.File because fs.File is an interface that implements Read
func processFile(file fs.File) int {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	count := 0

	for inputLeft := scanner.Scan(); inputLeft ; inputLeft = scanner.Scan(){
		count++
		line := scanner.Text()
		fmt.Printf("line: %d:%s\n", count, line)

	}

	return count
}

/*
grep [flags] patterns [file]
goals: 
support -r, -i
support stdin and file
*/
func main() {


	// check if input is piped
	f, e := os.Stdin.Stat()
	if e != nil {
		panic(e)
	}

	if f.Mode() & fs.ModeCharDevice == 0 {
		// grab from stdin
		fmt.Println("From stdin")
		processFile(os.Stdin)
	} else {
		// grab from file
		fmt.Println("From file")
	}


}
