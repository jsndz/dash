package directory

import (
	"os"
	"path/filepath"
	"testing"
)

func touch(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, nil, 0o644)
}

func TestScanDir(t *testing.T) {
	t.Run("test ", func(t *testing.T) {
		root := t.TempDir()
		touch(filepath.Join(root, "cmd", "internal", "main.go"))
		touch(filepath.Join(root, "cmd", "internal", "man.go"))
		//   /cmd/internal/m
		// -> /cmd/internal/main.go
		// -> /cmd/internal/man.go

		expected := []string{"main.go", "man.go"}
		actual, _ := ScanDirectory(filepath.Join(root, "cmd", "internal", "ma"))
		if len(actual) != len(expected) {
			t.Fatalf("expected %v, got %v", expected, actual)
		}

	})
}
