package main

import (
	"dash/internal/history"
	"dash/internal/search"
	"dash/internal/tree"
	"fmt"
	"os"
	"strings"
)

func main() {

	args := os.Args
	command := args[1]
	input := args[2]
	tree := tree.NewTree(128, 64)
	engine := search.NewSearchEngine(tree)
	err := history.Import(tree)
	if err != nil {
		panic(err)
	}

	switch command {
	case "autocomplete":
		result := engine.Autocomplete(input)
		if input[len(input)-1] == ' ' {
			output, _ := strings.CutPrefix(result[0], input)
			fmt.Print(output)
		} else {
			parts := strings.Split(input, " ")
			check := parts[len(parts)-1]
			index := strings.Index(result[0], check)
			output := result[0][index:]
			fmt.Print(output)
		}

	}

}
