package selection

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/femnad/barn/entity"
	"github.com/femnad/mare"
	"github.com/femnad/mare/cmd"
)

const gitTopLevelCmd = "git rev-parse --show-toplevel"

type choice struct {
	entity.Entry
	entity.Selector
	Selection string
}

type env struct {
	bucket string
}

func (e env) GitRoot() (string, error) {
	gitOut, err := cmd.RunFormatError(cmd.Input{Command: gitTopLevelCmd})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(gitOut.Stdout), nil
}

func (e env) Id() (string, error) {
	return e.bucket, nil
}

func (e env) Pwd() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return pwd, nil
}

func expandBucketTemplate(selector entity.Selector) (string, error) {
	bucket := os.ExpandEnv(selector.Bucket)

	// Primitive check to see if the bucket name is templated so that we can return early.
	if !strings.Contains(bucket, "{{") {
		return bucket, nil
	}

	tmpl, err := template.New("bucket").Parse(bucket)
	if err != nil {
		return "", err
	}

	out := bytes.Buffer{}
	e := env{bucket: selector.Id}

	err = tmpl.Execute(&out, e)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func getBucket(id string, selector entity.Selector) (string, error) {
	bucket := selector.Bucket
	if bucket == "" {
		return id, nil
	}

	return expandBucketTemplate(selector)
}

func getDisplayName(entryName, targetPath string, includeParents int) string {
	if includeParents <= 0 {
		return entryName
	}

	parents := strings.Split(targetPath, "/")
	parents = mare.Filter[string](parents, func(d string) bool {
		return d != ""
	})
	numParents := len(parents)
	startIndex := numParents - includeParents - 1
	if startIndex < 0 {
		return targetPath
	}

	prefix := parents[numParents-includeParents-1 : numParents]

	name := strings.Join(prefix, "/")
	if includeParents >= numParents-1 {
		name = "/" + name
	}
	return name
}

func buildEntry(line string, settings entity.ActionSettings) entity.Entry {
	displayName := line

	if settings.RemovePrefix != "" {
		prefix := os.ExpandEnv(settings.RemovePrefix)
		displayName = strings.TrimPrefix(displayName, prefix)
	}

	if settings.RemoveSuffix != "" {
		displayName = strings.TrimSuffix(displayName, settings.RemoveSuffix)
	}

	return entity.Entry{FullName: line, DisplayName: displayName}
}

func readdir(target string, settings entity.ActionSettings) ([]entity.Entry, error) {
	var out []entity.Entry
	target = mare.ExpandUser(target)
	entries, err := os.ReadDir(target)
	if err != nil {
		return out, fmt.Errorf("error reading contents of directory %s: %v", target, err)
	}

	for _, i := range entries {
		name := i.Name()
		fullPath := path.Join(target, name)
		displayName := getDisplayName(name, fullPath, settings.IncludeParents)
		e := entity.Entry{DisplayName: displayName, FullName: fullPath}
		out = append(out, e)
	}

	return out, nil
}

func getActionFn(action string) (func(string, entity.ActionSettings) ([]entity.Entry, error), error) {
	switch action {
	case "exec":
		return execCmd, nil
	case "readdir":
		return readdir, nil
	case "walkdir":
		return walkdir, nil
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
		return cfg, fmt.Errorf("error unmarshalling config: %v", err)
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
