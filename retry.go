package main

import (
	"fmt"
	"path/filepath"
	"os"
	"io/ioutil"
	"strings"
)

// get_retry takes a relative filename of a deck file and returns the absolute
// filename of the related retry file
func get_retry(filename string) string {
	absolute, err := filepath.Abs(filename)
	handle(err, "Error: broken deck path.")

	tmp_dir := os.Getenv("TMPDIR")
	if tmp_dir == "" {
		tmp_dir = "/tmp"
	}

	path := strings.Split(absolute, "/")
	output_path := fmt.Sprintf("%v/tunnel%v", tmp_dir, strings.Join(path[:len(path)-1], "/"))

	// Make sure retry file's parent directory exists
	err = os.MkdirAll(output_path, 0755)
	handle(err, fmt.Sprintf("Error: couldn't create %v.", output_path))

	return fmt.Sprintf("%v/tunnel%v", tmp_dir, absolute)
}

// update_retry manipulates the deck retry file to account for a review.
// If this function is running, the assumption is that this card was
// *validly* reviewed just now. So, either the card was due, or the card
// was on the first cycle of the fail list.
func update_retry(filename, deck_index string, grade int) {

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
	var retry_index int

	// If the line is already there, we need to add it to the next retry cycle
	for i, line := range lines {

		if line == "-" {
			if !second_cycle {
				second_cycle = true
				continue
			}

			// There is more than one "-"
			fmt.Fprintln(os.Stderr, "Error: broken retry file.")
			os.Exit(1)
		}

		if line == deck_index {

			// It's impossible to have more than two retry cycles by default.
			// So, if we're seeing the line number in the retry cycle, it has
			// to be in the first cycle, not the second cycle.
			if second_cycle {
				fmt.Fprintln(os.Stderr, "Error: broken retry file.")
				os.Exit(1)
			}

			// We can't stop the loop because we still need to know if
			// the retry file has one cycle or two cycles
			if !card_found {
				card_found = true
				retry_index = i

			// Multiple of the same card in the retry file
			} else {
				fmt.Fprintln(os.Stderr, "Error: broken retry file.")
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
		lines = append(lines[:retry_index], lines[retry_index+1:]...)

		if grade < 4 {
			if second_cycle {
				lines = append(lines, deck_index)
			} else {
				lines = append(lines, "-", deck_index)
			}
		}

		// If we leave the cycle indicator after the first cycle is complete,
		// it will be flagged as a broken file next time. Also, we don't need
		// the cycle indicator if the reviewed card was the sole card in the
		// retry file.
		if len(lines) > 0 && lines[0] == "-" {
			lines = lines[1:]
		}

	} else {
		lines = append(lines, deck_index)
	}

	if len(lines) == 0 {
		err := os.Remove(filename)
		handle(err, "Error: couldn't remove retry file.")
		return
	}

	// Keep a newline at the end
	// - It's still considered a text file this way
	// - It makes it easy to edit in Kakoune without problems
	lines = append(lines, "")

	write_lines(filename, lines)
}

// check_retry checks if a given index exists in a given retry file
func check_retry(filename, deck_index string) bool {
	file, err := os.Open(filename)
	defer file.Close()

	// Let's just assume that the error means the file doesn't exist
	if err != nil {
		return false
	}

	var line string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()

		// All the cards from the first retry cycle have to be retried before
		// the cards from the second retry cycle
		if line == "-" {
			return false
		} else if line == deck_index {
			return true
		}
	}

	return false
}
