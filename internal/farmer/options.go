package farmer

import (
	"runtime"
	"time"

	"github.com/koalacrypt/contributions-farmer/internal/common"
	"github.com/koalacrypt/contributions-farmer/internal/vcs"
)

const maxConcurrency = 16

var defaultConcurrency = runtime.NumCPU()

type Config struct {
	vcs.VCSProvider
	Repo        string
	Concurrency int
	FirstDay    time.Time
	LastDay     time.Time
}

type OptionFn = common.Options[Config]

func WithCommonOptions(vcs vcs.VCSProvider, repo string, concurrency int, firstDay, lastDay time.Time) OptionFn {
	return func(o *Config) {
		o.VCSProvider = vcs
		o.Repo = repo
		o.Concurrency = concurrency
		o.FirstDay = firstDay
		o.LastDay = lastDay
	}
}
