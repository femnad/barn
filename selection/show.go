package selection

import (
	"fmt"
	"os"

	"github.com/femnad/barn/entity"
)

func accumulate(selector entity.Selector) ([]entity.Entry, error) {
	var output []entity.Entry

	fn, err := getActionFn(selector.Action)
	if err != nil {
		return nil, err
	}

	for _, target := range selector.Target {
		target = os.ExpandEnv(target)
		selections, fErr := fn(target, selector.Settings)
		if fErr != nil {
			return nil, fErr
		}
		output = append(output, selections...)
	}

	return output, nil
}

func Show(configFile, id string, reverse bool) error {
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

	bucket, err := getBucket(id, selector)
	if err != nil {
		return err
	}

	storedEntries, err := getSelectionMap(cfg, bucket, selections)
	if err != nil {
		return err
	}

	sorted := sortEntries(storedEntries, reverse)
	for _, selection := range sorted {
		fmt.Println(selection.value.DisplayName)
	}

	return nil
}
