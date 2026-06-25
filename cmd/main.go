package main

import (
	"dash/internal/history"
	"dash/internal/tree"
)

func main() {
	tree := tree.NewTree(128, 64)
	err := history.Import(tree)
	if err != nil {
		panic(err)
	}
	tree.PrintTreeColor()
}
