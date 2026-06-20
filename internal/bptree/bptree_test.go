package bptree

import (
	"testing"
)

func TestRootAsLeaf(t *testing.T) {
	t.Run("check if basic insert is working for 1 element", func(t *testing.T) {
		tree := NewTree()

		tree.Insert("ls")

		cmd := tree.Search("ls", tree.Root)
		expected := "ls"
		if cmd.Text != expected {
			t.Errorf("didn't get expected form search for %s got %s", expected, cmd.Text)
		}
	})
	t.Run("check if basic insert is working for more element", func(t *testing.T) {
		tree := NewTree()

		tree.Insert("ls")
		tree.Insert("echo")
		tree.Insert("docker ps")

		cmd := tree.Search("ls", tree.Root)
		expected := "ls"
		if cmd.Text != expected {
			t.Errorf("didn't get expected form search for %s got %s", expected, cmd.Text)
		}
	})
	t.Run("search for non-existing element", func(t *testing.T) {
		tree := NewTree()

		tree.Insert("ls")
		cmd := tree.SearchLeafNode("la", tree.Root)

		if cmd != nil {
			t.Errorf(" expected nil search  %s got %v", "la", cmd)
		}
	})
}

// func TestSearchInternalNode(t *testing.T) {
// 	t.Run("")
// }
