package main

import (
	"os"
	"fmt"
	"strconv"
	"strings"
	"math"
	"io/ioutil"
	"time"
)

func reviewCard(card string, grade, currentTime int) string {

	if grade < 0 || grade > 5 {
		fmt.Fprintf(os.Stderr, "Error: invalid grade '%v'.\n", grade)
		os.Exit(1)
	}

	fields := strings.Split(card, "	")

	if len(fields) != 6 {
		fmt.Fprintf(os.Stderr, "Error: card '%v' is invalid.\n", card)
		os.Exit(1)
	}

	// n: repetition number
	// EF: easiness factor
	// I: inter-repetition interval

	n, err := strconv.Atoi(fields[2])
	handle(err, fmt.Sprintf("Error: card '%v' is invalid.\n", card))
	EF, err := strconv.ParseFloat(fields[3], 64)
	handle(err, fmt.Sprintf("Error: card '%v' is invalid.\n", card))
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

	EF = EF + (0.1 - (5-float64(grade))*(0.08+(5-float64(grade))*0.02))
	if EF < 1.3 {
		EF = 1.3
	}

	// Convert back to strings and return
	fields[2] = strconv.Itoa(n)
	fields[3] = strconv.FormatFloat(EF, 'f', -1, 64)
	fields[4] = strconv.Itoa(I)
	fields[5] = strconv.Itoa(currentTime)

	return strings.Join(fields, "	")
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
			retryFilename := retryFilename(filename)

			// indexStr: No point converting back to a
			// string again when writing to the file later

			isDue := isCardDue(line, currentTime)
			isRetry := checkRetry(retryFilename, indexStr)

			if isDue || isRetry {
				lines[i] = reviewCard(line, grade, currentTime)
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
