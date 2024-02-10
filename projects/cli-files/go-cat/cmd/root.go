package cmd

import (
	"fmt"
	"os"
	"strings"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <program> <file1> [<file2> ...]")
		return
	}

	commands := os.Args[1:]
	var appendLine bool = false
	for _, command := range commands {
		if strings.HasPrefix(command, "-") {
			if command == "-n" {
				appendLine = true
			} else {
				fmt.Printf("Unknown option: %s\n", command)
			}
			continue
		}
		fileContent, err := os.ReadFile(command)
		if err != nil {
			fmt.Print(err)
			continue
		}
		if appendLine {
			lines := strings.Split(string(fileContent), "\n")
			for i, line := range lines {
				fmt.Printf("%6d  %s\n", i+1, line)
			}
		} else {
			if _, err := os.Stdout.Write(fileContent); err != nil {
				fmt.Printf("Error writing content of file %s to stdout: %v\n", command, err)
			}
			fmt.Println("")
		}
	}
}
