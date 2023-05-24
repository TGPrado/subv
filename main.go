package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func readFile(file *os.File, lines map[string]int) int {
	input := bufio.NewScanner(file)
	var cont int
	for input.Scan() {
		text := input.Text()
		if text == "" {
			continue
		}

		lines[text]++
		if lines[text] == 1 {
			cont++
		}
	}
	return cont
}

func getDomains(fileName string, lines map[string]int) int {
	if fileName == "" {
		return readFile(os.Stdin, lines)
	} else {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal("Error reading file; ", err)
		}
		return readFile(file, lines)
	}

}

func main() {
	var dFlag = flag.String("d", "", "File with a list of domains.")
	//var fFlag = flag.String("f", "", "File with a list of regex.")
	//var tFlag = flag.Int("t", 10, "Number of threads.")
	flag.Parse()
	lines := make(map[string]int)
	lineNumber := getDomains(*dFlag, lines)
	fmt.Println(lineNumber)
}
