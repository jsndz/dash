package tree

import (
	"fmt"
	"math/rand"
	"testing"
)

func generateKeys(n int) []string {
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = fmt.Sprintf("command-path-%08d", i)
	}
	return keys
}

func generateRandomKeys(n int, seed int64) []string {
	r := rand.New(rand.NewSource(seed))
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = fmt.Sprintf("cmd-%08d", r.Intn(n*10))
	}
	return keys
}

func BenchmarkTreeInsertSequential(b *testing.B) {
	for _, n := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			keys := generateKeys(n)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tr := NewTree(128, 64)
				for j := 0; j < n; j++ {
					tr.Insert(keys[j])
				}
			}
		})
	}
}

func BenchmarkTreeInsertRandom(b *testing.B) {
	for _, n := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			keys := generateRandomKeys(n, 42)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tr := NewTree(128, 64)
				for j := 0; j < n; j++ {
					tr.Insert(keys[j])
				}
			}
		})
	}
}

func BenchmarkTreeSearchHit(b *testing.B) {
	for _, n := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			tr := NewTree(128, 64)
			keys := generateKeys(n)
			for _, k := range keys {
				tr.Insert(k)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				target := keys[i%n]
				cmd := tr.Search(target, tr.Root)
				if cmd == nil {
					b.Fatalf("expected key %s to be found", target)
				}
			}
		})
	}
}

func BenchmarkTreeSearchMiss(b *testing.B) {
	for _, n := range []int{1000, 10000} {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			tr := NewTree(128, 64)
			keys := generateKeys(n)
			for _, k := range keys {
				tr.Insert(k)
			}
			missKeys := make([]string, 100)
			for i := 0; i < 100; i++ {
				missKeys[i] = fmt.Sprintf("nonexistent-key-%08d", i)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				target := missKeys[i%100]
				cmd := tr.Search(target, tr.Root)
				if cmd != nil {
					b.Fatalf("expected key %s to not be found", target)
				}
			}
		})
	}
}

func BenchmarkTreeRangeScan(b *testing.B) {
	tr := NewTree(128, 64)
	// Insert 500 git commands, 500 docker commands, and 5000 general commands
	for i := 0; i < 500; i++ {
		tr.Insert(fmt.Sprintf("git checkout branch-%04d", i))
		tr.Insert(fmt.Sprintf("docker container run-%04d", i))
	}
	for i := 0; i < 5000; i++ {
		tr.Insert(fmt.Sprintf("z-misc-cmd-%05d", i))
	}

	b.Run("PrefixSmallMatch_5", func(b *testing.B) {
		for i := 0; i < 5; i++ {
			tr.Insert(fmt.Sprintf("npm run test-%d", i))
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res := tr.RangeScan("npm run")
			if len(res) == 0 {
				b.Fatal("expected results")
			}
		}
	})

	b.Run("PrefixMediumMatch_500", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res := tr.RangeScan("git checkout")
			if len(res) != 500 {
				b.Fatalf("expected 500 results, got %d", len(res))
			}
		}
	})

	b.Run("PrefixNoMatch", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res := tr.RangeScan("kubectl")
			if len(res) != 0 {
				b.Fatalf("expected 0 results, got %d", len(res))
			}
		}
	})
}

func BenchmarkTreeDelete(b *testing.B) {
	n := 1000
	keys := generateKeys(n)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tr := NewTree(128, 64)
		for _, k := range keys {
			tr.Insert(k)
		}
		b.StartTimer()

		for _, k := range keys {
			_ = tr.Delete(k)
		}
	}
}
