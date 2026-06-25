package history

import (
	"dash/internal/tree"
	"os"
	"path/filepath"
	"strings"
)

func Import(tree *tree.Bptree) error {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(home, ".bash_history"))

	if err != nil {
		return err
	}
	commands := strings.Split(string(data), "\n")

	for _, command := range commands {
		// fmt.Println(command)
		tree.Insert(command)
	}
	return nil
}
