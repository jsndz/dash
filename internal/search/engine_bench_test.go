package search

import (
	"dash/internal/tree"
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkAutocomplete(b *testing.B) {
	tr := tree.NewTree(128, 64)
	prefixes := []string{"git checkout", "git commit", "docker run", "docker ps", "npm start", "npm test"}

	// Insert realistic commands with varying frequencies
	r := rand.New(rand.NewSource(123))
	for _, p := range prefixes {
		for i := 0; i < 100; i++ {
			cmd := fmt.Sprintf("%s --option-%d", p, i)
			freq := r.Intn(50) + 1
			for f := 0; f < freq; f++ {
				tr.Insert(cmd)
			}
		}
	}

	engine := NewSearchEngine(tr)

	b.Run("Prefix_git", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res := engine.Autocomplete("git")
			if len(res) == 0 {
				b.Fatal("expected results")
			}
		}
	})

	b.Run("Prefix_docker", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res := engine.Autocomplete("docker")
			if len(res) == 0 {
				b.Fatal("expected results")
			}
		}
	})

	b.Run("Prefix_npm_start", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res := engine.Autocomplete("npm start")
			if len(res) == 0 {
				b.Fatal("expected results")
			}
		}
	})

	b.Run("Prefix_nomatch", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res := engine.Autocomplete("kubectl")
			if len(res) != 0 {
				b.Fatal("expected no results")
			}
		}
	})
}
