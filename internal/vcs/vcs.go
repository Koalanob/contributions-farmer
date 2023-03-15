package vcs

import (
	"context"
	"errors"
	"time"
)

var (
	ErrGetReposFailure   = errors.New("error when trying to get repositories list")
	ErrRepoAlreadyExists = errors.New("initial repository already exists")
	ErrCreateRepoFailure = errors.New("error when trying to create repo")
	ErrDeleteRepoFailure = errors.New("error when trying to delete repository")
	ErrCloneFailure      = errors.New("error when trying to clone")
	ErrAddFailure        = errors.New("error when trying to add")
	ErrCommitFailure     = errors.New("error when trying to commit")
	ErrPushFailure       = errors.New("error when trying to push")
)

type VCSProvider interface {
	GetFarmerRepos(ctx context.Context, prefix string) ([]string, error)
	CreateInitialRepo(ctx context.Context, name string) (bool, error)
	CreateRepo(ctx context.Context, name string) error
	DeleteRepo(ctx context.Context, targetRepo string) error
	DeleteAllRepos(ctx context.Context, prefix string) error

	Clone(ctx context.Context, repo string) error
	Commit(ctx context.Context, message string, date time.Time) error
	Push(ctx context.Context, repo string) error
}
