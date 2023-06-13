package selection

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/anmitsu/go-shlex"

	"github.com/femnad/barn/entity"
	marecmd "github.com/femnad/mare/cmd"
)

func execCmd(target string, settings entity.ActionSettings) ([]entity.Entry, error) {
	var entries []entity.Entry
	if target == "" {
		return entries, fmt.Errorf("given command is empty")
	}

	cmdSlice, err := shlex.Split(target, true)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)

	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr

	if settings.SetPwdCmd != "" {
		out, sErr := marecmd.RunFormatError(marecmd.Input{Command: settings.SetPwdCmd})
		if sErr != nil {
			return nil, sErr
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
		stderrContent, rErr := io.ReadAll(&stderr)
		if rErr != nil {
			return nil, fmt.Errorf("error reading stderr of command with error exit %s: %v", target, rErr)
		}
		return nil, fmt.Errorf("error running command %s, stderr: %s, error: %v", target, stderrContent, rErr)
	}

	return entries, nil
}
