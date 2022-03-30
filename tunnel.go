package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func newCard(card string) string {

	// A new card only has two fields: front and back
	if len(strings.Split(card, "	")) != 2 {
		return card
	}

	return card + "	0	2.5	0	1617249600"
}

func checkDue(card string, currentTime int) bool {

	// To be due, (last review date) + (inter-repetition interval)
	// has to be before (current date)

	fields := strings.Split(card, "	")

	// Invalid cards can't be due
	if len(fields) != 6 {
		return false
	}

	interval, err := strconv.Atoi(fields[4])
	handle(err, fmt.Sprintf("Error: card '%v' is invalid.\n", card))
	lastReview, err := strconv.Atoi(fields[5])
	handle(err, fmt.Sprintf("Error: card '%v' is invalid.\n", card))

	// Interval is in days, so we multiply by the number of seconds
	// in a day, which is 86400

	if lastReview + interval*86400 < currentTime {
		return true
	}
	return false
}

func review(card string, grade, currentTime int) string {

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

func main() {
	argsLength := len(os.Args)

	if argsLength == 1 {
		userError()
	}

	switch os.Args[1] + strconv.Itoa(argsLength) {

	case "new_cards3":
		// filename
		newCards(os.Args[2])

	case "due3":
		// filename
		deckDue(os.Args[2])

	case "front4", "back4":
		// side, line, filename
		getSide(os.Args[1], os.Args[2], os.Args[3])

	case "review5":
		// index (string), grade (string), filename
		reviewDeck(os.Args[2], os.Args[3], os.Args[4])

	case "retry3":
		// filename
		listRetry(os.Args[2])

	default:
		userError()
	}
}
