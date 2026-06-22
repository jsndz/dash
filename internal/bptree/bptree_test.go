package bptree

import (
	"testing"
)

func TestRootAsLeaf(t *testing.T) {
	t.Run("check if basic insert is working for 1 element", func(t *testing.T) {
		tree := NewTree()

		node := &Node{
			IsLeaf: true,
			Keys:   []string{"ls"},
			Values: []*Command{{Text: "ls"}},
		}
		tree.Root = node
		cmd := tree.SearchLeafNode("ls", tree.Root)
		expected := "ls"
		if cmd.Text != expected {
			t.Errorf("didn't get expected form search for %s got %s", expected, cmd.Text)
		}
	})
	t.Run("check if basic insert is working for more element", func(t *testing.T) {
		tree := NewTree()

		node := &Node{
			IsLeaf: true,
			Keys:   []string{"ls"},
			Values: []*Command{{Text: "ls"}, {Text: "la"}, {Text: "cd"}},
		}
		tree.Root = node
		cmd := tree.SearchLeafNode("ls", tree.Root)
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

func TestSearchInternalNode(t *testing.T) {
	t.Run("search with internal nodes", func(t *testing.T) {

		leaf1 := &Node{
			IsLeaf: true,
			Keys:   []string{"apple", "banana", "cat"},
			Values: []*Command{
				{Text: "apple"},
				{Text: "banana"},
				{Text: "cat"},
			},
		}

		leaf2 := &Node{
			IsLeaf: true,
			Keys:   []string{"dog", "elephant", "fish"},
			Values: []*Command{
				{Text: "dog"},
				{Text: "elephant"},
				{Text: "fish"},
			},
		}

		leaf3 := &Node{
			IsLeaf: true,
			Keys:   []string{"goat", "horse", "iguana"},
			Values: []*Command{
				{Text: "goat"},
				{Text: "horse"},
				{Text: "iguana"},
			},
		}

		leaf1.Next = leaf2
		leaf2.Next = leaf3

		root := &Node{
			IsLeaf: false,
			Keys:   []string{"dog", "goat"},
			Children: []*Node{
				leaf1, // < dog
				leaf2, // >= dog && < goat
				leaf3, // >= goat
			},
		}

		tree := &Bptree{
			Root: root,
		}
		cmd := tree.Search("goat", tree.Root)
		expected := "goat"
		if cmd.Text != expected {
			t.Errorf("didn't get expected form search for %s got %s", expected, cmd.Text)
		}
	})
}
