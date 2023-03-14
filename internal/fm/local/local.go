package local

import (
	"context"
	"fmt"
	"os"

	"github.com/robotiksuperb/contributions-farmer/pkg/utils/random"
)

type localFileManager struct {
	reposFolder string
	repoFolder  string
}

func (l localFileManager) CreateAndEditFile(ctx context.Context, filename string) error {
	fullPath := fmt.Sprintf("./%s/%s/%s", l.reposFolder, l.repoFolder, filename)

	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	if _, err := file.WriteString(random.MakeString(10)); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
