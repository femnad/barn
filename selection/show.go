package selection

import (
	"fmt"
	"sort"

	"github.com/femnad/barn/entity"
)

type pair struct {
	key   string
	value entity.Entry
}

func accumulate(selector entity.Selector) ([]entity.Entry, error) {
	var output []entity.Entry

	fn, err := getActionFn(selector.Action)
	if err != nil {
		return nil, err
	}

	for _, arg := range selector.Args {
		selections, fErr := fn(arg)
		if fErr != nil {
			return nil, fErr
		}
		output = append(output, selections...)
	}

	return output, nil
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

	selections, err := accumulate(selector)
	if err != nil {
		return err
	}

	storedEntries, err := getSelectionMap(cfg, id, selections)
	if err != nil {
		return err
	}

	for _, selection := range selections {
		_, ok := storedEntries[selection.FullName]
		if !ok {
			storedEntries[selection.FullName] = entity.Entry{
				DisplayName: selection.DisplayName,
				FullName:    selection.FullName,
				Count:       selection.Count,
			}
		}
	}

	var sorted []pair
	for k, v := range storedEntries {
		sorted = append(sorted, pair{key: k, value: v})
	}
	// Reverse order as that's what fzf expects by default.
	sort.Slice(sorted, func(i, j int) bool {
		itemI := sorted[i]
		itemJ := sorted[j]
		if itemI.value.Count == itemJ.value.Count {
			return itemI.key > itemJ.key
		}
		return itemI.value.Count > itemJ.value.Count
	})

	for _, selection := range sorted {
		fmt.Println(selection.value.DisplayName)
	}

	return nil
}
