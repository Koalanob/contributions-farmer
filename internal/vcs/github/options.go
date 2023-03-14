package github

import (
	"github.com/robotiksuperb/contributions-farmer/internal/common"
	"github.com/robotiksuperb/contributions-farmer/internal/fm"
)

type Config struct {
	fm.FileManager

	AccessToken  string
	ClassicToken string

	FarmerPrefix string

	ReposFolder string
	RepoFolder  string
	FileName    string

	Username string
	Email    string
}

type OptionFn = common.Options[Config]

func WithAccessToken(token string) OptionFn {
	return func(c *Config) {
		c.AccessToken = token
	}
}

func WithCommonOptions(fm fm.FileManager) OptionFn {
	return func(o *Config) {
		o.FileManager = fm
	}
}

func WithClassicCredentials(token, username, email string) OptionFn {
	return func(c *Config) {
		c.ClassicToken = token
		c.Username = username
		c.Email = email
	}
}

func WithPath(reposFolder string, repoFolder string, filename string, farmerPrefix string) OptionFn {
	return func(c *Config) {
		c.ReposFolder = reposFolder
		c.RepoFolder = repoFolder
		c.FileName = filename

		c.FarmerPrefix = farmerPrefix
	}
}
