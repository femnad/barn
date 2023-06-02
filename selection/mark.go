package selection

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/femnad/barn/entity"
	"github.com/femnad/mare"
	"github.com/femnad/mare/cmd"
)

func onSelectOutput(entry entity.Entry, selector entity.Selector, selection string) (string, error) {
	var outStr string

	if selector.OnSelect == "" {
		return outStr, nil
	}

	ch := choice{Entry: entry, Selector: selector, Selection: selection}

	tmpl, err := template.New("on-select").Parse(ch.OnSelect)
	if err != nil {
		return outStr, err
	}

	out := bytes.Buffer{}
	err = tmpl.Execute(&out, ch)
	if err != nil {
		return "", err
	}
	outStr = mare.ExpandUser(out.String())

	if selector.ExecOnSelect {
		cmdOut, cErr := cmd.RunFormatError(cmd.Input{Command: outStr})
		if cErr != nil {
			return outStr, cErr
		}
		outStr = strings.TrimSpace(cmdOut.Stdout)
	}

	return outStr, nil
}

func Mark(configFile, id, selection string) (int, error) {
	var exitCode int
	cfg, err := getConfig(configFile)
	if err != nil {
		return exitCode, err
	}

	selector, err := getSelector(cfg, id)
	if err != nil {
		return exitCode, err
	}

	bucket := selector.Bucket
	if bucket == "" {
		bucket = id
	}
	exitCode = selector.ExitOnSelect

	entry, err := incrementEntryCount(cfg, bucket, selection)
	if err != nil {
		return exitCode, err
	}

	outStr, err := onSelectOutput(entry, selector, selection)
	if err != nil {
		return exitCode, err
	}

	if outStr == "" {
		return exitCode, nil
	}

	outStream := os.Stdout
	if selector.StderrOutput {
		outStream = os.Stderr
	}

	_, err = fmt.Fprintf(outStream, "%s\n", outStr)
	if err != nil {
		return exitCode, err
	}

	return selector.ExitOnSelect, nil
}
