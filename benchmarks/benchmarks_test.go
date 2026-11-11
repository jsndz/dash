package benchmarks_test

import (
	"dash/internal/search"
	"dash/internal/storage"
	"dash/internal/tree"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
)

// generateRealisticHistory creates a simulated shell history with command clusters
// (e.g. repeated git commands, docker commands, navigation, scripts).
func generateRealisticHistory(total int) []string {
	commandTemplates := []string{
		"git status",
		"git add .",
		"git commit -m 'wip'",
		"git checkout -b feature/%d",
		"git push origin feature/%d",
		"git pull --rebase origin main",
		"git log -n %d --oneline",
		"git diff HEAD~%d",
		"docker ps -a",
		"docker build -t app:%d .",
		"docker compose up -d",
		"docker compose logs -f service-%d",
		"docker stop container-%d",
		"npm run build",
		"npm test -- -u",
		"npm install pkg-%d",
		"cargo check",
		"cargo test test_%d",
		"cargo run --release",
		"cd /var/log/%d",
		"ls -la /home/user/%d",
		"curl -X GET http://localhost:8080/api/v1/%d",
		"cat /etc/hosts",
		"ssh server-%d.prod.internal",
		"kubectl get pods -n ns-%d",
		"kubectl logs pod-%d",
		"kubectl describe pod-%d",
	}

	r := rand.New(rand.NewSource(999))
	commands := make([]string, total)
	for i := 0; i < total; i++ {
		tmpl := commandTemplates[r.Intn(len(commandTemplates))]
		if strings.Contains(tmpl, "%d") {
			commands[i] = fmt.Sprintf(tmpl, r.Intn(200))
		} else {
			commands[i] = tmpl
		}
	}
	return commands
}

// linearScanAutocomplete simulates the naive shell history linear search and ranking baseline.
type historyEntry struct {
	text      string
	frequency int
	lastUsed  time.Time
}

func linearScanAutocomplete(history []historyEntry, prefix string) []string {
	var matches []*historyEntry
	for i := range history {
		if strings.HasPrefix(history[i].text, prefix) {
			matches = append(matches, &history[i])
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].frequency == matches[j].frequency {
			return matches[i].text < matches[j].text
		}
		return matches[i].frequency > matches[j].frequency
	})
	results := make([]string, len(matches))
	for i, m := range matches {
		results[i] = m.text
	}
	return results
}

// BenchmarkAutocomplete_BptreeVsLinearScan benchmarks Dash's B+ Tree against
// a naive linear scan over history of sizes 1K, 10K, and 50K commands.
func BenchmarkAutocomplete_BptreeVsLinearScan(b *testing.B) {
	for _, size := range []int{1000, 10000, 50000} {
		historyRaw := generateRealisticHistory(size)

		// 1. Build B+ Tree
		tr := tree.NewTree(128, 64)
		for _, cmd := range historyRaw {
			tr.Insert(cmd)
		}
		engine := search.NewSearchEngine(tr)

		// 2. Build Linear Scan entries (aggregated by unique command)
		freqMap := make(map[string]int)
		for _, cmd := range historyRaw {
			freqMap[cmd]++
		}
		var linearEntries []historyEntry
		for cmd, freq := range freqMap {
			linearEntries = append(linearEntries, historyEntry{
				text:      cmd,
				frequency: freq,
				lastUsed:  time.Now(),
			})
		}

		query := "docker compose"

		b.Run(fmt.Sprintf("Bptree/HistorySize=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				res := engine.Autocomplete(query)
				if len(res) == 0 {
					b.Fatal("expected results")
				}
			}
		})

		b.Run(fmt.Sprintf("LinearScan/HistorySize=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				res := linearScanAutocomplete(linearEntries, query)
				if len(res) == 0 {
					b.Fatal("expected results")
				}
			}
		})
	}
}

// BenchmarkEndToEndIndexing benchmarks the cost to import and index command batches.
func BenchmarkEndToEndIndexing(b *testing.B) {
	for _, size := range []int{1000, 5000, 10000} {
		raw := generateRealisticHistory(size)
		b.Run(fmt.Sprintf("Commands=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tr := tree.NewTree(128, 64)
				for _, cmd := range raw {
					tr.Insert(cmd)
				}
			}
		})
	}
}

// BenchmarkEndToEndStoragePipeline benchmarks generating, serializing, and deserializing the index.
func BenchmarkEndToEndStoragePipeline(b *testing.B) {
	for _, size := range []int{1000, 5000} {
		raw := generateRealisticHistory(size)
		tr := tree.NewTree(128, 64)
		for _, cmd := range raw {
			tr.Insert(cmd)
		}

		b.Run(fmt.Sprintf("FullSaveAndLoad/Commands=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// 1. Generate DiskTree
				dt := storage.Generate(tr)
				// 2. Serialize
				data, err := storage.SerializeDisktree(dt)
				if err != nil {
					b.Fatal(err)
				}
				// 3. Deserialize
				_ = storage.DeserializeDisktree(data)
			}
		})
	}
}
