package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	shellexpand "github.com/ganbarodigital/go_shellexpand"
)

func ExpandComposeYaml(composePath string, vars map[string]string) (*string, bool, error) {
	if !strings.Contains(strings.ToLower(composePath), ".buang.yml") && !strings.Contains(strings.ToLower(composePath), ".buang.yaml") {
		return &composePath, false, nil
	}

	cb := shellexpand.ExpansionCallbacks{
		LookupVar: func(s string) (string, bool) {
			v, ok := vars[s]

			return v, ok
		},
	}
	dir := filepath.Dir(composePath)
	expandedFileName := fmt.Sprintf("expanded-%v", filepath.Base(composePath))
	expandedFilePath := filepath.Join(dir, expandedFileName)

	b, err := os.ReadFile(composePath)
	if err != nil {
		return nil, false, err
	}

	contents := string(b)

	expanded, err := shellexpand.Expand(contents, cb)
	if err != nil {
		return nil, false, err
	}

	expandedFile, err := os.Create(expandedFilePath)
	if err != nil {
		return nil, false, err
	}

	if _, err := fmt.Fprintf(expandedFile, "%v", expanded); err != nil {
		return nil, false, err
	}

	if err := expandedFile.Sync(); err != nil {
		return nil, false, err
	}

	return &expandedFilePath, true, nil
}
