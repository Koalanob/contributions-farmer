package config

import (
	"log"
	"os"
	"strings"
)

const (
	defaultReposPath = "repos"
	defaultFileName  = "data.txt"

	defaultTargetRepo = "mustafa"
	defaultRepoPrefix = "farmer_"
)

type Config struct {
	AccessToken  string `mapstructure:"ACCESS_TOKEN"`
	ClassicToken string `mapstructure:"CLASSIC_TOKEN"`

	UserName  string `mapstructure:"USER_NAME"`
	UserEmail string `mapstructure:"USER_EMAIL"`

	ReposPath string `mapstructure:"REPOSITORIES_PATH"`
	FileName  string `mapstructure:"FILE_NAME"`

	RepositoryPrefix string `mapstructure:"REPOSITORY_PREFIX"`
	TargetRepo       string `mapstructure:"TARGET_REPOSITORY"`
}

func New(path string) (*Config, error) {
	config := &Config{}
	viperConfig, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}

	if viperConfig.AccessToken == "" {
		log.Fatalf("Access token cannot be nil")
		os.Exit(1)
	}

	if viperConfig.ClassicToken == "" {
		log.Fatalf("Classic token cannot be nil")
		os.Exit(1)
	}

	if viperConfig.UserName == "" {
		log.Fatalf("User name cannot be nil")
		os.Exit(1)
	}

	if viperConfig.UserEmail == "" {
		log.Fatalf("User email cannot be nil")
		os.Exit(1)
	}

	config.AccessToken = viperConfig.AccessToken
	config.ClassicToken = viperConfig.ClassicToken
	config.UserName = viperConfig.UserName
	config.UserEmail = viperConfig.UserEmail
	config.FileName = viperConfig.FileName
	config.ReposPath = viperConfig.ReposPath
	config.RepositoryPrefix = viperConfig.RepositoryPrefix

	if strings.Contains(viperConfig.TargetRepo, defaultRepoPrefix) {
		config.TargetRepo = viperConfig.TargetRepo
	} else {
		config.TargetRepo = config.RepositoryPrefix + viperConfig.TargetRepo
	}

	return config, nil
}
