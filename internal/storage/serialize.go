package storage

import (
	"bytes"
	"dash/internal/tree"
	"encoding/binary"
)

func SerializeDisktree(tree *Disktree) ([]byte, error) {
	var buf bytes.Buffer
	for _, node := range tree.Nodes {
		serializedNode, _ := SerializeDiskNode(node)
		buf.Write(serializedNode)
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
	binary.Write(&buf, binary.LittleEndian, len(node.Keys))
	for _, key := range node.Keys {
		binary.Write(&buf, binary.LittleEndian, len(key))
		buf.WriteString(key)
	}
	if node.IsLeaf {
		for _, val := range node.Values {
			serializedVal := SerializeCommand(val)
			binary.Write(&buf, binary.LittleEndian, len(serializedVal))
			buf.Write(serializedVal)
		}
	} else {
		for _, childID := range node.ChildrenID {
			binary.Write(&buf, binary.LittleEndian, childID)
		}
	}
	return buf.Bytes(), nil
}
func SerializeCommand(cmd tree.Command) []byte {
	var buf bytes.Buffer
	text := []byte(cmd.Text)

	binary.Write(&buf, binary.LittleEndian, uint32(len(text)))
	buf.Write(text)
	binary.Write(&buf, binary.LittleEndian, cmd.LastUsed.Unix())
	binary.Write(&buf, binary.LittleEndian, int64(cmd.Frequency))

	return buf.Bytes()
}
