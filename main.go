package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"sync"
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

func processFile(file *os.File, pattern string, wg *sync.WaitGroup) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for inputLeft := scanner.Scan(); inputLeft ; inputLeft = scanner.Scan(){
		line := scanner.Text()
		match, matchFound := processLine(line, pattern)
		if matchFound {
			fmt.Printf("%s: %s\n", file.Name(), match)
		}
	}
}

func processPath(file *os.File, pattern string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer file.Close()
	stat, e := file.Stat()
	if e != nil {
		panic(e)
	}
	if !stat.Mode().IsDir() {
		processFile(file, pattern, wg)
		return
	}

	fileNames, e := file.Readdirnames(0)
	if e != nil {
		panic(e)
	}
	for _, fileName := range(fileNames) {
		newFile, e := os.Open(file.Name() + "/" + fileName)
		if e != nil {
			panic(e)
		}
		wg.Add(1)
		go processPath(newFile, pattern, wg)
	}
}

/*
grep [flags] patterns [file]
goals: 
support -r, -i
support stdin and file
*/
func main() {
	if len(os.Args) < 2 {
		fmt.Println("no pattern provided")
		os.Exit(1)
	}
	pattern := os.Args[1]
	var target *os.File

	// check if input is piped
	f, e := os.Stdin.Stat()
	if e != nil {
		panic(e)
	}

	if f.Mode() & fs.ModeCharDevice == 0 {
		// grab from stdin
		fmt.Println("From stdin")
		target = os.Stdin
	} else {
		// grab from file
		fmt.Println("From file")
		if len(os.Args) < 3 {
			fmt.Println("no file provided")
			os.Exit(1)
		}

		target, e = os.Open(os.Args[2])
		if e != nil {
			panic(e)
		}
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go processPath(target, pattern, wg)
	wg.Wait()
}
