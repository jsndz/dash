package storage

import (
	"dash/internal/tree"
	"time"
)

type Disktree struct {
	RootID  uint64
	MaxSize int
	MinSize int
	Nodes   []*DiskNode
}

type Command struct {
	Text      string
	LastUsed  time.Time
	Frequency int
}

type DiskNode struct {
	ID         uint64
	IsLeaf     bool
	Keys       []uint64
	ChildrenID []uint64
	Values     []*Command
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

func Generate(tree *tree.Bptree) *Disktree {
	rootId := 1
	disktree := NewDisktree(uint64(rootId), 128, 64)

	return disktree
}
