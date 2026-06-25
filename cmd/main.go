package main

import (
	"dash/internal/history"
	"dash/internal/search"
	"dash/internal/tree"
	"fmt"
)

func main() {
	tree := tree.NewTree(128, 64)
	engine := search.NewSearchEngine(tree)
	err := history.Import(tree)
	if err != nil {
		panic(err)
	}
	tree.PrintTreeColor()
	result := engine.Autocomplete("git")
	for _, cmd := range result {
		fmt.Println(cmd)
	}
}
