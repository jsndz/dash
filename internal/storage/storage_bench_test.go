package storage

import (
	"dash/internal/tree"
	"fmt"
	"testing"
)

func buildTree(n int) *tree.Bptree {
	tr := tree.NewTree(128, 64)
	for i := 0; i < n; i++ {
		tr.Insert(fmt.Sprintf("git commit -m 'commit message %06d'", i))
	}
	return tr
}

func BenchmarkGenerateDisktree(b *testing.B) {
	for _, n := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("Commands=%d", n), func(b *testing.B) {
			tr := buildTree(n)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Generate(tr)
			}
		})
	}
}

func BenchmarkSerializeDisktree(b *testing.B) {
	for _, n := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("Commands=%d", n), func(b *testing.B) {
			tr := buildTree(n)
			dt := Generate(tr)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := SerializeDisktree(dt)
				if err != nil {
					b.Fatalf("serialization failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkDeserializeDisktree(b *testing.B) {
	for _, n := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("Commands=%d", n), func(b *testing.B) {
			tr := buildTree(n)
			dt := Generate(tr)
			data, err := SerializeDisktree(dt)
			if err != nil {
				b.Fatalf("serialization failed: %v", err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(data)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = DeserializeDisktree(data)
			}
		})
	}
}
