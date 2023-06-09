package selection

import (
	"fmt"
	"os"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/femnad/barn/entity"
	"github.com/femnad/mare"
)

func accumulate(selector entity.Selector, extraArgs string) ([]entity.Entry, error) {
	var output []entity.Entry

	fn, err := getActionFn(selector.Action)
	if err != nil {
		return nil, err
	}

	for _, target := range selector.Target {
		target = mare.ExpandUser(os.ExpandEnv(target))
		if extraArgs != "" {
			target = fmt.Sprintf("%s %s", target, extraArgs)
		}
		selections, fErr := fn(target, selector.Settings)
		if fErr != nil {
			return nil, fErr
		}
		output = append(output, selections...)
	}

	return output, nil
}

func getSelections(cfg entity.Config, bucket string, selections []entity.Entry, eager bool) (selectionMap, error) {
	if eager {
		return getSelectionMap(cfg, bucket, selections)
	}

	var selectionList []string
	for _, selection := range selections {
		selectionList = append(selectionList, selection.DisplayName)
	}
	validSelections := mapset.NewSet[string](selectionList...)

	storedSelections, err := getLazySelectionMap(cfg, bucket, validSelections)
	if err != nil {
		return nil, err
	}

	for _, entry := range selections {
		// Check if there is persisted entry which could have a non-zero count.
		_, ok := storedSelections[entry.DisplayName]
		if !ok {
			storedSelections[entry.DisplayName] = entry
		}
	}

	return storedSelections, nil
}

func Show(configFile, id, extraArgs string) error {
	cfg, err := getConfig(configFile)
	if err != nil {
		return err
	}

	selector, err := getSelector(cfg, id)
	if err != nil {
		return err
	}

	selections, err := accumulate(selector, extraArgs)
	if err != nil {
		return err
	}

	bucket, err := getBucket(id, extraArgs, selector)
	if err != nil {
		return err
	}

	storedEntries, err := getSelections(cfg, bucket, selections, selector.Settings.Eager)
	if err != nil {
		return err
	}

	sorted := sortEntries(storedEntries)
	for _, selection := range sorted {
		fmt.Println(selection.DisplayName)
	}

	return nil
}
