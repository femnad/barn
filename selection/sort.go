package selection

import (
	"sort"

	"golang.org/x/exp/constraints"

	"github.com/femnad/barn/entity"
)

type pair struct {
	key   string
	value entity.Entry
}

func less[T constraints.Ordered](i, j T, reverse bool) bool {
	if reverse {
		return i > j
	}
	return i < j
}

func sortEntries(entries selectionMap, reverse bool) []pair {
	var sorted []pair
	for k, v := range entries {
		sorted = append(sorted, pair{key: k, value: v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		itemI := sorted[i]
		itemJ := sorted[j]
		if itemI.value.Count == itemJ.value.Count {
			return less[string](itemJ.key, itemJ.key, reverse)
		}
		return less[int64](itemI.value.Count, itemJ.value.Count, reverse)
	})

	return sorted
}
