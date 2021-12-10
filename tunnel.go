package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
)

var usage = `Valid commands:
	tunnel new_cards <deck filename>
	tunnel due <deck filename>
	tunnel <front|back> <line number> <deck filename>
	tunnel review <line number> <score> <deck filename>
See README.md for more information.`

func main() {

	len_of_args := len(os.Args)

	if len_of_args == 1 {
		fmt.Println(usage)
		os.Exit(1)
	}

	switch os.Args[1]+strconv.Itoa(len_of_args) {

	case "new_cards3":
		filename := os.Args[2]

		// Open deck file
		file, _ := ioutil.ReadFile(filename)
		lines := strings.Split(string(file), "\n")

		for _, line := range lines {
			fmt.Println(line)
		}

	case "due3":
		fmt.Println("due")

	case "front4":
		fmt.Println("front")

	case "back4":
		fmt.Println("back")

	case "review5":
		fmt.Println("review")

	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}
