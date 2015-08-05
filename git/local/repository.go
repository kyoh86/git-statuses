package local

import (
	"os"
	"path/filepath"
)

func WalkOnRepositories(rootPath string, walker func(repositoryPath string) error) error {
	return filepath.Walk(rootPath, func(children string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() || info.Name() != ".git" {
			return nil
		}

		repositoryPath := filepath.Dir(children)
		if err := walker(repositoryPath); err != nil {
			return err
		}
		return filepath.SkipDir // .git ディレクトリの下は見ない
	})
}
