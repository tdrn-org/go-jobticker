//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package jobticker_test

import (
	"context"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-jobticker"
)

func TestSequentialQueueStartShutdown(t *testing.T) {
	queue := jobticker.StartSequential(t.Context(), time.Minute, time.Minute)
	require.NoError(t, queue.Shutdown(t.Context()))
}

func TestSequentialQueueJobs(t *testing.T) {
	var tickCount int64
	job := func(ctx context.Context, last, now time.Time) error {
		count := atomic.AddInt64(&tickCount, 1)
		slog.Info("job run", slog.String("last", last.Format(time.TimeOnly)), slog.String("now", now.Format(time.TimeOnly)), slog.Int64("count", count))
		return nil
	}
	queue := jobticker.StartSequential(t.Context(), 10*time.Second, time.Minute, job)
	time.Sleep(25 * time.Second)
	require.NoError(t, queue.Shutdown(t.Context()))
	require.Equal(t, int64(2), tickCount)
}
