// Package scheduler runs functions at fixed intervals until stopped.
package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Job is a function run immediately on start and then every Interval.
type Job struct {
	Name     string
	Interval time.Duration
	// Singleton prevents a run from starting while the previous one is still
	// running. Ticks missed meanwhile are collapsed into a single run.
	Singleton bool
	Run       func(context.Context)
}

// Scheduler runs a fixed set of jobs, see Start.
type Scheduler struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// Start validates all the jobs and, only if they are all valid, starts them.
// Jobs receive a context that is cancelled when ctx is done or Stop is called.
func Start(ctx context.Context, jobs ...Job) (*Scheduler, error) {
	var errs []error
	for _, j := range jobs {
		if err := j.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	s := &Scheduler{cancel: cancel}

	s.wg.Add(len(jobs))
	for _, j := range jobs {
		go s.loop(ctx, j)
	}

	return s, nil
}

// Stop cancels the jobs' context and waits for the running ones to return.
func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
}

func (j Job) validate() error {
	if j.Interval <= 0 {
		return fmt.Errorf("job %q: interval must be positive, got %s", j.Name, j.Interval)
	}
	if j.Run == nil {
		return fmt.Errorf("job %q: nil function", j.Name)
	}
	return nil
}

func (s *Scheduler) loop(ctx context.Context, j Job) {
	defer s.wg.Done()

	t := time.NewTicker(j.Interval)
	defer t.Stop()

	for ctx.Err() == nil {
		if j.Singleton {
			j.Run(ctx)
		} else {
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				j.Run(ctx)
			}()
		}

		select {
		case <-ctx.Done():
		case <-t.C:
		}
	}
}
