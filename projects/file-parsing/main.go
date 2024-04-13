package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// {"name": "Prisha", "high_score": 30},
type Player struct {
	Name      string `json:"name"`
	HighScore int    `json:"high_score"`
}

func parseJsonFile() {
	file, err := os.Open("examples/json.txt")
	if err != nil {
		fmt.Printf("Unable to open file %s", err)
		return
	}

	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	var players []Player

	for _, line := range strings.Split(string(fileContent), "\n") {
		if line == "[" || line == "]" {
			continue
		}
		if strings.TrimSpace(line) == "" {
            continue
        }
		line = strings.Replace(line, "},", "}", -1)
		line = strings.TrimSpace(line)

		var player Player

		if err := json.Unmarshal([]byte(line), &player); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		players = append(players, player)
	}

	fmt.Printf("len of players: %d ", len(players))

}

func main() {
	parseJsonFile()
}
