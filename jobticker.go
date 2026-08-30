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
	"time"
)

// JobFunc defines the jobs to execute. It is called with the
// execution context as well as the last and current tick of
// the job ticker.
type JobFunc func(ctx context.Context, last, now time.Time) error

const scheduleJitter time.Duration = 5 * time.Second
