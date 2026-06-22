package main

import (
	"dash/internal/tree"
	"fmt"
)

func main() {
	tree := tree.NewTree(4, 2)

	// Insert 100 random keys
	for i := 0; i < 1000; i++ {

		key := fmt.Sprintf("cmd-%d", i)
		tree.Insert(key)

	}

	tree.PrintTreeColor()
}
