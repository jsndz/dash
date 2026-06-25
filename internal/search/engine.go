package search

import (
	"dash/internal/ranking"
	"dash/internal/tree"
)

type SearchEngine struct {
	tree *tree.Bptree
}

func NewSearchEngine(tree *tree.Bptree) *SearchEngine {
	return &SearchEngine{
		tree: tree,
	}
}

func (e *SearchEngine) Autocomplete(prefix string) []string {
	commands := e.tree.RangeScan(prefix)
	ranks := ranking.GetRanking(commands)
	return ranks
}
