package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
)

var usage = `Valid commands:
	tunnel new_cards <deck filename>
	tunnel due <deck filename>
	tunnel <front|back> <line number> <deck filename>
	tunnel review <line number> <score> <deck filename>
See README.md for more information.`

func handle(err error, message string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, message)
		os.Exit(1)
	}
}

func write_lines(filename string, lines []string) {
	err := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
	handle(err, "Error: couldn't write to habit file")
}

func new_card(card string) string {

	// Only add fields to valid new cards
	if len(strings.Split(card, " ")) != 2 {
		return card
	}

	return card+"	0	2.5	0	2021-04-01"
}

func main() {

	len_of_args := len(os.Args)

	if len_of_args == 1 {
		fmt.Println(usage)
		os.Exit(1)
	}

	switch os.Args[1]+strconv.Itoa(len_of_args) {

	case "new_cards3":
		filename := os.Args[2]

		// Open deck file
		file, err := ioutil.ReadFile(filename)
		handle(err, "Error: couldn't read deck file")
		lines := strings.Split(string(file), "\n")

		// We only want to update the file if it's changed, just in case someone
		// has a problem with the last modified date being updated)
		changed := false

		// Set every line to its new_card() value
		for i, line := range lines {
			lines[i] = new_card(line)

			if lines[i] != line {
				changed = true
			}
		}

		if changed {
			write_lines(filename, lines)
		}

	case "due3":
		fmt.Println("due")

	case "front4":
		fmt.Println("front")

	case "back4":
		fmt.Println("back")

	case "review5":
		fmt.Println("review")

	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}
