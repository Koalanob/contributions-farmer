package farmer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robotiksuperb/contributions-farmer/internal/vcs"
)

type Farmer interface {
	Run(context.Context) error
}

type activityFarmer struct {
	m   *sync.Mutex
	vcs vcs.VCSProvider

	lapCounter     atomic.Uint64
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

		currentDay = currentDay.AddDate(0, 0, -1)
		if currentDay == a.end {
			currentDay = a.start
			a.lapCounter.Add(1)
			log.Printf("successfully commited from %s to %s | %d round is complete | commits count: %d",
				a.start, a.end, a.lapCounter.Load(), a.commitsCounter.Load())
		}

		a.m.Lock()
		if err := a.vcs.Commit(ctx, fmt.Sprintf("feat: my cool feature. %d", j), currentDay); err != nil {
			log.Fatalln(err)
		}
		a.m.Unlock()

	}
}

func (a *activityFarmer) seedJobs(ctx context.Context, ch chan int) {
	defer close(ch)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\n\n\n You have reached %d commits. I'll try to push it now ~~~\n\n\n", a.commitsCounter.Load())
			if err := a.vcs.Push(ctx, a.repo); err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("Successfully pushed all commits\n")
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
