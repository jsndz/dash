package directory

import (
	"os"
	"path/filepath"
	"strings"
)

func ScanDirectory(path string) ([]string, error) {
	dir := filepath.Dir(path)
	prefix := filepath.Base(path)

	if !filepath.IsAbs(dir) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(wd, dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), prefix) {
			result = append(result, entry.Name())
		}
	}

	return result, nil
}
