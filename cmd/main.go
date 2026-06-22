package main

import (
	"dash/internal/tree"
	"fmt"
	"math/rand"
)

func main() {
	tree := tree.NewTree(4, 2)

	inserted := make(map[string]bool)

	// Insert 100 random keys
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("cmd-%d", rand.Intn(1000))
		tree.Insert(key)
		inserted[key] = true
	}

	tree.PrintTreeColor()
}
