package storage

import (
	"bytes"
	"dash/internal/tree"
	"encoding/binary"
)

func SerializeDisktree(tree *Disktree) ([]byte, error) {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, tree.RootID)
	binary.Write(&buf, binary.LittleEndian, uint32(tree.MaxSize))
	binary.Write(&buf, binary.LittleEndian, uint32(tree.MinSize))
	binary.Write(&buf, binary.LittleEndian, uint32(len(tree.Nodes)))
	for _, node := range tree.Nodes {
		b, _ := SerializeDiskNode(node)
		binary.Write(&buf, binary.LittleEndian, uint32(len(b)))
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func SerializeDiskNode(node *DiskNode) ([]byte, error) {
	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, node.ID)
	if node.IsLeaf {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}
	binary.Write(&buf, binary.LittleEndian, node.ParentID)
	binary.Write(&buf, binary.LittleEndian, node.PrevID)
	binary.Write(&buf, binary.LittleEndian, node.NextID)
	binary.Write(&buf, binary.LittleEndian, uint32(len(node.Keys)))
	for _, key := range node.Keys {
		binary.Write(&buf, binary.LittleEndian, uint32(len(key)))
		buf.WriteString(key)
	}
	if node.IsLeaf {
		binary.Write(&buf, binary.LittleEndian, uint32(len(node.Values)))
		for _, val := range node.Values {
			serializedVal := SerializeCommand(val)
			binary.Write(&buf, binary.LittleEndian, len(serializedVal))
			buf.Write(serializedVal)
		}
	} else {
		binary.Write(&buf, binary.LittleEndian, uint32(len(node.ChildrenID)))
		for _, childID := range node.ChildrenID {
			binary.Write(&buf, binary.LittleEndian, childID)
		}
	}
	return buf.Bytes(), nil
}
func SerializeCommand(cmd tree.Command) []byte {
	var buf bytes.Buffer
	text := []byte(cmd.Text)
	binary.Write(&buf, binary.LittleEndian, cmd.LastUsed.Unix())
	binary.Write(&buf, binary.LittleEndian, uint64(cmd.Frequency))
	binary.Write(&buf, binary.LittleEndian, uint32(len(text)))
	buf.Write(text)

	return buf.Bytes()
}
