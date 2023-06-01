package selection

import (
	"fmt"
	"github.com/femnad/barn/config"
	"github.com/femnad/mare"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type choice struct {
	config.Selector
	Selection string
}

func readdir(arg string) ([]string, error) {
	var out []string
	arg = mare.ExpandUser(arg)
	entries, err := os.ReadDir(arg)
	if err != nil {
		return out, err
	}

	for _, entry := range entries {
		out = append(out, entry.Name())
	}

	return out, nil
}

func getActionFn(action string) (func(string) ([]string, error), error) {
	switch action {
	case "readdir":
		return readdir, nil
	default:
		return nil, fmt.Errorf("no function found for %s", action)
	}
}

func getConfig(file string) (config.Config, error) {
	var cfg config.Config

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

func getSelector(cfg config.Config, id string) (config.Selector, error) {
	for _, selector := range cfg.Selectors {
		if selector.Id != id {
			continue
		}

		return selector, nil
	}

	return config.Selector{}, fmt.Errorf("no selector defined for id %s", id)
}
