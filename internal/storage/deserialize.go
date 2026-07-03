package storage

import (
	"bytes"
	"dash/internal/tree"
	"encoding/binary"
	"io"
	"time"
)

func DeserializeDisktree(b []byte) *Disktree {
	var tree Disktree
	r := bytes.NewReader(b)
	binary.Read(r, binary.LittleEndian, &tree.RootID)
	var maxSize, minSize uint32
	binary.Read(r, binary.LittleEndian, &maxSize)
	binary.Read(r, binary.LittleEndian, &minSize)
	tree.MaxSize = int(maxSize)
	tree.MinSize = int(minSize)
	var nodeCount uint32
	binary.Read(r, binary.LittleEndian, &nodeCount)
	tree.Nodes = make([]*DiskNode, nodeCount)
	for i := 0; i < int(nodeCount); i++ {
		var nodeLen uint32
		binary.Read(r, binary.LittleEndian, &nodeLen)
		nodeBytes := make([]byte, nodeLen)
		io.ReadFull(r, nodeBytes)
		tree.Nodes[i] = DeserializeDiskNode(nodeBytes)
	}
	return &tree
}

func DeserializeDiskNode(b []byte) *DiskNode {
	var node DiskNode
	r := bytes.NewReader(b)
	binary.Read(r, binary.LittleEndian, &node.ID)
	isLeaf, _ := r.ReadByte()
	node.IsLeaf = isLeaf == 1
	binary.Read(r, binary.LittleEndian, &node.ParentID)
	binary.Read(r, binary.LittleEndian, &node.PrevID)
	binary.Read(r, binary.LittleEndian, &node.NextID)
	var keyCount uint32
	binary.Read(r, binary.LittleEndian, &keyCount)
	node.Keys = make([]string, keyCount)
	for i := 0; i < int(keyCount); i++ {
		var keyLen uint32
		binary.Read(r, binary.LittleEndian, &keyLen)
		keyBytes := make([]byte, keyLen)
		io.ReadFull(r, keyBytes)
		node.Keys[i] = string(keyBytes)
	}
	if node.IsLeaf {
		var valCount uint32
		binary.Read(r, binary.LittleEndian, &valCount)
		node.Values = make([]tree.Command, valCount)
		for i := 0; i < int(valCount); i++ {
			var valLen int
			binary.Read(r, binary.LittleEndian, &valLen)
			valBytes := make([]byte, valLen)
			io.ReadFull(r, valBytes)
			node.Values[i] = *DeserializeCommand(valBytes)
		}
	} else {
		var childCount uint32
		binary.Read(r, binary.LittleEndian, &childCount)
		node.ChildrenID = make([]uint64, childCount)
		for i := 0; i < int(childCount); i++ {
			binary.Read(r, binary.LittleEndian, &node.ChildrenID[i])
		}
	}
	return &node
}
func DeserializeCommand(b []byte) *tree.Command {

	var cmdLen uint32
	var ts int64
	var frq uint64
	r := bytes.NewReader(b)
	binary.Read(r, binary.LittleEndian, &ts)
	binary.Read(r, binary.LittleEndian, &frq)
	binary.Read(r, binary.LittleEndian, &cmdLen)
	text := make([]byte, cmdLen)
	io.ReadFull(r, text)
	cmd := tree.Command{
		Text:      string(text),
		LastUsed:  time.Unix(ts, 0),
		Frequency: int(frq),
	}
	return &cmd
}
