package tree

import (
	"errors"
	"fmt"
	"sort"
	"strings"
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
	if node.IsLeaf {
		return node
	}
	index := sort.Search(len(node.Keys), func(i int) bool {
		return node.Keys[i] > cmd
	})

	return tree.FindLeaf(node.Children[index], cmd)
}

func (tree *Bptree) Insert(cmdText string) {
	node := tree.Root
	leafNode := tree.FindLeaf(node, cmdText)
	if len(leafNode.Keys) > tree.MaxSize {
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
	copy(parent.Children[i+2:], parent.Children[i+1:])
	parent.Children[i+1] = child
	child.Parent = parent
	if len(parent.Keys) > tree.MaxSize {
		tree.SplitAndInsertToInternal(parent)
	}
}

func (tree *Bptree) SplitAndInsertToInternal(node *Node) {
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
		tree.InsertToInternal(node.Parent, promotedKey, rightNode)
	} else {
		tree.CreateNewParentNode(node, rightNode, promotedKey)
	}
}

func (tree *Bptree) CreateNewParentNode(leftNode, rightNode *Node, promotedKey string) {
	parentNode := NewNode(false)
	leftNode.Parent = parentNode
	rightNode.Parent = parentNode
	parentNode.Keys = append(parentNode.Keys, promotedKey)
	parentNode.Children = append(parentNode.Children, leftNode, rightNode)

	tree.Root = parentNode
}

func (tree *Bptree) CreateNewParentNodeForLeaf(leftNode, rightNode *Node) {
	parentNode := NewNode(false)

	leftNode.Parent = parentNode
	rightNode.Parent = parentNode
	parentNode.Keys = append(parentNode.Keys, rightNode.Keys[0])
	parentNode.Children = append(parentNode.Children, leftNode, rightNode)

	tree.Root = parentNode
}

func (tree *Bptree) SplitAndInsertToLeaf(leafNode *Node, cmdText string) {
	tree.InsertIntoLeaf(leafNode, cmdText)
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

	if leafNode.Parent != nil {
		rightNode.Parent = leafNode.Parent
		tree.InsertToInternal(leafNode.Parent, rightNode.Keys[0], rightNode)
	} else {
		tree.CreateNewParentNodeForLeaf(leafNode, rightNode)
	}
}

func (tree *Bptree) InsertIncreaseFrequency(command *Command) {
	command.Frequency++
	command.LastUsed = time.Now()
}

func (tree *Bptree) InsertIntoLeaf(node *Node, cmdText string) {
	i := sort.Search(len(node.Keys), func(i int) bool {
		return node.Keys[i] >= cmdText
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

func (tree *Bptree) RefreshSeparator(node *Node, childIndex int) {
	if childIndex == 0 {
		return
	}
	newKey := node.Children[childIndex].Keys[0]
	node.Keys[childIndex-1] = newKey
}

func (tree *Bptree) GetSeparatorIndex(Parent *Node, childIndex int) string {
	return Parent.Keys[childIndex]
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

func (node *Node) getSiblings() (*Node, *Node) {
	if node.Parent == nil {
		return nil, nil
	}
	idx := ChildIndex(node)
	if idx == -1 {
		return nil, nil
	}
	var prev, next *Node
	if idx > 0 {
		prev = node.Parent.Children[idx-1]
	}
	if idx < len(node.Parent.Children)-1 {
		next = node.Parent.Children[idx+1]
	}
	return prev, next
}

func (tree *Bptree) updateSeparator(node *Node, oldKey, newKey string) {
	if oldKey == "" || oldKey == newKey {
		return
	}
	curr := node
	for curr.Parent != nil {
		parent := curr.Parent
		idx := ChildIndex(curr)
		if idx > 0 {
			if parent.Keys[idx-1] == oldKey {
				parent.Keys[idx-1] = newKey
				return
			}
		}
		curr = parent
	}
}

func (tree *Bptree) BorrowOrMerge(node *Node, oldKey string) {
	prev, next := node.getSiblings()

	if prev != nil && len(prev.Keys) > tree.MinSize {
		// Borrow from left sibling
		node.Keys = append(node.Keys, "")
		copy(node.Keys[1:], node.Keys[:len(node.Keys)-1])
		node.Values = append(node.Values, nil)
		copy(node.Values[1:], node.Values[:len(node.Values)-1])

		borrowedKey := prev.Keys[len(prev.Keys)-1]
		borrowedVal := prev.Values[len(prev.Keys)-1]
		prev.Keys = prev.Keys[:len(prev.Keys)-1]
		prev.Values = prev.Values[:len(prev.Values)-1]

		node.Keys[0] = borrowedKey
		node.Values[0] = borrowedVal

		childIndex := ChildIndex(node)
		node.Parent.Keys[childIndex-1] = node.Keys[0]
		return
	}

	if next != nil && len(next.Keys) > tree.MinSize {
		// Borrow from right sibling
		borrowedKey := next.Keys[0]
		borrowedVal := next.Values[0]

		node.Keys = append(node.Keys, borrowedKey)
		node.Values = append(node.Values, borrowedVal)

		copy(next.Keys[:len(next.Keys)-1], next.Keys[1:])
		next.Keys = next.Keys[:len(next.Keys)-1]
		copy(next.Values[:len(next.Values)-1], next.Values[1:])
		next.Values = next.Values[:len(next.Values)-1]

		childIndex := ChildIndex(next)
		node.Parent.Keys[childIndex-1] = next.Keys[0]

		tree.updateSeparator(node, oldKey, node.Keys[0])
		return
	}

	tree.MergeLeaf(node, oldKey)
}

func (tree *Bptree) BorrowOrMergeInternal(node *Node) {
	prev, next := node.getSiblings()

	if prev != nil && len(prev.Keys) > tree.MinSize {
		// Borrow from left sibling
		childIndex := ChildIndex(prev)
		separator := tree.GetSeparatorIndex(node.Parent, childIndex)

		node.Keys = append([]string{separator}, node.Keys...)
		node.Parent.Keys[childIndex] = prev.Keys[len(prev.Keys)-1]

		borrowedChild := prev.Children[len(prev.Children)-1]
		node.Children = append([]*Node{borrowedChild}, node.Children...)
		borrowedChild.Parent = node

		prev.Keys = prev.Keys[:len(prev.Keys)-1]
		prev.Children = prev.Children[:len(prev.Children)-1]
		return
	}

	if next != nil && len(next.Keys) > tree.MinSize {
		// Borrow from right sibling
		childIndex := ChildIndex(node)
		separator := tree.GetSeparatorIndex(node.Parent, childIndex)

		node.Keys = append(node.Keys, separator)
		borrowedChild := next.Children[0]
		node.Children = append(node.Children, borrowedChild)
		borrowedChild.Parent = node

		node.Parent.Keys[childIndex] = next.Keys[0]

		next.Keys = next.Keys[1:]
		next.Children = next.Children[1:]
		return
	}

	// Merge
	if prev != nil {
		childIndex := ChildIndex(prev)
		separator := tree.GetSeparatorIndex(node.Parent, childIndex)

		copy(node.Parent.Keys[childIndex:], node.Parent.Keys[childIndex+1:])
		node.Parent.Keys = node.Parent.Keys[:len(node.Parent.Keys)-1]
		copy(node.Parent.Children[childIndex+1:], node.Parent.Children[childIndex+2:])
		node.Parent.Children = node.Parent.Children[:len(node.Parent.Children)-1]

		prev.Keys = append(prev.Keys, separator)
		prev.Keys = append(prev.Keys, node.Keys...)
		prev.Children = append(prev.Children, node.Children...)
		for _, child := range node.Children {
			child.Parent = prev
		}

		if node.Parent != nil && node.Parent != tree.Root && len(node.Parent.Keys) < tree.MinSize {
			tree.BorrowOrMergeInternal(node.Parent)
		}
		return
	}

	if next != nil {
		childIndex := ChildIndex(node)
		separator := tree.GetSeparatorIndex(node.Parent, childIndex)

		node.Keys = append(node.Keys, separator)
		for _, child := range next.Children {
			child.Parent = node
		}
		node.Keys = append(node.Keys, next.Keys...)
		node.Children = append(node.Children, next.Children...)

		copy(node.Parent.Keys[childIndex:], node.Parent.Keys[childIndex+1:])
		node.Parent.Keys = node.Parent.Keys[:len(node.Parent.Keys)-1]
		copy(node.Parent.Children[childIndex+1:], node.Parent.Children[childIndex+2:])
		node.Parent.Children = node.Parent.Children[:len(node.Parent.Children)-1]

		if node.Parent != nil && node.Parent != tree.Root && len(node.Parent.Keys) < tree.MinSize {
			tree.BorrowOrMergeInternal(node.Parent)
		}
		return
	}
}

func (tree *Bptree) MergeLeaf(node *Node, oldKey string) {
	prev, next := node.getSiblings()

	if prev != nil {
		leftNode := prev
		leftNode.Keys = append(leftNode.Keys, node.Keys...)
		leftNode.Values = append(leftNode.Values, node.Values...)

		childIndex := ChildIndex(node)

		node.Parent.Keys = append(node.Parent.Keys[:childIndex-1], node.Parent.Keys[childIndex:]...)
		node.Parent.Children = append(node.Parent.Children[:childIndex], node.Parent.Children[childIndex+1:]...)

		leftNode.Next = node.Next
		if node.Next != nil {
			node.Next.Prev = leftNode
		}

		if leftNode.Parent != nil && leftNode.Parent != tree.Root && len(leftNode.Parent.Keys) < tree.MinSize {
			tree.BorrowOrMergeInternal(leftNode.Parent)
		}
		return
	}

	if next != nil {
		nextNode := next
		node.Keys = append(node.Keys, nextNode.Keys...)
		node.Values = append(node.Values, nextNode.Values...)

		childIndex := ChildIndex(node)

		node.Parent.Keys = append(node.Parent.Keys[:childIndex], node.Parent.Keys[childIndex+1:]...)
		node.Parent.Children = append(node.Parent.Children[:childIndex+1], node.Parent.Children[childIndex+2:]...)

		node.Next = nextNode.Next
		if nextNode.Next != nil {
			nextNode.Next.Prev = node
		}

		tree.updateSeparator(node, oldKey, node.Keys[0])

		if node.Parent != nil && node.Parent != tree.Root && len(node.Parent.Keys) < tree.MinSize {
			tree.BorrowOrMergeInternal(node.Parent)
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

	oldKey := node.Keys[0]

	// Delete from leaf
	copy(node.Keys[i:], node.Keys[i+1:])
	copy(node.Values[i:], node.Values[i+1:])
	node.Keys = node.Keys[:len(node.Keys)-1]
	node.Values = node.Values[:len(node.Values)-1]

	if node != tree.Root && len(node.Keys) < tree.MinSize {
		tree.BorrowOrMerge(node, oldKey)
	} else {
		if len(node.Keys) > 0 && node.Keys[0] != oldKey {
			tree.updateSeparator(node, oldKey, node.Keys[0])
		}
	}

	if tree.Root.Parent != nil {
		tree.Root.Parent = nil
	}

	if !tree.Root.IsLeaf && len(tree.Root.Keys) == 0 {
		tree.Root = tree.Root.Children[0]
		tree.Root.Parent = nil
	}
	return nil
}

func (tree *Bptree) RangeScan(prefix string) []*Command {
	node := tree.FindLeaf(tree.Root, prefix)
	if node == nil {
		return nil
	}

	i := sort.Search(len(node.Keys), func(i int) bool {
		return node.Keys[i] >= prefix
	})

	var result []*Command
	for node != nil {
		for i < len(node.Keys) {
			if !strings.HasPrefix(node.Keys[i], prefix) {
				return result
			}
			result = append(result, node.Values[i])
			i++
		}
		node = node.Next
		i = 0
	}
	return result
}

func validate(node *Node) {
	if node.IsLeaf {
		return
	}

	if len(node.Children) != len(node.Keys)+1 {
		panic(fmt.Sprintf(
			"invalid node: keys=%v keys=%d children=%d",
			node.Keys,
			len(node.Keys),
			len(node.Children),
		))
	}

	for _, child := range node.Children {
		if child.Parent != node {
			panic("bad parent pointer")
		}
		validate(child)
	}
}
