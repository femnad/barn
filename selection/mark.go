package selection

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/femnad/mare"
)

func Mark(configFile, id, selection string) error {
	cfg, err := getConfig(configFile)
	if err != nil {
		return err
	}

	selector, err := getSelector(cfg, id)
	if err != nil {
		return err
	}

	err = storeSelection(cfg, id, selection)
	if err != nil {
		return err
	}

	ch := choice{Selector: selector, Selection: selection}

	tmpl, err := template.New("on-select").Parse(ch.OnSelect)
	if err != nil {
		return err
	}

	out := bytes.Buffer{}
	err = tmpl.Execute(&out, ch)
	if err != nil {
		return err
	}
	outStr := mare.ExpandUser(out.String())

	fmt.Printf("%s\n", outStr)
	return nil
}
