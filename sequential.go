//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

// Package jobticker provides functions and types to setup and
// trigger background jobs on a pre-defined ticker.
package jobticker

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// SequentialQueue defines a simple sequential job ticker. Jobs are
// executed in sequential order every time the ticker fires.
type SequentialQueue struct {
	schedule   time.Duration
	deadline   time.Duration
	ticker     *time.Ticker
	lastTick   time.Time
	jobs       []JobFunc
	cancel     context.CancelFunc
	logger     *slog.Logger
	shutdownWG sync.WaitGroup
	mutex      sync.Mutex
}

// StartSequential starts a new sequential job ticker using the given
// schedule, deadline limit and jobs.
func StartSequential(ctx context.Context, schedule, deadline time.Duration, jobs ...JobFunc) *SequentialQueue {
	logger := slog.With(slog.String("jobticker", schedule.String()))
	logger.Info("starting...")
	cancelCtx, cancel := context.WithCancel(ctx)
	q := &SequentialQueue{
		schedule: schedule,
		deadline: deadline,
		ticker:   time.NewTicker(schedule),
		lastTick: time.Now(),
		jobs:     append([]JobFunc{}, jobs...),
		cancel:   cancel,
		logger:   logger,
	}
	q.shutdownWG.Go(func() {
		q.run(cancelCtx)
	})
	return q
}

func (q *SequentialQueue) run(ctx context.Context) {
	defer q.ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			q.logger.Info("stopped")
			return
		case tick := <-q.ticker.C:
			q.runJobs(ctx, tick)
		}
	}
}

func (q *SequentialQueue) runJobs(ctx context.Context, tick time.Time) {
	jobs, last := q.preJobs(tick)
	deadlineCtx, cancel := context.WithDeadline(ctx, tick.Add(q.deadline))
	defer cancel()
	for _, job := range jobs {
		err := job(deadlineCtx, last, tick)
		if err != nil {
			q.logger.Warn("job failure", slog.Any("err", err))
		}
	}
	q.postJobs(tick)
}

func (q *SequentialQueue) preJobs(tick time.Time) ([]JobFunc, time.Time) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if tick.Sub(q.lastTick)-q.schedule > scheduleJitter {
		q.logger.Warn("tick(s) skipped; due to excessive job runtime")
	}
	jobs := make([]JobFunc, len(q.jobs))
	copy(jobs, q.jobs)
	return jobs, q.lastTick
}

func (q *SequentialQueue) postJobs(tick time.Time) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.lastTick = tick
}

// AddJobs adds the given jobs to this job ticker. The
// jobs are picked up the next time the job ticker fires.
func (q *SequentialQueue) AddJobs(jobs ...JobFunc) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.jobs = append(q.jobs, jobs...)
}

// Shutdown stops the job ticker and releases any related
// resources.
func (q *SequentialQueue) Shutdown(ctx context.Context) error {
	q.logger.Info("shutting down...")
	waitCh := make(chan any)
	go func() {
		q.cancel()
		q.shutdownWG.Wait()
		close(waitCh)
	}()
	select {
	case <-waitCh:
		q.logger.Info("shutdown complete")
	case <-ctx.Done():
		q.logger.Info("shutdown pending")
	}
	return nil
}
