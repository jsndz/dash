package bptree

import (
	"sort"
	"time"
)

type Bptree struct {
	Root    *Node
	MaxSize int
	MinSize int
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
	Parent   *Node
	Next     *Node
}

func NewNode(isLeaf bool) *Node {
	// split logic is in the code not in array constraint
	return &Node{
		IsLeaf:   isLeaf,
		Keys:     make([]string, 0),
		Values:   make([]*Command, 0),
		Children: make([]*Node, 0),
	}
}

func NewTree(maxSize, minSize int) *Bptree {
	root := NewNode(true)
	return &Bptree{Root: root, MaxSize: maxSize, MinSize: minSize}
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

	node := tree.Root
	// you have root search for >= in the node keys
	leafNode := FindLeaf(node, cmdText)
	if len(leafNode.Keys) == tree.MaxSize {
		tree.SplitAndInsertToLeaf(leafNode, cmdText)
		return
	}
	tree.InsertIntoLeaf(leafNode, cmdText)

}

func (tree *Bptree) InsertToInternal(parent *Node, key string, child *Node) {
	if parent.IsLeaf {
		return
	}
	// add the key to the keys
	i := sort.Search(len(parent.Keys), func(i int) bool {
		return parent.Keys[i] > key
	})
	parent.Keys = append(parent.Keys, "")
	copy(parent.Keys[i+1:], parent.Keys[i:])
	parent.Keys[i] = key

	parent.Children = append(parent.Children, nil)
	copy(parent.Children[i+1:], parent.Children[i:])
	parent.Children[i] = child
	if parent.Parent != nil && len(parent.Keys) > tree.MaxSize {
		// split the parent node and insert the key to the grandparent
		tree.SplitAndInsertToInternal(parent.Parent)
	}
}

func (tree *Bptree) SplitAndInsertToInternal(node *Node) {
	rightNode := NewNode(false)
	mid := len(node.Keys) / 2

	rightNode.Keys = append(rightNode.Keys, node.Keys[mid:]...)
	node.Keys = node.Keys[:mid]

	rightNode.Values = append(rightNode.Values, node.Values[mid:]...)
	node.Values = node.Values[:mid]

	if node.Parent != nil {
		tree.InsertToInternal(node.Parent, rightNode.Keys[0], rightNode)
	} else {
		tree.CreateNewParentNode(node, rightNode)
	}
}

func (tree *Bptree) CreateNewParentNode(leftNode, rightNode *Node) {
	parentNode := NewNode(false)
	leftNode.Parent = parentNode
	rightNode.Parent = parentNode
	parentNode.Keys = append(parentNode.Keys, rightNode.Keys[0])
	parentNode.Children = append(parentNode.Children, leftNode, rightNode)
}
func (tree *Bptree) SplitAndInsertToLeaf(leafNode *Node, cmdText string) {
	// insert to the leaf node then split

	// if any node is spliting it should start with leaf node
	rightNode := NewNode(true)
	mid := len(leafNode.Keys) / 2

	//keys
	rightNode.Keys = append(rightNode.Keys, leafNode.Keys[mid:]...)
	leafNode.Keys = leafNode.Keys[:mid]
	//values
	rightNode.Values = append(rightNode.Values, leafNode.Values[mid:]...)
	leafNode.Values = leafNode.Values[:mid]

	// two nodes are ready leftnode(leafnode) and rightnode
	// now i need to add the left most in the right node to the parent node
	// if parent exist add to it
	if leafNode.Parent != nil {
		rightNode.Parent = leafNode.Parent
		tree.InsertToInternal(leafNode.Parent, rightNode.Keys[0], rightNode)
	} else {
		tree.CreateNewParentNode(leafNode, rightNode)
	}
}

func (tree *Bptree) InsertIncreaseFrequency(command *Command) {
	command.Frequency++
	command.LastUsed = time.Now()
}

func (tree *Bptree) InsertIntoLeaf(node *Node, cmdText string) {

	i := sort.Search(len(node.Values), func(i int) bool {
		return node.Values[i].Text >= cmdText
	})
	// if key exist directly increase frequency
	if i < len(node.Keys) && node.Keys[i] == cmdText {
		tree.InsertIncreaseFrequency(node.Values[i])
		return
	}

	command := &Command{
		Text:      cmdText,
		LastUsed:  time.Now(),
		Frequency: 1,
	}

	// add keys also since they are direct pointer to values
	node.Keys = append(node.Keys, "")
	copy(node.Keys[i+1:], node.Keys[i:])
	node.Keys[i] = command.Text

	node.Values = append(node.Values, nil)
	copy(node.Values[i+1:], node.Values[i:])
	node.Values[i] = command

}
