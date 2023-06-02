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

func sortEntries(entries selectionMap, reverse bool) []entity.Entry {
	var pairs []pair
	for k, v := range entries {
		pairs = append(pairs, pair{key: k, value: v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		itemI := pairs[i]
		itemJ := pairs[j]
		if itemI.value.Count == itemJ.value.Count {
			return less[string](itemI.key, itemJ.key, reverse)
		}
		return less[int64](itemI.value.Count, itemJ.value.Count, reverse)
	})

	var items []entity.Entry
	for _, p := range pairs {
		items = append(items, p.value)
	}

	return items
}
