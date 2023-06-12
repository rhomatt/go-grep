package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"sync"
)

type matchInfo struct {
	lineNumber int
	line string
	matching string
}

type safeResults struct {
	mu *sync.Mutex
	matches map[string]matchInfo
}

func (sr *safeResults) addResult(lineNumber int, file, line, matching string) {
	sr.mu.Lock()
	sr.matches[file] = matchInfo{lineNumber, line, matching}
	sr.mu.Unlock()
}

// return a match if found
func processLine(line , pattern string) (string, bool) {
	re := regexp.MustCompilePOSIX(pattern)
	match := re.FindString(line)

	if match != "" {
		return match, true
	}

	return "", false
}

func processFile(file *os.File, pattern string, wg *sync.WaitGroup, results *safeResults) {
	defer wg.Done()
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	lineNumber := 1

	for inputLeft := scanner.Scan(); inputLeft ; inputLeft = scanner.Scan() {
		line := scanner.Text()
		match, matchFound := processLine(line, pattern)
		if matchFound {
			results.addResult(lineNumber, file.Name(), line, match)
			//fmt.Printf("%s: %s\n", file.Name(), match)
		}
		lineNumber++
	}
}

func processPath(file *os.File, pattern string, wg *sync.WaitGroup, results *safeResults) {
	stat, e := file.Stat()
	if e != nil {
		panic(e)
	}
	if !stat.Mode().IsDir() {
		wg.Add(1)
		go processFile(file, pattern, wg, results)
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
		processPath(newFile, pattern, wg, results)
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

	results := &safeResults{mu: new(sync.Mutex), matches: make(map[string]matchInfo)}
	wg := new(sync.WaitGroup)
	processPath(target, pattern, wg, results)
	wg.Wait()

	for file, result := range(results.matches) {
		fmt.Printf("%s: %d: %s", file, result.lineNumber, result.matching)
	}
}
