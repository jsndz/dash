package main

import (
	"dash/internal/history"
	"dash/internal/search"
	"dash/internal/tree"
	"fmt"
	"os"
)

func main() {

	args := os.Args
	command := args[1]

	tree := tree.NewTree(128, 64)
	engine := search.NewSearchEngine(tree)
	err := history.Import(tree)
	if err != nil {
		panic(err)
	}
	result := engine.Autocomplete(command)
	for _, r := range result {
		fmt.Println(r)
	}

}
