package storage

import (
	"dash/internal/tree"
)

type Disktree struct {
	RootID  uint64
	MaxSize int
	MinSize int
	Nodes   []*DiskNode
}

type DiskNode struct {
	ID         uint64
	IsLeaf     bool
	Keys       []string
	ChildrenID []uint64
	Values     []tree.Command
	ParentID   uint64
	NextID     uint64
	PrevID     uint64
}

func NewDisktree(rootID uint64, maxSize int, minSize int) *Disktree {
	return &Disktree{
		RootID:  rootID,
		MaxSize: maxSize,
		MinSize: minSize,
		Nodes:   []*DiskNode{},
	}
}

func Generate(t *tree.Bptree) *Disktree {
	rootId := 1
	disktree := NewDisktree(uint64(rootId), 128, 64)
	mapper := make(map[*tree.Node]uint64)
	id := uint64(1)
	SetupMap(mapper, t.Root, &id)
	traverse(t.Root, disktree, mapper)
	return disktree
}

func SetupMap(treemap map[*tree.Node]uint64, node *tree.Node, id *uint64) {
	treemap[node] = *id
	*id++

	for _, node := range node.Children {
		SetupMap(treemap, node, id)
	}

}

func traverse(node *tree.Node, disktree *Disktree, mapper map[*tree.Node]uint64) {
	var children []uint64
	var vals []tree.Command

	if node.IsLeaf {
		for _, cmd := range node.Values {
			vals = append(vals, *cmd)
		}
	} else {
		for _, ch := range node.Children {
			children = append(children, mapper[ch])
		}
	}
	treeNode := &DiskNode{
		ID:         mapper[node],
		IsLeaf:     node.IsLeaf,
		Keys:       node.Keys,
		ChildrenID: children,
		Values:     vals,
		ParentID:   mapper[node.Parent],
		PrevID:     mapper[node.Prev],
		NextID:     mapper[node.Next],
	}
	if !node.IsLeaf {
		for _, ch := range node.Children {
			traverse(ch, disktree, mapper)
		}
	}
	disktree.Nodes = append(disktree.Nodes, treeNode)
}
