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
		Keys:     make([]string, 0, MAX),
		Children: make([]*Node, 0, MAX+1),
		Values:   make([]*Command, 0, MAX),
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

func (tree *Bptree) Search(command string, node *Node) *Command {
	if node.IsLeaf {
		return tree.SearchLeafNode(command, node)
	}
	index := sort.Search(len(node.Keys), func(i int) bool {
		return node.Keys[i] > command
	})
	return tree.Search(command, node.Children[index])
}

func (tree *Bptree) SearchLeafNode(command string, node *Node) *Command {

	index := sort.Search(len(node.Keys), func(i int) bool {
		return node.Keys[i] >= command
	})
	if index < len(node.Keys) && node.Keys[index] == command {
		return node.Values[index]
	}
	return nil
}

func (tree *Bptree) InsertInRecursion() {

}

// from my idea of insert
// first lets find the node through searching
// if already exist increase the frequency
// if does not exist take its index and parent index
// and insert
// insert has 3 conditions
// insert and max size wont be crossed
// insert and max size will cross
// where you have to split the node and

func FindLeaf(node *Node, cmd string) *Node {
	if node.IsLeaf {
		return node
	}
	index := sort.Search(len(node.Keys), func(i int) bool {
		return node.Keys[i] > cmd
	})
	// from this moment there are two options
	// if the value we are searching is same as key
	// or the value we are searching is pointing to some children
	if len(node.Keys) > 1 {
	}
	return FindLeaf(node.Children[index], cmd)
}

func (tree *Bptree) Insert(cmdText string) {
	command := &Command{
		Text:      cmdText,
		LastUsed:  time.Now(),
		Frequency: 1,
	}

	node := tree.Root
	// you have root search for >= in the node keys
	leafNode := FindLeaf(node, cmdText)
	InsertIntoLeaf(leafNode, command)
}

func InsertIncreaseFrequency(command *Command) {
	command.Frequency++
	command.LastUsed = time.Now()
}

func InsertIntoLeaf(node *Node, command *Command) {
	if len(node.Values) == 0 {
		node.Values = append(node.Values, command)
		return
	}

	i := sort.Search(len(node.Values), func(i int) bool {
		return node.Values[i].Text >= command.Text
	})
	if i < len(node.Keys) && node.Keys[i] == command.Text {
		InsertIncreaseFrequency(command)
		return
	}

	// add keys also since they are direct pointer to values
	node.Keys = append(node.Keys, "")
	copy(node.Keys[i+1:], node.Keys[i:])
	node.Keys[i] = command.Text

	node.Values = append(node.Values, nil)
	copy(node.Values[i+1:], node.Values[i:])
	node.Values[i] = command

}
