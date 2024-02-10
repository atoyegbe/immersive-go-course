package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <program> <file1> [<file2> ...]")
		return
	}
	
	commands := os.Args[1:]
	for _, command := range commands {
		fileContent, err := os.ReadFile(command)
		if err != nil {
			fmt.Print(err)
			continue
		}

		if _, err := os.Stdout.Write(fileContent); err != nil {
			fmt.Printf("Error writing content of file %s to stdout: %v\n", command, err)
		}
		fmt.Println("")
	}
}
