package selection

import (
	"fmt"
	"sort"
)

type pair struct {
	key   string
	count int64
}

func Show(configFile, id string) error {
	cfg, err := getConfig(configFile)
	if err != nil {
		return err
	}

	selector, err := getSelector(cfg, id)
	if err != nil {
		return err
	}

	fn, err := getActionFn(selector.Action)
	if err != nil {
		return err
	}

	selections, err := fn(selector.Arg)
	if err != nil {
		return err
	}

	countMap, err := getSelectionMap(cfg, id)
	if err != nil {
		return err
	}

	for _, selection := range selections {
		_, ok := countMap[selection]
		if !ok {
			countMap[selection] = 0
		}
	}

	var sorted []pair
	for k, v := range countMap {
		sorted = append(sorted, pair{key: k, count: v})
	}
	// Reverse order as that's what fzf expects by default.
	sort.Slice(sorted, func(i, j int) bool {
		itemI := sorted[i]
		itemJ := sorted[j]
		if itemI.count == itemJ.count {
			return itemI.key > itemJ.key
		}
		return itemI.count > itemJ.count
	})

	for _, selection := range sorted {
		fmt.Println(selection.key)
	}

	return nil
}
