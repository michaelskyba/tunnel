package main

import (
	"strings"
	"strconv"
	"fmt"
	"os"
	"bufio"
	"time"
)

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
