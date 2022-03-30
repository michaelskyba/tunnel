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
func newCards(filename string) {
	file, err := ioutil.ReadFile(filename)
	handle(err, "Error: couldn't read deck file.")
	lines := strings.Split(string(file), "\n")

	// We only want to update the file if it's changed, just in case someone
	// has a problem with the last modified date being updated
	changed := false

	for i, line := range lines {
		lines[i] = newCard(line)

		if lines[i] != line {
			changed = true
		}
	}

	if changed {
		writeLines(filename, lines)
	}
}

// tunnel due
func deckDue(filename string) {
	file, err := os.Open(filename)
	defer file.Close()
	handle(err, "Error: couldn't read deck file.")

	i := 0
	scanner := bufio.NewScanner(file)

	currentTime := int(time.Now().Unix())
	for scanner.Scan() {
		if checkDue(scanner.Text(), currentTime) {
			fmt.Println(i)
		}

		i++
	}
}

// tunnel front and tunnel back
func getSide(side, line, filename string) {
	i, err := strconv.Atoi(line)
	handle(err, "Error: non-integer card number provided.")

	card := getLine(filename, i)
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
func reviewDeck(indexStr, gradeStr, filename string) {
	file, err := ioutil.ReadFile(filename)
	handle(err, "Error: couldn't read deck file.")
	lines := strings.Split(string(file), "\n")

	deckIndex, err := strconv.Atoi(indexStr)
	handle(err, "Error: non-integer card number provided.")
	grade, err := strconv.Atoi(gradeStr)
	handle(err, "Error: non-integer grade number provided.")

	currentTime := int(time.Now().Unix())

	// The deck file is expected to end with a newline, so
	// e.g. len(lines) will be five if there are four cards.
	// Accessing this fourth card would use the "3" index
	// so we check if >= len(lines) - 1.
	if deckIndex < 0 || deckIndex >= len(lines) - 1{
		fmt.Fprintf(os.Stderr, "Error: no line %v in deck.\n", deckIndex)
		os.Exit(1)
	}

	// Not worth using getLine() because we need to update "lines"
	for i, line := range lines {
		if i == deckIndex {
			retryFilename := getRetry(filename)

			// indexStr: No point converting back to a
			// string again when writing to the file later

			isDue := checkDue(line, currentTime)
			isRetry := checkRetry(retryFilename, indexStr)

			if isDue || isRetry {
				lines[i] = review(line, grade, currentTime)
			} else {
				fmt.Fprintf(os.Stderr, "Error: card %v is not due for review.\n", deckIndex)
				os.Exit(1)
			}

			if isRetry || (isDue && grade < 4) {

				// There's no DRY benefit to having this as a function
				// but I feel like it makes the contents of this switch case
				// quite a bit more organized. Feel free to bully me if this
				// is wrong style-wise.
				updateRetry(retryFilename, indexStr, grade)
			}
		}
	}

	writeLines(filename, lines)
}

// tunnel retry
func listRetry(deckFilename string) {
	filename := getRetry(deckFilename)
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
