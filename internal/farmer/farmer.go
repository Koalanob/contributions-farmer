package farmer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koalacrypt/contributions-farmer/internal/vcs"
	"github.com/koalacrypt/contributions-farmer/pkg/utils/fm"
)

type Farmer interface {
	Run(context.Context) error
}

type activityFarmer struct {
	m   *sync.Mutex
	vcs vcs.VCSProvider

	commitsCounter atomic.Uint64

	repo        string
	concurrency int

	start time.Time
	end   time.Time
}

func New(opts ...OptionFn) (Farmer, error) {
	cfg := Config{}

	for _, fn := range opts {
		fn(&cfg)
	}

	if cfg.Concurrency < 1 || cfg.Concurrency > maxConcurrency {
		cfg.Concurrency = defaultConcurrency
	}

	return &activityFarmer{
		m:   &sync.Mutex{},
		vcs: cfg.VCSProvider,

		repo:        cfg.Repo,
		concurrency: cfg.Concurrency,

		start: cfg.FirstDay,
		end:   cfg.LastDay,
	}, nil
}

func (a *activityFarmer) startWorker(id int, ctx context.Context, wg *sync.WaitGroup, ch chan int) {
	defer wg.Done()

	currentDay := a.start

	for j := range ch {
		a.commitsCounter.Add(1)

		a.m.Lock()

		currentDay = currentDay.AddDate(0, 0, -1)
		if currentDay == a.end {
			currentDay = a.start
			log.Printf("successfully commited from %s to %s | commits count: %d",
				a.start, a.end, a.commitsCounter.Load())
		}

		if err := a.vcs.Commit(ctx, fmt.Sprintf("feat: my cool feature. %d", j), currentDay); err != nil {
			log.Fatalln(err)
		}

		if a.commitsCounter.Load()%1000 == 0 && a.commitsCounter.Load() != 0 {
			fmt.Printf("\n\n\nYou have reached %d commits. Trying to push...\n", a.commitsCounter.Load())
			if err := a.vcs.Push(ctx, a.repo); err != nil {
				log.Fatalln(err)
			}
		}

		if a.commitsCounter.Load()%20000 == 0 && a.commitsCounter.Load() != 0 {
			fmt.Printf("\n\n\nYou have reached %d commits. Restarting the app to improve perfomance...\n", a.commitsCounter.Load())

			if err := fm.RemoveReposFolder(a.repo); err != nil {
				log.Fatalln(err)
			}
			if err := a.vcs.Clone(ctx, a.repo); err != nil {
				log.Fatalln(err)
			}
		}

		a.m.Unlock()
	}
}

func (a *activityFarmer) seedJobs(ctx context.Context, ch chan int) {
	defer close(ch)
	for {
		select {
		case <-ctx.Done():
			if err := fm.RemoveReposFolder(a.repo); err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("successfully deleted %s folder\n", a.repo)
			fmt.Printf("seeder is being closed... \n")
			return
		default:
		}

		ch <- time.Now().Nanosecond()
	}
}

func (a *activityFarmer) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	jobs := make(chan int)

	go a.seedJobs(ctx, jobs)

	if _, err := a.vcs.CreateInitialRepo(ctx, a.repo); err != nil {
		log.Fatalln(err)
	}

	if err := a.vcs.Clone(ctx, a.repo); err != nil {
		log.Fatalln(err)
	}

	for i := 1; i <= a.concurrency; i++ {
		wg.Add(1)
		go a.startWorker(i, ctx, wg, jobs)
	}

	wg.Wait()
	return nil
}
