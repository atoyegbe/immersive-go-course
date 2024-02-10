package cmd

import (
	"os"
	"fmt"
)

func Execute() {

	dir := os.Args[1]
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Print("Error while reading currecnt directory")
		os.Exit(1)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
