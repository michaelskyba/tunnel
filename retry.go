package main

import (
	"fmt"
	"path/filepath"
	"os"
	"io/ioutil"
	"strings"
	"bufio"
	"errors"
)

// retryFilename takes a relative filename of a deck file and returns the absolute
// filename of the related retry file
func retryFilename(filename string) string {
	absolute, err := filepath.Abs(filename)
	handle(err, "Error: broken deck path.")

	tmpDir := os.Getenv("TMPDIR")
	if tmpDir == "" {
		tmpDir = "/tmp"
	}

	path := strings.Split(absolute, "/")
	outputPath := fmt.Sprintf("%v/tunnel%v", tmpDir, strings.Join(path[:len(path)-1], "/"))

	// Make sure retry file's parent directory exists
	err = os.MkdirAll(outputPath, 0755)
	handle(err, fmt.Sprintf("Error: couldn't create %v.", outputPath))

	return fmt.Sprintf("%v/tunnel%v", tmpDir, absolute)
}

// updateRetry manipulates the deck retry file to account for a review.
// If this function is running, the assumption is that this card was
// *validly* reviewed just now. So, either the card was due, or the card
// was on the first cycle of the fail list.
func updateRetry(filename, deckIndex string, grade int) {

	// The discarding of the error here is deliberate. If the file doesn't
	// exist, it's fine, because we'll create it in writeLines() later. If the
	// permissions are wrong, we'll get an error when we write to the file, which
	// means that the user will still know about it. It would be a better UX to catch
	// it early but I don't know how to do that idiomatically... I can't find any
	// error variables in ioutil's online documentation.
	file, _ := ioutil.ReadFile(filename)

	// If the file doesn't exist, lines == []string{""}, which works fine.
	lines := strings.Split(string(file), "\n")

	secondCycle := false
	cardFound := false
	var retryIndex int

	// If the line is already there, we need to add it to the next retry cycle
	for i, line := range lines {

		if line == "-" {
			if !secondCycle {
				secondCycle = true
				continue
			}

			// There is more than one "-"
			fmt.Fprintln(os.Stderr, "Error: broken retry file.")
			os.Exit(1)
		}

		if line == deckIndex {

			// It's impossible to have more than two retry cycles by default.
			// So, if we're seeing the line number in the retry cycle, it has
			// to be in the first cycle, not the second cycle.
			if secondCycle {
				fmt.Fprintln(os.Stderr, "Error: broken retry file.")
				os.Exit(1)
			}

			// We can't stop the loop because we still need to know if
			// the retry file has one cycle or two cycles
			if !cardFound {
				cardFound = true
				retryIndex = i

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

	if cardFound {

		// We're moving the line to the second cycle, so we
		// want to remove it from the first cycle
		lines = append(lines[:retryIndex], lines[retryIndex+1:]...)

		if grade < 4 {
			if secondCycle {
				lines = append(lines, deckIndex)
			} else {
				lines = append(lines, "-", deckIndex)
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
		lines = append(lines, deckIndex)
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

	writeLines(filename, lines)
}

// I don't like the ambiguous names "checkRetry" and "listRetry" but I can't
// think of anything more descriptive.

// checkRetry checks if a given index exists in a given retry file
func checkRetry(filename, deckIndex string) bool {
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
		} else if line == deckIndex {
			return true
		}
	}

	return false
}

// tunnel retry command
func listRetry(deckFilename string) {
	filename := retryFilename(deckFilename)
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
