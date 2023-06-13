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

const defaultShell = "sh"

func getCmd(command string, shell bool) (*exec.Cmd, error) {
	if shell {
		return exec.Command(defaultShell, "-c", command), nil
	}

	cmdSlice, err := shlex.Split(command, true)
	if err != nil {
		return nil, err
	}

	return exec.Command(cmdSlice[0], cmdSlice[1:]...), nil
}

func runCmd(target string, settings entity.ActionSettings, shell bool) ([]entity.Entry, error) {
	var entries []entity.Entry
	if target == "" {
		return entries, fmt.Errorf("given command is empty")
	}

	cmd, err := getCmd(target, shell)
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("error running command %s, stderr: %s, error: %v", target, stderrContent, err)
	}

	return entries, nil
}

func execCmd(target string, settings entity.ActionSettings) ([]entity.Entry, error) {
	return runCmd(target, settings, false)
}

func shellCmd(target string, settings entity.ActionSettings) ([]entity.Entry, error) {
	return runCmd(target, settings, true)
}
