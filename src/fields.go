package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func newCard(card string, currentTime int) string {
	// A new card only has two fields: front and back
	if len(strings.Split(card, "	")) != 2 {
		return card
	}

	return fmt.Sprintf("%v	0	2.5	0	%v", card, currentTime)
}

// tunnel new_cards
func newCards(filename string) {
	file, err := ioutil.ReadFile(filename)
	handle(err, "Error: couldn't read deck file.")
	lines := strings.Split(string(file), "\n")

	// We only want to update the file if it's changed, just in case someone
	// has a problem with the last modified date being updated
	changed := false

	// If it's one second before the current (creation) date, it guarantees
	// that it will be due. This seems more elegant than having a random
	// past date like 1617249600
	lastReview := int(time.Now().Unix()) - 1

	for i, line := range lines {
		lines[i] = newCard(line, lastReview)

		if lines[i] != line {
			changed = true
		}
	}

	if changed {
		writeLines(filename, lines)
	}
}

// tunnel front and tunnel back
func getSide(side, line, filename string) {
	i, err := strconv.Atoi(line)
	handle(err, "Error: non-integer card number provided.")

	card := getLine(filename, i)
	fields := strings.Split(card, "	")

	if len(fields) < 2 {
		printError(fmt.Sprintf("Error: line %v is not a valid card.\n", i))
	}

	if side == "front" {
		fmt.Println(fields[0])
	} else {
		fmt.Println(fields[1])
	}
}
