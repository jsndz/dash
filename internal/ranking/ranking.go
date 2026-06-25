package ranking

import (
	"dash/internal/tree"
	"sort"
)

func GetRanking(commands []*tree.Command) []string {
	sort.Slice(commands, func(i, j int) bool {
		if commands[i].Frequency == commands[j].Frequency {
			return commands[i].Text < commands[j].Text
		}
		return commands[i].Frequency > commands[j].Frequency
	})
	var results []string
	for _, cmd := range commands {
		results = append(results, cmd.Text)
	}
	return results
}
