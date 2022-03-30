package main

import (
	"os"
	"strconv"
)

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
