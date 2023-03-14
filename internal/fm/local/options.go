package local

import (
	"github.com/robotiksuperb/contributions-farmer/internal/common"
	"github.com/robotiksuperb/contributions-farmer/internal/fm"
)

type Config struct {
	ReposFolder string
	RepoFolder  string
}

type OptionFn = common.Options[Config]

func NewLocalFileManager(opts ...OptionFn) fm.FileManager {
	cfg := Config{}

	for _, fn := range opts {
		fn(&cfg)
	}

	return &localFileManager{
		reposFolder: cfg.ReposFolder,
		repoFolder:  cfg.RepoFolder,
	}
}

func WithPath(reposFolder, repoFolder string) OptionFn {
	return func(c *Config) {
		c.ReposFolder = reposFolder
		c.RepoFolder = repoFolder
	}
}
