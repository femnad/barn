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
	settings := selector.Settings

	if settings.OnSelect == "" {
		return outStr, nil
	}

	ch := choice{Entry: entry, Selector: selector, Selection: selection}

	tmpl, err := template.New("on-select").Parse(settings.OnSelect)
	if err != nil {
		return outStr, err
	}

	out := bytes.Buffer{}
	err = tmpl.Execute(&out, ch)
	if err != nil {
		return "", err
	}
	outStr = mare.ExpandUser(out.String())

	if settings.ExecOnSelect {
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
	exitCode = selector.Settings.ExitOnSelect

	bucket, err := getBucket(id, selector)
	if err != nil {
		return exitCode, err
	}

	entry, err := incrementEntryCount(cfg, bucket, selection, selector.Settings.Eager)
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
	if selector.Settings.StderrOutput {
		outStream = os.Stderr
	}

	_, err = fmt.Fprintf(outStream, "%s\n", outStr)
	if err != nil {
		return exitCode, err
	}

	return exitCode, nil
}
