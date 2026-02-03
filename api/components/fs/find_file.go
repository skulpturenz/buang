package fs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"skulpture/buang/app"
)

type FindFileParams struct {
	File string
	Dir  string
}

type FindFileResult struct {
	Path string
}

func (f FindFileParams) FindFile(ctx context.Context, s *app.ApplicationServices) (*FindFileResult, error) {
	p := ""
	err := filepath.WalkDir(f.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == f.File {
			p = path

			return errors.New("success")
		}

		return nil
	})

	if err != nil {
		if err.Error() == "success" {
			return &FindFileResult{
				Path: p,
			}, nil
		}

		return nil, err
	}

	return nil, fmt.Errorf("file %s not found", f.File)
}
