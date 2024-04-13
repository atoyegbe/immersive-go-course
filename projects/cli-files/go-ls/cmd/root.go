package cmd

import (
	"fmt"
	"os"
	"strings"
)

func Execute() {

	commands := os.Args[1:]
	var addComma = false
	for _, command := range commands {
		if strings.HasPrefix(command, "-") {
			if command == "-m" {
				addComma = true
			} else {
				fmt.Printf("Unknown option: %s\n", command)
			}
			continue
		}
		files, err := os.ReadDir(command)
		if err != nil {
			fmt.Print("Error while reading current directory")
			os.Exit(1)
		}
		if addComma {
			for _, file := range files {
				fmt.Print(file, ",")
			}
		} else {
			for _, file := range files {
				fmt.Println(file)
			}
		}
	}
}
