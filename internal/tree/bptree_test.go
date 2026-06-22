package tree

import (
	"fmt"
	"math/rand"
	"testing"
)

const MAX = 4
const MIN = 2

func TestRootAsLeaf(t *testing.T) {
	t.Run("check if basic insert is working for 1 element", func(t *testing.T) {
		tree := NewTree(MAX, MIN)

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
		tree := NewTree(MAX, MIN)

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
		tree := NewTree(MAX, MIN)

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

func TestInsert(t *testing.T) {
	t.Run("Insert 1 element", func(t *testing.T) {
		tree := NewTree(MAX, MIN)

		tree.Insert("ls")

		lenOfKeys := len(tree.Root.Keys)
		lenOfVals := len(tree.Root.Values)

		node := tree.Root.Values[0]

		if lenOfKeys != 1 {
			t.Errorf("Invalid KeyLength %d expected 1", lenOfKeys)
		}
		if lenOfVals != 1 {
			t.Errorf("Invalid ValLength %d expected 1", lenOfVals)
		}
		if node.Text != "ls" {
			t.Errorf("Invalid Text %s expected ls", node.Text)
		}
	})
	t.Run("Insert multiple elements", func(t *testing.T) {
		tree := NewTree(MAX, MIN)

		commands := []string{"ls", "cd", "mkdir", "rm", "touch"}

		for _, cmd := range commands {
			tree.Insert(cmd)
		}

		for _, cmd := range commands {
			result := tree.Search(cmd, tree.Root)
			if result == nil || result.Text != cmd {
				t.Errorf("Search failed for %s, got %v", cmd, result)
			}
		}
	})
	t.Run("Insert duplicate elements", func(t *testing.T) {
		tree := NewTree(MAX, MIN)

		tree.Insert("ls")
		tree.Insert("ls") // Duplicate insert
		cmd := tree.Search("ls", tree.Root)
		if cmd.Frequency != 2 {
			t.Errorf("Expected frequency 2, got %d", cmd.Frequency)
		}
	})
	t.Run("random inserts and invariant checks", func(t *testing.T) {
		tree := NewTree(MAX, MIN)
		inserted := make(map[string]bool)

		// Insert 100 random keys
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("cmd-%d", rand.Intn(1000))
			tree.Insert(key)
			inserted[key] = true
		}

		// 1. Check search correctness
		for key := range inserted {
			res := tree.Search(key, tree.Root)
			if res == nil || res.Text != key {
				t.Errorf("Expected to find key %s", key)
			}
		}
	})

}

func TestBptree_SearchAfterInternalSplit(t *testing.T) {
	// Construct a tree with MaxSize = 2
	tr := NewTree(2, 1)

	// Insert elements sequentially to force an internal node split
	tr.Insert("A")
	tr.Insert("B")
	tr.Insert("C")
	tr.Insert("D")
	tr.Insert("E")

	// Validate structural integrity
	validate(tr.Root)

	// Try to find "C"
	cmd := tr.Search("C", tr.Root)
	if cmd == nil {
		t.Fatalf("FAIL: Command 'C' was not found in the tree, but it should be present!")
	} else {
		t.Logf("SUCCESS: Found command %s", cmd.Text)
	}
}

func TestBptree_StructureAndSearch(t *testing.T) {
	// 1. Initialize tree with MaxSize = 2
	tr := NewTree(2, 1)

	// 2. Insert sequential elements to force a split up to the root level
	commands := []string{"A", "B", "C", "D", "E"}
	for _, cmd := range commands {
		tr.Insert(cmd)
	}

	// 3. Verify tree structure properties via your validate function
	validate(tr.Root)

	// 4. Test searching for existing items
	for _, cmd := range commands {
		found := tr.Search(cmd, tr.Root)
		if found == nil {
			t.Errorf("Expected to find key %q, but got nil", cmd)
		} else if found.Text != cmd {
			t.Errorf("Expected key %q, but retrieved %q", cmd, found.Text)
		}
	}
}
