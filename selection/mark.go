package selection

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/femnad/mare"
	"github.com/femnad/mare/cmd"
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

	bucket := selector.Bucket
	if bucket == "" {
		bucket = id
	}

	entry, err := incrementEntryCount(cfg, bucket, selection)
	if err != nil {
		return err
	}

	ch := choice{Entry: entry, Selector: selector, Selection: selection}

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

	if selector.ExecOnSelect {
		cmdOut, cErr := cmd.RunFormatError(cmd.Input{Command: outStr})
		if cErr != nil {
			return cErr
		}
		outStr = strings.TrimSpace(cmdOut.Stdout)
	}

	fmt.Printf("%s\n", outStr)
	return nil
}
