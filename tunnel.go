package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

func get_line(filename string, target int) string {

	// O(n) time; don't use get_line in a loop

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

func write_lines(filename string, lines []string) {
	err := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
	handle(err, fmt.Sprintf("Error: couldn't write to %v.", filename))
}

// fail_list manipulates the deck retry file to account for a review.
// If this function is running, the assumption is that this card was
// *validly* reviewed just now. So, either the card was due, or the card
// was on the first cycle of the fail list.
func update_retry(filename, line_number string, grade int) {

	// The discarding of the error here is deliberate. If the file doesn't
	// exist, it's fine, because we'll create it in write_lines() later. If the
	// permissions are wrong, we'll get an error when we write to the file, which
	// means that the user will still know about it. It would be a better UX to catch
	// it early but I don't know how to do that idiomatically... I can't find any
	// error variables in ioutil's online documentation.
	file, _ := ioutil.ReadFile(filename)

	// If the file doesn't exist, lines == []string{""}, which works fine.
	lines := strings.Split(string(file), "\n")

	second_cycle := false
	card_found := false
	var card_index int

	// If the line is already there, we need to add it to the next retry cycle
	for i, line := range lines {

		if line == "-" {
			if !second_cycle {
				second_cycle = true
				continue
			}

			// There is more than one "-"
			fmt.Fprintln(os.Stderr, "Error: broken retry file")
			os.Exit(1)
		}

		if line == line_number {

			// It's impossible to have more than two retry cycles by default.
			// So, if we're seeing the line number in the retry cycle, it has
			// to be in the first cycle, not the second cycle.
			if second_cycle {
				fmt.Fprintln(os.Stderr, "Error: broken retry file")
				os.Exit(1)
			}

			// We can't stop the loop because we still need to know if
			// the retry file has one cycle or two cycles
			if !card_found {
				card_found = true
				card_index = i

			// Multiple of the same card in the retry file
			} else {
				fmt.Fprintln(os.Stderr, "Error: broken retry file")
				os.Exit(1)
			}
		}
	}

	// We don't want a random newline in the middle, we want it at the end
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if card_found {

		// We're moving the line to the second cycle, so we
		// want to remove it from the first cycle
		lines = append(lines[:card_index], lines[card_index+1:]...)

		if grade < 4 {
			if second_cycle {
				lines = append(lines, line_number)
			} else {
				lines = append(lines, "-", line_number)
			}
		}

		// If we leave the cycle indicator after the first cycle is complete,
		// it will be flagged as a broken file next time. Also, we don't need
		// the cycle indicator if the reviewed card was the sole card in the
		// retry file.
		if lines[0] == "-" {
			lines = lines[1:]
		}

	} else {
		lines = append(lines, line_number)
	}

	// Keep a newline at the end
	// - It's still considered a text file this way
	// - It makes it easy to edit in Kakoune without problems
	lines = append(lines, "")

	write_lines(filename, lines)
}

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

func check_retry(filename, line_number string) bool {
	// Makeshift
	return true
}

func review(card string, grade, current_time int) string {
	fields := strings.Split(card, "	")

	if len(fields) != 6 {
		fmt.Fprintf(os.Stderr, "Error: card '%v' is invalid.\n", card)
		os.Exit(1)
	}

	// n: repetition number
	// EF: easiness factor
	// I: inter-repition interval

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

		line_number, err := strconv.Atoi(os.Args[2])
		handle(err, "Error: non-integer card number provided.")
		grade, err := strconv.Atoi(os.Args[3])
		handle(err, "Error: non-integer grade number provided.")

		current_time := int(time.Now().Unix())

		// Not worth using get_line() because we need to update "lines"
		for i, line := range lines {
			if i == line_number {
				absolute, err := filepath.Abs(filename)
				handle(err, "Error: broken deck path?")

				// Make sure retry file's parent directory exists

				tmp_dir := os.Getenv("TMPDIR")
				if tmp_dir == "" {
					tmp_dir = "/tmp"
				}

				path := strings.Split(absolute, "/")
				output_path := fmt.Sprintf("%v/tunnel%v", tmp_dir, strings.Join(path[:len(path)-1], "/"))

				err = os.MkdirAll(output_path, 0755)
				handle(err, fmt.Sprintf("Error: couldn't create %v.", output_path))

				retry_filename := fmt.Sprintf("%v/tunnel%v", tmp_dir, absolute)

				// os.Args[2]: No point converting back to a
				// string again when writing to the file later

				is_due := check_due(line, current_time)
				is_retry := check_retry(retry_filename, os.Args[2])

				if is_due || is_retry {
					lines[i] = review(line, grade, current_time)
				} else {
					fmt.Fprintf(os.Stderr, "Error: card %v is not due for review.\n", line_number)
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

	default:
		user_error()
	}
}
