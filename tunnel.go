package main

import (
	"fmt"
	"os"
)

var usage = `Valid commands:
	tunnel new_cards <deck filename>
	tunnel due <deck filename>
	tunnel <front|back> <line number> <deck filename>
	tunnel review <line number> <score> <deck filename>
See README.md for more information.`

func main() {

	if len(os.Args) == 1 {
		fmt.Println(usage)
		os.Exit(1)
	}

	switch os.Args[1] {

	case "new_cards":
		fmt.Println("new_cards")

	case "due":
		fmt.Println("due")

	case "front":
		fmt.Println("front")

	case "back":
		fmt.Println("back")

	case "review":
		fmt.Println("review")

	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}
