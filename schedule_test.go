//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package jobticker_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-jobticker"
)

func TestScheduleMarshalingValid(t *testing.T) {
	valid := []byte("1 * * * *")

	var spec jobticker.ScheduleSpec
	err := spec.UnmarshalText(valid)
	require.NoError(t, err)
	marshaled, err := spec.MarshalText()
	require.NoError(t, err)
	require.Equal(t, valid, marshaled)
}

func TestScheduleMarshalingInvalid(t *testing.T) {
	invalid := []byte("1 A * * *")

	var spec jobticker.ScheduleSpec
	err := spec.UnmarshalText(invalid)
	require.Error(t, err)
	marshaled, err := spec.MarshalText()
	require.NoError(t, err)
	require.Equal(t, []byte(""), marshaled)
}

func TestScheduleInInterval(t *testing.T) {
	from := time.Date(2026, 8, 30, 0, 0, 0, 0, time.Local)
	to := from.Add(5 * time.Minute)
	scheduleBefore := jobticker.ScheduleSpec("0 * * * *")
	scheduleAfter := jobticker.ScheduleSpec("6 * * * *")
	scheduleIn := jobticker.ScheduleSpec("2 * * * *")
	require.False(t, scheduleBefore.InInterval(from, to))
	require.False(t, scheduleAfter.InInterval(from, to))
	require.True(t, scheduleIn.InInterval(from, to))
}
