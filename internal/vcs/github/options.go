package github

import (
	"github.com/koalacrypt/contributions-farmer/internal/common"
)

type Config struct {
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
