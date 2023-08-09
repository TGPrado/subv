package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

func readFileLines(file *os.File, linesDict map[string]int) {
	input := bufio.NewScanner(file)

	for input.Scan() {
		text := input.Text()
		if text == "" {
			continue
		}

		linesDict[text]++

	}
}

func readFile(fileName string, linesDict map[string]int) {
	if fileName == "" {
		readFileLines(os.Stdin, linesDict)
	} else {
		file, err := os.Open(fileName)

		if err != nil {
			log.Fatal("Error reading file; ", err)
		}

		readFileLines(file, linesDict)
	}

}

func checkIfMatch(regexs map[string]int, domain string, output chan string) {
	for regex, _ := range regexs {
		match, _ := regexp.MatchString(regex, domain)
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
	domains := make(map[string]int)
	regexs := make(map[string]int)

	readFile(*dFlag, domains)

	output := make(chan string, len(domains))

	if *rFlag == "" {
		log.Fatal("Please provide a valid regex file")
	}

	readFile(*rFlag, regexs)

	for domain, _ := range domains {
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
