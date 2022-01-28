package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
	"time"
	"math"
	"path/filepath"
)


func user_error() {
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

// O(n) time; don't use get_line in a loop
func get_line(filename string, target int) string {
	if target >= 0 {

		// Open
		file, err := os.Open(filename)
		defer file.Close()
		handle(err, "Error: couldn't read deck file.")

		i := 0
		scanner := bufio.NewScanner(file)

		// Iterate while incrementing until we hit (target)
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

func write_lines(filename string, lines []string) {
	err := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
	handle(err, "Error: couldn't write to habit file.")
}

// Mark a card for a retry in the deck retry file
func fail_list(absolute string, line_number int) {
	fmt.Println("You failed the card lol")
	fmt.Printf("Path: %v\n", absolute)
	fmt.Printf("Line number: %v\n", line_number)
}

func new_card(card string) string {

	// Only add fields to valid new cards
	if len(strings.Split(card, "	")) != 2 {
		return card
	}

	return card+"	0	2.5	0	1617249600"
}

func due(card string, current_time int) bool {
	fields := strings.Split(card, "	")

	// Not a valid card, so we already know it can't be due
	if len(fields) != 6 {
		return false
	}

	// To be due, (last review date) + (inter-repetition interval)
	// has to be before (current date)

	interval, err1 := strconv.Atoi(fields[4])
	last_review, err2 := strconv.Atoi(fields[5])

	handle(err1, fmt.Sprintf("Error: card '%v' is invalid.\n", card))
	handle(err2, fmt.Sprintf("Error: card '%v' is invalid.\n", card))

	// Interval is in days, so we multiply by the number of seconds
	// in a day, which is 86400

	if last_review + interval * 86400 < current_time {
		return true
	}
	return false
}

func review(card string, grade int, current_time int) string {
	fields := strings.Split(card, "	")

	// Invalid card detection
	if len(fields) != 6 {
		fmt.Fprintf(os.Stderr, "Error: card '%v' is invalid.\n", card)
		os.Exit(1)
	}

	// Repetition number
	n, err := strconv.Atoi(fields[2])
	handle(err, fmt.Sprintf("Error: card '%v' is invalid.\n", card))

	// Easiness factor
	EF, err := strconv.ParseFloat(fields[3], 64)
	handle(err, fmt.Sprintf("Error: card '%v' is invalid.\n", card))

	// Inter-repetition interval
	I, err := strconv.Atoi(fields[4])
	handle(err, fmt.Sprintf("Error: card '%v' is invalid.\n", card))

	// SM-2
	if grade >= 3 {

		if n == 0 {
			I = 1
		} else if n == 1 {
			I = 6
		} else {
			I = int(math.Round(float64(I) * EF))
		}

		n++

	} else {
		n = 0
		I = 1
	}

	EF = EF + (0.1 - (5 - float64(grade)) * (0.08 + (5 - float64(grade)) * 0.02))
    if EF < 1.3 {
	    EF = 1.3
    }

	// Convert to strings and return
    fields[2] = strconv.Itoa(n)
    fields[3] = strconv.FormatFloat(EF, 'f', -1, 64)
    fields[4] = strconv.Itoa(I)
    fields[5] = strconv.Itoa(current_time)

    return strings.Join(fields, "	")
}

func main() {

	len_of_args := len(os.Args)

	if len_of_args == 1 {
		user_error()
	}

	switch os.Args[1]+strconv.Itoa(len_of_args) {

	case "new_cards3":
		filename := os.Args[2]

		// Open deck file
		file, err := ioutil.ReadFile(filename)
		handle(err, "Error: couldn't read deck file.")
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
		// Open deck file
		file, err := os.Open(os.Args[2])
		defer file.Close()
		handle(err, "Error: couldn't read deck file.")

		i := 0
		scanner := bufio.NewScanner(file)

		// Iterate and print due for review
		current_time := int(time.Now().Unix())
		for scanner.Scan() {
			if due(scanner.Text(), current_time) {
				fmt.Println(i)
			}

			i++
		}

	case "front4", "back4":
		i, err := strconv.Atoi(os.Args[2])
		handle(err, "Error: non-integer card number provided.")

		card := get_line(os.Args[3], i)
		fields := strings.Split(card, "	")

		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "Error: line %v is not a valid card.\n", i)
			os.Exit(1)
		}

		// Print first or second field depending on the argument
		if os.Args[1] == "front" {
			fmt.Println(fields[0])
		} else {
			fmt.Println(fields[1])
		}

	case "review5":
		filename := os.Args[4]

		// Open deck file
		file, err := ioutil.ReadFile(filename)
		handle(err, "Error: couldn't read deck file.")

		lines := strings.Split(string(file), "\n")

		line_number, err := strconv.Atoi(os.Args[2])
		handle(err, "Error: non-integer card number provided.")

		grade, err := strconv.Atoi(os.Args[3])
		handle(err, "Error: non-integer grade number provided.")

		current_time := int(time.Now().Unix())

		// Not worth using get_line() because we need to update "lines"
		for i, line := range lines {
			if i == line_number {

				if !due(line, current_time) {
					fmt.Fprintf(os.Stderr, "Error: card %v is not due for review.\n", line_number)
					os.Exit(1)
				}

				lines[i] = review(line, grade, current_time)

				// Keep track of which cards need to be retried
				if grade < 4 {
					absolute, err := filepath.Abs(filename)
					handle(err, "Error: broken deck path?")

					fail_list(absolute, line_number)
				}
			}
		}

		write_lines(filename, lines)

	default:
		user_error()
	}
}
