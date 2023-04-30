package main

import "os"

func main() {
	length := len(os.Args)
	if length == 1 {
		commandError()
	}

	commandName := os.Args[1]
	validateCommand(commandName, length)

	switch commandName {
	case "new_cards":
		// filename
		newCards(os.Args[2])

	case "due":
		// filename
		deckDue(os.Args[2])

	case "front", "back":
		// side, line, filename
		getSide(commandName, os.Args[2], os.Args[3])

	case "review":
		// index (string), grade (string), filename
		reviewDeck(os.Args[2], os.Args[3], os.Args[4])

	case "retry":
		// filename
		listRetry(os.Args[2])
	}
}
