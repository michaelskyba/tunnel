package main

import (
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"strings"
)

func userError() {
	fmt.Println(`Valid commands:
	tunnel new_cards <deck filename>
	tunnel due <deck filename>
	tunnel <front|back> <line number> <deck filename>
	tunnel review <line number> <score> <deck filename>
	tunnel retry <deck filename>
See README.md for more information.`)
	os.Exit(1)
}

func handle(err error, message string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, message)
		os.Exit(1)
	}
}

func getLine(filename string, target int) string {

	// O(n) time; don't use getLine in a loop

	if target >= 0 {
		file, err := os.Open(filename)
		defer file.Close()
		handle(err, "Error: couldn't read deck file.")

		i := 0
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			if i == target {
				return scanner.Text()
			}

			i++
		}
	}

	fmt.Fprintf(os.Stderr, "Error: no line %v in deck.\n", target)
	os.Exit(1)

	return ""
}

func writeLines(filename string, lines []string) {
	err := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
	handle(err, fmt.Sprintf("Error: couldn't write to %v.", filename))
}
