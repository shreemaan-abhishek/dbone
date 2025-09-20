package main

import (
	"fmt"
	"strings"
)

func main() {
	var input string
	for {
		fmt.Print("dbone > ")
		fmt.Scan(&input)
		if strings.HasPrefix(input, ".") {
			if doMetaCommand(input) != 0 {
				return
			}
			continue
		}

		prepareCommand(input)

	}
}

func prepareCommand(statement string) {
	command := strings.Fields(statement)[0]
	switch command {
	case "select":
		fmt.Println("selecting...")
	case "insert":
		fmt.Println("inserting...")
	default:
		fmt.Println("unrecognized command")
	}
}

func doMetaCommand(input string) int {
	switch input {
	case ".exit":
		return -1
	default:
		fmt.Println("unrecognized command")
		return 0
	}
}
