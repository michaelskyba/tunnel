package main

import (
	"os"
	"io/ioutil"
	"errors"
	"strings"
	"bufio"
	"time"
	"fmt"
	"strconv"
)

// tunnel new_cards
func new_cards(filename string) {
	file, err := ioutil.ReadFile(filename)
	handle(err, "Error: couldn't read deck file.")
	lines := strings.Split(string(file), "\n")

	// We only want to update the file if it's changed, just in case someone
	// has a problem with the last modified date being updated
	changed := false

	for i, line := range lines {
		lines[i] = new_card(line)

		if lines[i] != line {
			changed = true
		}
	}

	if changed {
		write_lines(filename, lines)
	}
}

// tunnel due
func deck_due(filename string) {
	file, err := os.Open(filename)
	defer file.Close()
	handle(err, "Error: couldn't read deck file.")

	i := 0
	scanner := bufio.NewScanner(file)

	current_time := int(time.Now().Unix())
	for scanner.Scan() {
		if check_due(scanner.Text(), current_time) {
			fmt.Println(i)
		}

		i++
	}
}

// tunnel front and tunnel back
func get_side(side, line, filename string) {
	i, err := strconv.Atoi(line)
	handle(err, "Error: non-integer card number provided.")

	card := get_line(filename, i)
	fields := strings.Split(card, "	")

	if len(fields) < 2 {
		fmt.Fprintf(os.Stderr, "Error: line %v is not a valid card.\n", i)
		os.Exit(1)
	}

	if side == "front" {
		fmt.Println(fields[0])
	} else {
		fmt.Println(fields[1])
	}
}

// tunnel review
func review_deck(index_str, grade_str, filename string) {
	file, err := ioutil.ReadFile(filename)
	handle(err, "Error: couldn't read deck file.")
	lines := strings.Split(string(file), "\n")

	deck_index, err := strconv.Atoi(index_str)
	handle(err, "Error: non-integer card number provided.")
	grade, err := strconv.Atoi(grade_str)
	handle(err, "Error: non-integer grade number provided.")

	current_time := int(time.Now().Unix())

	// The deck file is expected to end with a newline, so
	// e.g. len(lines) will be five if there are four cards.
	// Accessing this fourth card would use the "3" index
	// so we check if >= len(lines) - 1.
	if deck_index < 0 || deck_index >= len(lines) - 1{
		fmt.Fprintf(os.Stderr, "Error: no line %v in deck.\n", deck_index)
		os.Exit(1)
	}

	// Not worth using get_line() because we need to update "lines"
	for i, line := range lines {
		if i == deck_index {
			retry_filename := get_retry(filename)

			// os.Args[2]: No point converting back to a
			// string again when writing to the file later

			is_due := check_due(line, current_time)
			is_retry := check_retry(retry_filename, index_str)

			if is_due || is_retry {
				lines[i] = review(line, grade, current_time)
			} else {
				fmt.Fprintf(os.Stderr, "Error: card %v is not due for review.\n", deck_index)
				os.Exit(1)
			}

			if is_retry || (is_due && grade < 4) {

				// There's no DRY benefit to having this as a function
				// but I feel like it makes the contents of this switch case
				// quite a bit more organized. Feel free to bully me if this
				// is wrong style-wise.
				update_retry(retry_filename, index_str, grade)
			}
		}
	}

	write_lines(filename, lines)
}

// tunnel retry
func list_retry(deck_filename string) {
	filename := get_retry(deck_filename)
	file, err := os.Open(filename)
	defer file.Close()

	// If the file doesn't exist, there are no retries to do
	if errors.Is(err, os.ErrNotExist) {
		os.Exit(0)
	}
	handle(err, "Error: couldn't read retry file.")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// The user only needs to worry about the first cycle
		if line == "-" {
			break
		}
		fmt.Println(line)
	}
}
