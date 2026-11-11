package storage

import (
	"dash/internal/tree"
	"testing"
)

func TestStorageRoundTrip(t *testing.T) {
	tr := tree.NewTree(4, 2)
	tr.Insert("git status")
	tr.Insert("git commit -m 'init'")
	tr.Insert("git push origin main")
	tr.Insert("docker ps")
	tr.Insert("docker build .")

	diskTree := Generate(tr)
	data, err := SerializeDisktree(diskTree)
	if err != nil {
		t.Fatalf("SerializeDisktree failed: %v", err)
	}

	deserialized := DeserializeDisktree(data)
	if deserialized.RootID != diskTree.RootID {
		t.Errorf("expected rootID %d, got %d", diskTree.RootID, deserialized.RootID)
	}
	if len(deserialized.Nodes) != len(diskTree.Nodes) {
		t.Errorf("expected %d nodes, got %d", len(diskTree.Nodes), len(deserialized.Nodes))
	}
}
