package selection

import (
	"fmt"
	"io"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/femnad/barn/entity"
	"github.com/femnad/mare"
)

type choice struct {
	entity.Entry
	entity.Selector
	Selection string
}

func readdir(arg string) ([]entity.Entry, error) {
	var out []entity.Entry
	arg = mare.ExpandUser(arg)
	entries, err := os.ReadDir(arg)
	if err != nil {
		return out, err
	}

	for _, i := range entries {
		name := i.Name()
		fullPath := path.Join(arg, name)
		e := entity.Entry{DisplayName: name, FullName: fullPath}
		out = append(out, e)
	}

	return out, nil
}

func getActionFn(action string) (func(string) ([]entity.Entry, error), error) {
	switch action {
	case "readdir":
		return readdir, nil
	default:
		return nil, fmt.Errorf("no function found for %s", action)
	}
}

func getConfig(file string) (entity.Config, error) {
	var cfg entity.Config

	file = mare.ExpandUser(file)
	f, err := os.Open(file)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func getSelector(cfg entity.Config, id string) (entity.Selector, error) {
	for _, selector := range cfg.Selectors {
		if selector.Id != id {
			continue
		}

		return selector, nil
	}

	return entity.Selector{}, fmt.Errorf("no selector defined for id %s", id)
}
