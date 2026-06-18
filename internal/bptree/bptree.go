package bptree

import (
	"sort"
	"time"
)

const MAX = 127
const MIN = 63

type Bptree struct {
	Root *Node
}

type Command struct {
	Text      string
	LastUsed  time.Time
	Frequency int
}

type Node struct {
	IsLeaf   bool
	Keys     []string
	Children []*Node
	Values   []*Command
	Next     *Node
}

func NewTree() *Bptree {
	root := &Node{
		IsLeaf:   true,
		Keys:     make([]string, 0, 10),
		Children: make([]*Node, 0, 10),
		Values:   make([]*Command, 0, 10),
	}
	return &Bptree{Root: root}
}

/*
        [10, 20]
       /    |    \
      C0   C1    C2

	  so 15 will be in c1 which is index 1 same as
	        [30]
         /    \
   [10 20]   [30 40]
*/

func (tree *Bptree) Search(command string) *Command {
	vals := tree.Root.Values
	index := sort.Search(len(vals), func(i int) bool {
		return tree.Root.Values[i].Text >= command
	})
	if index < len(vals) && vals[index].Text == command {
		return vals[index]
	}
	return nil
}

func (tree *Bptree) Insert(cmdText string) {
	command := &Command{
		Text:      cmdText,
		LastUsed:  time.Now(),
		Frequency: 1,
	}

	node := tree.Root
	if len(node.Values) == 0 {
		node.Values = append(node.Values, command)
		return
	}

	i := sort.Search(len(node.Values), func(i int) bool {
		return node.Values[i].Text > cmdText
	})

	node.Values = append(node.Values, nil)
	copy(node.Values[i+1:], node.Values[i:])
	node.Values[i] = command
}
