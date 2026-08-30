//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package jobticker

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/adhocore/gronx"
)

// ScheduleSpec defines a schedule based on a cron expression
// ([https://en.wikipedia.org/wiki/Cron]).
//
// ScheduleSpec implements [encoding.TextMarshaler] and
// [encoding.TextUnmarshaler].
type ScheduleSpec string

// String gets the schedule expression as string.
func (spec ScheduleSpec) String() string {
	return string(spec)
}

// see [encoding.TextMarshaler].
func (spec ScheduleSpec) MarshalText() ([]byte, error) {
	return []byte(spec), nil
}

// see [encoding.TextUnmarshaler].
func (spec *ScheduleSpec) UnmarshalText(text []byte) error {
	expressionString := string(text)
	if expressionString == "" {
		*spec = ScheduleSpec(expressionString)
		return nil
	}
	if !gronx.IsValid(expressionString) {
		return fmt.Errorf("invalid cron expression '%s'", expressionString)
	}
	*spec = ScheduleSpec(expressionString)
	return nil
}

// InInterval determines whether this schedule lies within the
// given time interval.
func (spec ScheduleSpec) InInterval(from, to time.Time) bool {
	if string(spec) == "" {
		return false
	}
	next, err := gronx.NextTickAfter(string(spec), from, false)
	if err != nil {
		slog.Warn("invalid cron expression", slog.String("expression", string(spec)), slog.Any("err", err))
		return false
	}
	return !next.After(to)
}
