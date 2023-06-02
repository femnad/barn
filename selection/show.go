package selection

import (
	"fmt"

	"github.com/femnad/barn/entity"
)

func accumulate(selector entity.Selector) ([]entity.Entry, error) {
	var output []entity.Entry

	fn, err := getActionFn(selector.Action)
	if err != nil {
		return nil, err
	}

	for _, targetPath := range selector.Target {
		selections, fErr := fn(targetPath, selector.Settings)
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

	bucket := selector.Bucket
	if bucket == "" {
		bucket = id
	}

	storedEntries, err := getSelectionMap(cfg, bucket, selections)
	if err != nil {
		return err
	}

	sorted := sortEntries(storedEntries, true)
	for _, selection := range sorted {
		fmt.Println(selection.value.DisplayName)
	}

	return nil
}
