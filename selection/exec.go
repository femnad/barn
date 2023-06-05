package selection

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/femnad/barn/entity"
	marecmd "github.com/femnad/mare/cmd"
)

func execCmd(target string, settings entity.ActionSettings) ([]entity.Entry, error) {
	var entries []entity.Entry
	if target == "" {
		return entries, fmt.Errorf("given command is empty")
	}

	cmdSlice := strings.Split(target, " ")
	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)

	if settings.SetPwdCmd != "" {
		out, err := marecmd.RunFormatError(marecmd.Input{Command: settings.SetPwdCmd})
		if err != nil {
			return nil, err
		}
		cmd.Dir = strings.TrimSpace(out.Stdout)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return entries, err
	}

	if err = cmd.Start(); err != nil {
		return entries, fmt.Errorf("error running command %s: %v", target, err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		entry := buildEntry(line, settings)
		entries = append(entries, entry)
	}

	if err = cmd.Wait(); err != nil {
		return entries, err
	}

	return entries, nil
}
