package selection

import (
	"github.com/femnad/barn/entity"
	"io/fs"
	"path/filepath"
	"strings"
)

func maybeBuildEntry(path string, d fs.DirEntry, settings entity.ActionSettings) *entity.Entry {
	if d.IsDir() {
		return nil
	}

	if settings.Extension != "" && !strings.HasSuffix(path, "."+settings.Extension) {
		return nil
	}

	entry := buildEntry(path, settings)
	return &entry
}

func walkdir(target string, settings entity.ActionSettings) ([]entity.Entry, error) {
	var entries []entity.Entry

	err := filepath.WalkDir(target, func(path string, d fs.DirEntry, err error) error {
		entry := maybeBuildEntry(path, d, settings)
		if entry != nil {
			entries = append(entries, *entry)
		}
		return nil
	})
	if err != nil {
		return entries, err
	}

	return entries, nil
}
