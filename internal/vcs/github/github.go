package github

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/robotiksuperb/contributions-farmer/internal/fm"
	"github.com/robotiksuperb/contributions-farmer/internal/vcs"

	goGithub "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type githubProvider struct {
	client *goGithub.Client
	auth   *http.BasicAuth

	fm fm.FileManager

	prefix string

	reposFolder string
	repoFolder  string
	filename    string

	username string
	email    string
}

func New(opts ...OptionFn) vcs.VCSProvider {
	ctx := context.Background()
	cfg := &Config{}

	for _, fn := range opts {
		fn(cfg)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.AccessToken})

	c := goGithub.NewClient(oauth2.NewClient(ctx, ts))

	flm := cfg.FileManager

	auth := &http.BasicAuth{
		Username: cfg.Username,
		Password: cfg.ClassicToken,
	}

	return &githubProvider{
		client:      c,
		auth:        auth,
		fm:          flm,
		prefix:      cfg.FarmerPrefix,
		reposFolder: cfg.ReposFolder,
		repoFolder:  cfg.RepoFolder,
		filename:    cfg.FileName,
		username:    cfg.Username,
		email:       cfg.Email,
	}
}

func (g *githubProvider) GetFarmerRepos(ctx context.Context, prefix string) ([]string, error) {
	repos, _, err := g.client.Repositories.List(ctx, g.client.BaseURL.User.Username(), &goGithub.RepositoryListOptions{})
	if err != nil {
		return nil, vcs.ErrGetReposFailure
	}

	var reposName []string
	for _, v := range repos {
		repo := v.GetName()
		re := regexp.MustCompile("^" + prefix)

		if re.Match([]byte(repo)) {
			reposName = append(reposName, repo)
		}
	}

	log.Printf("list of farmer repositories -> %v\n", reposName)

	return reposName, nil
}

func (g *githubProvider) CreateInitialRepo(ctx context.Context, name string) (bool, error) {
	repos, err := g.GetFarmerRepos(ctx, g.prefix)
	if err != nil {
		return false, err
	}

	for _, repo := range repos {
		re := regexp.MustCompile("^" + name)

		if re.Match([]byte(repo)) {
			log.Println(vcs.ErrRepoAlreadyExists)
			return false, nil
		}
	}

	log.Println("initial repository not found... trying to create one")
	if err := g.CreateRepo(ctx, name); err != nil {
		return false, err
	}

	return true, nil
}

func (g *githubProvider) CreateRepo(ctx context.Context, repo string) error {
	_, _, err := g.client.Repositories.Create(ctx, "", &goGithub.Repository{
		Name:     goGithub.String(repo),
		Private:  goGithub.Bool(true),
		AutoInit: goGithub.Bool(true),
	})
	if err != nil {
		return vcs.ErrCreateRepoFailure
	}

	log.Printf("successfully created %s repository", repo)
	return nil
}

func (g *githubProvider) DeleteRepo(ctx context.Context, targetRepo string) error {
	repos, err := g.GetFarmerRepos(ctx, g.prefix)
	if err != nil {
		return err
	}

	for _, repo := range repos {

		if repo == targetRepo {
			_, err := g.client.Repositories.Delete(ctx, g.username, repo)

			if err != nil {
				return vcs.ErrDeleteRepoFailure
			}
		}
	}

	log.Printf("successfully deleted %s repo", targetRepo)
	return nil
}

func (g *githubProvider) DeleteAllRepos(ctx context.Context, prefix string) error {
	repos, err := g.GetFarmerRepos(ctx, g.prefix)
	if err != nil {
		return err
	}

	var deletedRepos []string
	for _, repo := range repos {
		re := regexp.MustCompile("^" + prefix)

		if re.Match([]byte(repo)) {

			_, err := g.client.Repositories.Delete(ctx, g.username, repo)
			if err != nil {
				return vcs.ErrDeleteRepoFailure
			}

			deletedRepos = append(deletedRepos, repo)
		}
	}

	log.Printf("successfully deleted repos -> %v\n", deletedRepos)
	return nil
}

func (g *githubProvider) getRepo(ctx context.Context, repo string) (*git.Repository, error) {
	path := fmt.Sprintf("./%s/%s", g.reposFolder, repo)

	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (g *githubProvider) Clone(ctx context.Context, repo string) error {
	log.Printf("trying to clone %s repository...", repo)

	URL := fmt.Sprintf("https://github.com/%s/%s/", g.username, repo)
	path := fmt.Sprintf("./%s/%s/", g.reposFolder, repo)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err = git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
			URL:      URL,
			Auth:     g.auth,
			Progress: os.Stdout,
		})

		if err != nil {
			fmt.Println(err)
			return vcs.ErrCloneFailure
		}
	}

	log.Printf("successfully cloned %s repository", repo)

	return nil
}

func (g *githubProvider) Commit(ctx context.Context, repo, message string, date time.Time) error {
	r, err := g.getRepo(ctx, repo)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	if err := g.fm.CreateAndEditFile(ctx, g.filename); err != nil {
		log.Fatalln(err)
	}

	if _, err = w.Add(g.filename); err != nil {
		return fmt.Errorf("%w: %w", vcs.ErrAddFailure, err)

	}

	if _, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{Name: g.username, Email: g.email, When: date},
	}); err != nil {
		return fmt.Errorf("%w: %w", vcs.ErrCommitFailure, err)

	}

	return nil
}

func (g *githubProvider) Push(ctx context.Context, repo string) error {
	r, err := g.getRepo(ctx, repo)
	if err != nil {
		return err
	}

	if err := r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       g.auth,
		Force:      true,
	}); err != nil {
		return fmt.Errorf("%w: %w", vcs.ErrPushFailure, err)
	}

	return nil
}
