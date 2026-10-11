package scheduler

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

type recorder struct {
	start   time.Time
	mu      sync.Mutex
	starts  []time.Duration
	running int
	maxRun  int
}

func newRecorder() *recorder {
	return &recorder{start: time.Now()}
}

func (r *recorder) run(d time.Duration) func(context.Context) {
	return func(ctx context.Context) {
		r.mu.Lock()
		r.starts = append(r.starts, time.Since(r.start))
		r.running++
		r.maxRun = max(r.maxRun, r.running)
		r.mu.Unlock()

		select {
		case <-ctx.Done():
		case <-time.After(d):
		}

		r.mu.Lock()
		r.running--
		r.mu.Unlock()
	}
}

func (r *recorder) snapshot() (starts []time.Duration, maxRunning int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]time.Duration(nil), r.starts...), r.maxRun
}

func mustStart(t *testing.T, ctx context.Context, jobs ...Job) *Scheduler {
	t.Helper()
	s, err := Start(ctx, jobs...)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	t.Cleanup(s.Stop)
	return s
}

func TestStartRunsImmediatelyThenEveryInterval(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := newRecorder()
		mustStart(t, t.Context(), Job{Name: "j", Interval: time.Second, Run: r.run(0)})

		synctest.Wait()
		if got, _ := r.snapshot(); len(got) != 1 {
			t.Fatalf("runs before the first interval = %v, want one immediate run", got)
		}

		time.Sleep(3500 * time.Millisecond)
		synctest.Wait()

		got, _ := r.snapshot()
		want := []time.Duration{0, time.Second, 2 * time.Second, 3 * time.Second}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run start times = %v, want %v", got, want)
		}
	})
}

func TestSingletonDoesNotOverlapAndCollapsesMissedTicks(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := newRecorder()
		mustStart(t, t.Context(), Job{Name: "j", Interval: time.Second, Singleton: true, Run: r.run(2500 * time.Millisecond)})

		time.Sleep(9 * time.Second)
		synctest.Wait()

		got, maxRunning := r.snapshot()
		if maxRunning != 1 {
			t.Fatalf("max concurrent runs = %d, want 1", maxRunning)
		}
		want := []time.Duration{0, 2500 * time.Millisecond, 5 * time.Second, 7500 * time.Millisecond}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run start times = %v, want %v", got, want)
		}
	})
}

func TestNonSingletonRunsOverlap(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := newRecorder()
		mustStart(t, t.Context(), Job{Name: "j", Interval: time.Second, Run: r.run(2500 * time.Millisecond)})

		time.Sleep(3500 * time.Millisecond)
		synctest.Wait()

		got, maxRunning := r.snapshot()
		if maxRunning != 3 {
			t.Fatalf("max concurrent runs = %d, want 3", maxRunning)
		}
		want := []time.Duration{0, time.Second, 2 * time.Second, 3 * time.Second}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("run start times = %v, want %v", got, want)
		}
	})
}

func TestStopCancelsRunningJobsAndWaitsForThem(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var running atomic.Int32
		block := func(ctx context.Context) {
			running.Add(1)
			defer running.Add(-1)
			<-ctx.Done()
		}
		r := newRecorder()
		s := mustStart(t, t.Context(),
			Job{Name: "singleton", Interval: time.Second, Singleton: true, Run: block},
			Job{Name: "overlapping", Interval: time.Second, Run: block},
			Job{Name: "counted", Interval: time.Second, Run: r.run(0)},
		)

		time.Sleep(2500 * time.Millisecond)
		synctest.Wait()
		if n := running.Load(); n != 4 {
			t.Fatalf("blocked runs before Stop = %d, want 4", n)
		}

		s.Stop()
		if n := running.Load(); n != 0 {
			t.Fatalf("runs still in flight after Stop = %d, want 0", n)
		}

		before, _ := r.snapshot()
		time.Sleep(time.Hour)
		synctest.Wait()
		if after, _ := r.snapshot(); len(after) != len(before) {
			t.Fatalf("runs after Stop = %d, want 0", len(after)-len(before))
		}
	})
}

func TestParentContextCancelStopsScheduling(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		r := newRecorder()
		mustStart(t, ctx, Job{Name: "j", Interval: time.Second, Run: r.run(0)})

		synctest.Wait()
		cancel()
		time.Sleep(time.Hour)
		synctest.Wait()

		if got, _ := r.snapshot(); len(got) != 1 {
			t.Fatalf("runs after the parent context was cancelled = %v, want only the immediate one", got)
		}
	})
}

func TestStartRejectsInvalidJobsWithoutStartingAny(t *testing.T) {
	noop := func(context.Context) {}
	tests := []struct {
		name string
		job  Job
	}{
		{"zero interval", Job{Name: "bad", Run: noop}},
		{"negative interval", Job{Name: "bad", Interval: -time.Second, Run: noop}},
		{"nil function", Job{Name: "bad", Interval: time.Second}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var ran atomic.Bool
				valid := Job{Name: "valid", Interval: time.Second, Run: func(context.Context) { ran.Store(true) }}

				s, err := Start(t.Context(), valid, tt.job)
				if err == nil {
					s.Stop()
					t.Fatal("Start() error = nil, want an error")
				}
				if !strings.Contains(err.Error(), `"bad"`) {
					t.Errorf("Start() error = %q, want it to name the invalid job", err)
				}

				synctest.Wait()
				if ran.Load() {
					t.Fatal("a valid job ran although Start failed")
				}
			})
		})
	}
}
