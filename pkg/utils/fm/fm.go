package fm

import (
	"fmt"
	"os"
	"path/filepath"
)

func RemoveReposFolder(repo string) error {
	dirPath := fmt.Sprintf("./repos/%s", repo)

	if err := filepath.Walk(filepath.Join(dirPath, ".git"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(path, 0777)
	}); err != nil {
		return err
	}

	if err := os.RemoveAll(dirPath); err != nil {
		return err
	}

	return nil
}
