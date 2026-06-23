package tree

import (
	"errors"
	"log"
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
	Prev     *Node
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
       /    |    \ujn6y6y
      C0   C1    C2

	  so 15 will be in c1 which is index 1 same as
	        [30]
         /    \
   [10 20]   [30 40]
*/

func (tree *Bptree) Search(command string, node *Node) *Command {
	log.Printf("Search: command=%q keys=%v leaf=%v",
		command, node.Keys, node.IsLeaf)

	if node.IsLeaf {
		return tree.SearchLeafNode(command, node)
	}

	index := sort.Search(len(node.Keys), func(i int) bool {
		result := node.Keys[i] > command
		log.Printf("Internal compare: key=%q > %q => %v",
			node.Keys[i], command, result)
		return result
	})

	log.Printf("Descending to child index=%d (children=%d)",
		index, len(node.Children))

	return tree.Search(command, node.Children[index])
}

func (tree *Bptree) SearchLeafNode(command string, node *Node) *Command {
	log.Printf("Leaf search: command=%q keys=%v",
		command, node.Keys)

	index := sort.Search(len(node.Keys), func(i int) bool {
		result := node.Keys[i] >= command
		log.Printf("Leaf compare: key=%q >= %q => %v",
			node.Keys[i], command, result)
		return result
	})

	log.Printf("Leaf search result index=%d", index)

	if index < len(node.Keys) && node.Keys[index] == command {
		log.Printf("Found command at index=%d", index)
		return node.Values[index]
	}

	log.Printf("Command not found")
	return nil
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

func (tree *Bptree) FindLeaf(node *Node, cmd string) *Node {
	log.Printf("FindLeaf: cmd=%q keys=%v leaf=%v", cmd, node.Keys, node.IsLeaf)
	if node.IsLeaf {
		return node
	}
	index := sort.Search(len(node.Keys), func(i int) bool {
		return node.Keys[i] > cmd
	})
	// from this moment there are two options
	// if the value we are searching is same as key
	// or the value we are searching is pointing to some children

	log.Printf("FindLeaf: descending to child index=%d", index)
	return tree.FindLeaf(node.Children[index], cmd)
}

func (tree *Bptree) Insert(cmdText string) {
	log.Printf("Insert: cmdText=%q Root.Keys=%v", cmdText, tree.Root.Keys)
	node := tree.Root
	// you have root search for >= in the node keys
	leafNode := tree.FindLeaf(node, cmdText)
	if len(leafNode.Keys) > tree.MaxSize {
		log.Printf("Insert: leafNode at max size, splitting")
		tree.SplitAndInsertToLeaf(leafNode, cmdText)
		return
	}
	log.Printf("Insert: inserting into leaf without split")
	tree.InsertIntoLeaf(leafNode, cmdText)

}

func (tree *Bptree) InsertToInternal(parent *Node, key string, child *Node) {
	log.Printf("InsertToInternal: parent.Keys=%v key=%q", parent.Keys, key)
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
	copy(parent.Children[i+2:], parent.Children[i+1:])
	parent.Children[i+1] = child
	if len(parent.Keys) > tree.MaxSize {

		log.Printf("InsertToInternal: parent keys exceed MaxSize, splitting parent")
		tree.SplitAndInsertToInternal(parent)

	}
}

func (tree *Bptree) SplitAndInsertToInternal(node *Node) {
	log.Printf("SplitAndInsertToInternal: node.Keys=%v", node.Keys)
	rightNode := NewNode(false)
	mid := len(node.Keys) / 2
	promotedKey := node.Keys[mid]
	rightNode.Keys = append(rightNode.Keys, node.Keys[mid+1:]...)
	node.Keys = node.Keys[:mid]

	rightNode.Children = append(rightNode.Children, node.Children[mid+1:]...)
	node.Children = node.Children[:mid+1]

	for _, child := range rightNode.Children {
		child.Parent = rightNode
	}
	if node.Parent != nil {
		log.Printf("SplitAndInsertToInternal: transferring key=%q to parent", promotedKey)
		tree.InsertToInternal(node.Parent, promotedKey, rightNode)
	} else {
		log.Printf("SplitAndInsertToInternal: creating new root parent node")
		tree.CreateNewParentNode(node, rightNode, promotedKey)
	}
}

func (tree *Bptree) CreateNewParentNode(leftNode, rightNode *Node, promotedKey string) {
	log.Printf("CreateNewParentNode: leftNode.Keys=%v, rightNode.Keys=%v", leftNode.Keys, rightNode.Keys)
	parentNode := NewNode(false)
	leftNode.Parent = parentNode
	rightNode.Parent = parentNode
	parentNode.Keys = append(parentNode.Keys, promotedKey)
	parentNode.Children = append(parentNode.Children, leftNode, rightNode)

	tree.Root = parentNode
	log.Printf("CreateNewParentNode: new root created with keys=%v", tree.Root.Keys)
}

func (tree *Bptree) CreateNewParentNodeForLeaf(leftNode, rightNode *Node) {
	log.Printf("CreateNewParentNode: leftNode.Keys=%v, rightNode.Keys=%v", leftNode.Keys, rightNode.Keys)
	parentNode := NewNode(false)

	leftNode.Parent = parentNode
	rightNode.Parent = parentNode
	parentNode.Keys = append(parentNode.Keys, rightNode.Keys[0])
	parentNode.Children = append(parentNode.Children, leftNode, rightNode)

	tree.Root = parentNode
	log.Printf("CreateNewParentNode: new root created with keys=%v", tree.Root.Keys)
}

func (tree *Bptree) SplitAndInsertToLeaf(leafNode *Node, cmdText string) {
	// insert to the leaf node then split
	log.Printf("SplitAndInsertToLeaf: leafNode.Keys=%v, cmdText=%q", leafNode.Keys, cmdText)

	tree.InsertIntoLeaf(leafNode, cmdText)
	// if any node is spliting it should start with leaf node
	rightNode := NewNode(true)
	mid := len(leafNode.Keys) / 2

	//keys
	rightNode.Keys = append(rightNode.Keys, leafNode.Keys[mid:]...)
	leafNode.Keys = leafNode.Keys[:mid]
	//values
	rightNode.Values = append(rightNode.Values, leafNode.Values[mid:]...)
	leafNode.Values = leafNode.Values[:mid]

	rightNode.Next = leafNode.Next
	leafNode.Next = rightNode
	rightNode.Prev = leafNode
	// two nodes are ready leftnode(leafnode) and rightnode
	// now i need to add the left most in the right node to the parent node
	// if parent exist add to it
	if leafNode.Parent != nil {
		log.Printf("SplitAndInsertToLeaf: parent exists, inserting rightNode.Keys[0]=%q to parent", rightNode.Keys[0])
		rightNode.Parent = leafNode.Parent
		tree.InsertToInternal(leafNode.Parent, rightNode.Keys[0], rightNode)
	} else {
		log.Printf("SplitAndInsertToLeaf: parent does not exist, creating new parent node with key=%q", rightNode.Keys[0])
		tree.CreateNewParentNodeForLeaf(leafNode, rightNode)
	}
}

func (tree *Bptree) InsertIncreaseFrequency(command *Command) {
	log.Printf("InsertIncreaseFrequency: command=%q, frequency=%d", command.Text, command.Frequency)
	command.Frequency++
	command.LastUsed = time.Now()
}

func (tree *Bptree) InsertIntoLeaf(node *Node, cmdText string) {
	log.Printf("InsertIntoLeaf: node.Keys=%v, cmdText=%q", node.Keys, cmdText)
	i := sort.Search(len(node.Keys), func(i int) bool {
		return node.Keys[i] >= cmdText
	})
	// if key exist directly increase frequency
	if i < len(node.Keys) && node.Keys[i] == cmdText {
		log.Printf("InsertIntoLeaf: command already exists at index=%d, increasing frequency", i)
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
	log.Printf("InsertIntoLeaf: command inserted at index=%d", i)
}

func (tree *Bptree) RefreshSeparator(node *Node, childIndex int) {
	if childIndex == 0 {
		return
	}
	newKey := node.Children[childIndex].Keys[0]
	node.Keys[childIndex-1] = newKey
}
func ChildIndex(node *Node) int {
	if node.Parent == nil {
		return -1
	}
	parent := node.Parent

	for i, child := range parent.Children {
		if child == node {
			return i
		}
	}

	return -1
}

func (tree *Bptree) BorrowOrMerge(node *Node) {
	// first try borrowing left
	if node.Prev != nil && len(node.Prev.Keys) > tree.MinSize {

		node.Keys = append(node.Keys, "")
		copy(node.Keys[1:], node.Keys[:len(node.Keys)-1])
		node.Values = append(node.Values, nil)
		copy(node.Values[1:], node.Values[:len(node.Values)-1])
		leftNode := node.Prev
		borrowedKey := leftNode.Keys[len(leftNode.Keys)-1]
		borrowedVal := leftNode.Values[len(leftNode.Keys)-1]
		leftNode.Keys = leftNode.Keys[:len(leftNode.Keys)-1]
		leftNode.Values = leftNode.Values[:len(leftNode.Values)-1]

		node.Keys[0] = borrowedKey
		node.Values[0] = borrowedVal
		childIndex := ChildIndex(node)
		tree.RefreshSeparator(node.Parent, childIndex)
		return
	} else if node.Next != nil && len(node.Next.Keys) > tree.MinSize {

		node.Keys = append(node.Keys, "")
		node.Values = append(node.Values, nil)
		nextNode := node.Next
		borrowedKey := nextNode.Keys[0]
		borrowedVal := nextNode.Values[0]
		//removing the borrowed node from the rightnode
		copy(nextNode.Keys[:len(nextNode.Keys)-1], nextNode.Keys[1:])
		nextNode.Keys = nextNode.Keys[:len(nextNode.Keys)-1]
		copy(nextNode.Values[:len(nextNode.Values)-1], nextNode.Values[1:])
		nextNode.Values = nextNode.Values[:len(nextNode.Values)-1]

		node.Keys[len(node.Keys)-1] = borrowedKey
		node.Values[len(node.Values)-1] = borrowedVal

		childIndex := ChildIndex(nextNode)

		tree.RefreshSeparator(node.Parent, childIndex)
		return
	}
	tree.MergeLeaf(node)
}

func (tree *Bptree) BorrowOrMergeInternal(node *Node) {
	//borrow
	if node.Prev != nil && len(node.Prev.Keys) > tree.MinSize {

		node.Keys = append(node.Keys, "")
		copy(node.Keys[1:], node.Keys[:len(node.Keys)-1])
		node.Children = append(node.Children, nil)
		copy(node.Children[1:], node.Children[:len(node.Children)-1])
		leftNode := node.Prev
		borrowedKey := leftNode.Keys[len(leftNode.Keys)-1]
		borrowedChildren := leftNode.Children[len(leftNode.Keys)-1]
		leftNode.Keys = leftNode.Keys[:len(leftNode.Keys)-1]
		leftNode.Children = leftNode.Children[:len(leftNode.Children)-1]

		node.Keys[0] = borrowedKey
		node.Children[0] = borrowedChildren
		childIndex := ChildIndex(node)
		tree.RefreshSeparator(node.Parent, childIndex)
		return
	} else if node.Next != nil && len(node.Next.Keys) > tree.MinSize {

		node.Keys = append(node.Keys, "")
		node.Children = append(node.Children, nil)
		nextNode := node.Next
		borrowedKey := nextNode.Keys[0]
		borrowedChildren := nextNode.Children[0]
		//removing the borrowed node from the rightnode
		copy(nextNode.Keys[:len(nextNode.Keys)-1], nextNode.Keys[1:])
		nextNode.Keys = nextNode.Keys[:len(nextNode.Keys)-1]
		copy(nextNode.Children[:len(nextNode.Children)-1], nextNode.Children[1:])
		nextNode.Children = nextNode.Children[:len(nextNode.Children)-1]

		node.Keys[len(node.Keys)-1] = borrowedKey
		node.Children[len(node.Children)-1] = borrowedChildren

		childIndex := ChildIndex(nextNode)

		tree.RefreshSeparator(node.Parent, childIndex)
		return
	}
	if node.Prev != nil {
		leftNode := node.Prev
		leftNode.Keys = append(leftNode.Keys, node.Keys...)
		leftNode.Children = append(leftNode.Children, node.Children...)

		childIndex := ChildIndex(node)

		node.Parent.Keys = append(node.Parent.Keys[:childIndex-1], node.Parent.Keys[childIndex:]...)
		node.Parent.Children = append(node.Parent.Children[:childIndex], node.Parent.Children[childIndex+1:]...)
		leftNode.Next = node.Next
		if node.Next != nil {
			node.Next.Prev = leftNode
		}
		return
	}
	if node.Next != nil {
		nextNode := node.Next
		node.Keys = append(node.Keys, nextNode.Keys...)
		node.Children = append(node.Children, nextNode.Children...)

		// just remove the separator the position will come into place
		childIndex := ChildIndex(node)

		node.Parent.Keys = append(node.Parent.Keys[:childIndex], node.Parent.Keys[childIndex+1:]...)
		node.Parent.Children = append(node.Parent.Children[:childIndex+1], node.Parent.Children[childIndex+2:]...)

		node.Next = nextNode.Next
		if nextNode.Next != nil {
			nextNode.Next.Prev = node
		}
		return
	}
	if node.Parent != nil && len(node.Parent.Keys) < tree.MinSize {
		tree.BorrowOrMergeInternal(node.Parent)
	}
}
func (tree *Bptree) MergeLeaf(node *Node) {
	if node.Prev != nil {
		leftNode := node.Prev
		leftNode.Keys = append(leftNode.Keys, node.Keys...)
		leftNode.Values = append(leftNode.Values, node.Values...)

		childIndex := ChildIndex(node)

		node.Parent.Keys = append(node.Parent.Keys[:childIndex-1], node.Parent.Keys[childIndex:]...)
		node.Parent.Children = append(node.Parent.Children[:childIndex], node.Parent.Children[childIndex+1:]...)
		leftNode.Next = node.Next
		if node.Next != nil {
			node.Next.Prev = leftNode
		}
		if leftNode.Parent != nil && len(leftNode.Parent.Keys) < tree.MinSize {
			tree.BorrowOrMergeInternal(leftNode.Parent)
		}
		return
	}
	if node.Next != nil {
		nextNode := node.Next
		node.Keys = append(node.Keys, nextNode.Keys...)
		node.Values = append(node.Values, nextNode.Values...)

		// just remove the separator the position will come into place
		childIndex := ChildIndex(node)

		node.Parent.Keys = append(node.Parent.Keys[:childIndex], node.Parent.Keys[childIndex+1:]...)
		node.Parent.Children = append(node.Parent.Children[:childIndex+1], node.Parent.Children[childIndex+2:]...)

		node.Next = nextNode.Next
		if nextNode.Next != nil {
			nextNode.Next.Prev = node
		}
		if nextNode.Parent != nil && len(nextNode.Parent.Keys) < tree.MinSize {
			tree.BorrowOrMergeInternal(nextNode.Parent)
		}
		return
	}
}
func (tree *Bptree) Delete(cmd string) error {
	node := tree.FindLeaf(tree.Root, cmd)

	i := sort.Search(len(node.Keys), func(i int) bool {
		return node.Keys[i] >= cmd
	})
	if i >= len(node.Keys) || node.Keys[i] != cmd {
		return errors.New("Command not found")
	}

	copy(node.Keys[i:], node.Keys[i+1:])
	copy(node.Values[i:], node.Values[i+1:])
	node.Keys = node.Keys[:len(node.Keys)-1]
	node.Values = node.Values[:len(node.Values)-1]

	if node.Parent != nil {
		node.Parent.Keys = append(node.Parent.Keys[:i], node.Parent.Keys[i+1:]...)
		node.Parent.Children = append(node.Parent.Children[:i], node.Parent.Children[i+1:]...)
	}
	if node.Parent == nil {
		return nil
	}
	if node != tree.Root && len(node.Keys) < tree.MinSize {
		tree.BorrowOrMerge(node)
	}
	return nil
}

func validate(node *Node) {
	if node.IsLeaf {
		return
	}

	if len(node.Children) != len(node.Keys)+1 {
		log.Fatalf(
			"invalid node: keys=%v keys=%d children=%d",
			node.Keys,
			len(node.Keys),
			len(node.Children),
		)
	}

	for _, child := range node.Children {
		validate(child)
	}
}
