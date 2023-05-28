package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

func readFileLines(file *os.File, lines []string) (int, []string) {
	linesDict := make(map[string]int)
	input := bufio.NewScanner(file)
	var cont int

	for input.Scan() {
		text := input.Text()
		if text == "" {
			continue
		}

		linesDict[text]++

		if linesDict[text] == 1 {
			lines = append(lines, text)
			cont++
		}
	}
	return cont, lines
}

func readFile(fileName string, lines []string) (int, []string) {
	if fileName == "" {
		return readFileLines(os.Stdin, lines)
	} else {
		file, err := os.Open(fileName)

		if err != nil {
			log.Fatal("Error reading file; ", err)
		}

		return readFileLines(file, lines)
	}

}

func checkIfMatch(regexs []string, domain string, output chan string) {
	for key := range regexs {
		match, _ := regexp.MatchString(regexs[key], domain)
		if match {
			output <- domain
			break
		}
	}
}

func main() {
	var dFlag = flag.String("d", "", "File with a list of domains.")
	var rFlag = flag.String("r", "", "File with a list of regex.")
	var tFlag = flag.Int("t", 10, "Number of threads.")
	flag.Parse()

	threads := make(chan struct{}, *tFlag)
	var domains, regexs []string

	_, domains = readFile(*dFlag, domains)

	output := make(chan string, len(domains))

	if *rFlag == "" {
		log.Fatal("Please provide a valid regex file")
	}

	_, regexs = readFile(*rFlag, regexs)

	for _, domain := range domains {
		go func(domain string) {

			threads <- struct{}{}
			defer func() { <-threads }()
			checkIfMatch(regexs, domain, output)

		}(domain)
	}

	for range domains {
		fmt.Println(<-output)

	}
}
