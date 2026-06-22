package tree

import (
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

func FindLeaf(node *Node, cmd string) *Node {
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
	return FindLeaf(node.Children[index], cmd)
}

func (tree *Bptree) Insert(cmdText string) {
	log.Printf("Insert: cmdText=%q Root.Keys=%v", cmdText, tree.Root.Keys)
	node := tree.Root
	// you have root search for >= in the node keys
	leafNode := FindLeaf(node, cmdText)
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
