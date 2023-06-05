package selection

import (
	"sort"

	"github.com/femnad/barn/entity"
)

type pair struct {
	key   string
	value entity.Entry
}

func sortPairs(pairs []pair) []pair {
	sort.Slice(pairs, func(i, j int) bool {
		itemI := pairs[i]
		itemJ := pairs[j]
		if itemI.value.Count == itemJ.value.Count {
			// Ascending order by key if counts are equal.
			return itemI.key < itemJ.key
		}

		// Descending order by count if counts are different.
		return itemI.value.Count > itemJ.value.Count
	})

	return pairs
}

func sortEntries(entries selectionMap) []entity.Entry {
	var nonZeroCounts []pair
	var zeroCounts []pair
	var merged []pair

	for k, v := range entries {
		p := pair{key: k, value: v}
		if v.Count > 0 {
			nonZeroCounts = append(nonZeroCounts, p)
		} else {
			zeroCounts = append(zeroCounts, p)
		}
	}

	nonZeroCounts = sortPairs(nonZeroCounts)
	zeroCounts = sortPairs(zeroCounts)
	merged = append(nonZeroCounts, zeroCounts...)

	var items []entity.Entry
	for _, p := range merged {
		items = append(items, p.value)
	}

	return items
}
