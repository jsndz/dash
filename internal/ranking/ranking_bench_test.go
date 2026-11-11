package ranking

import (
	"dash/internal/tree"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkGetRanking(b *testing.B) {
	for _, count := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("Candidates=%d", count), func(b *testing.B) {
			r := rand.New(rand.NewSource(42))
			baseList := make([]*tree.Command, count)
			for i := 0; i < count; i++ {
				baseList[i] = &tree.Command{
					Text:      fmt.Sprintf("cmd-%05d", i),
					Frequency: r.Intn(1000),
					LastUsed:  time.Now(),
				}
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Copy slice headers to benchmark realistic sorting
				cmds := make([]*tree.Command, count)
				copy(cmds, baseList)
				_ = GetRanking(cmds)
			}
		})
	}
}
