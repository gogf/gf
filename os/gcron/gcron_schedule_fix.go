// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/internal/intlog"
)

// getFixedSecond checks, fixes and returns the seconds that have delay fix in some seconds.
// Reference: https://github.com/golang/go/issues/14410
func (s *cronSchedule) getFixedSecond(ctx context.Context, t time.Time) int {
	return (t.Second() + s.getFixedTimestampDelta(ctx, t)) % 60
}

// getFixedTimestampDelta checks, fixes and returns the timestamp delta that have delay fix in some seconds.
// The tolerated timestamp delay is `3` seconds in default.
func (s *cronSchedule) getFixedTimestampDelta(ctx context.Context, t time.Time) int {
	var (
		currentTimestamp = t.Unix()
		lastTimestamp    = s.lastTimestamp.Val()
		delta            int
	)
	switch {
	case
		lastTimestamp == currentTimestamp-1:
		lastTimestamp = currentTimestamp

	case
		lastTimestamp == currentTimestamp-2,
		lastTimestamp == currentTimestamp-3,
		lastTimestamp == currentTimestamp:
		lastTimestamp += 1
		delta = 1

	default:
		// Too much delay, let's update the last timestamp to current one.
		intlog.Printf(
			ctx,
			`too much delay, last timestamp "%d", current "%d"`,
			lastTimestamp, currentTimestamp,
		)
		lastTimestamp = currentTimestamp
	}
	s.lastTimestamp.Set(lastTimestamp)
	return delta
}
