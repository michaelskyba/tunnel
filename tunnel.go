package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func new_card(card string) string {

	// A new card only has two fields: front and back
	if len(strings.Split(card, "	")) != 2 {
		return card
	}

	return card + "	0	2.5	0	1617249600"
}

func check_due(card string, current_time int) bool {

	// To be due, (last review date) + (inter-repetition interval)
	// has to be before (current date)

	fields := strings.Split(card, "	")

	// Invalid cards can't be due
	if len(fields) != 6 {
		return false
	}

	interval, err := strconv.Atoi(fields[4])
	handle(err, fmt.Sprintf("Error: card '%v' is invalid.\n", card))
	last_review, err := strconv.Atoi(fields[5])
	handle(err, fmt.Sprintf("Error: card '%v' is invalid.\n", card))

	// Interval is in days, so we multiply by the number of seconds
	// in a day, which is 86400

	if last_review + interval*86400 < current_time {
		return true
	}
	return false
}

func review(card string, grade, current_time int) string {

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
	fields[5] = strconv.Itoa(current_time)

	return strings.Join(fields, "	")
}

func main() {
	len_of_args := len(os.Args)

	if len_of_args == 1 {
		user_error()
	}

	switch os.Args[1] + strconv.Itoa(len_of_args) {

	case "new_cards3":
		filename := os.Args[2]

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

	case "due3":
		file, err := os.Open(os.Args[2])
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

	case "front4", "back4":
		i, err := strconv.Atoi(os.Args[2])
		handle(err, "Error: non-integer card number provided.")

		card := get_line(os.Args[3], i)
		fields := strings.Split(card, "	")

		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "Error: line %v is not a valid card.\n", i)
			os.Exit(1)
		}

		if os.Args[1] == "front" {
			fmt.Println(fields[0])
		} else {
			fmt.Println(fields[1])
		}

	case "review5":
		filename := os.Args[4]
		file, err := ioutil.ReadFile(filename)
		handle(err, "Error: couldn't read deck file.")
		lines := strings.Split(string(file), "\n")

		deck_index, err := strconv.Atoi(os.Args[2])
		handle(err, "Error: non-integer card number provided.")
		grade, err := strconv.Atoi(os.Args[3])
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
				is_retry := check_retry(retry_filename, os.Args[2])

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
					update_retry(retry_filename, os.Args[2], grade)
				}
			}
		}

		write_lines(filename, lines)

	case "retry3":
		filename := get_retry(os.Args[2])
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

	default:
		user_error()
	}
}
