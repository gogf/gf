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

// getAndUpdateLastCheckTimestamp checks fixes and returns the last timestamp that have delay fix in some seconds.
func (s *cronSchedule) getAndUpdateLastCheckTimestamp(ctx context.Context, t time.Time) int64 {
	var (
		currentTimestamp   = t.Unix()
		lastCheckTimestamp = s.lastCheckTimestamp.Val()
	)
	switch {
	// Often happens, timer triggers in the same second, but the millisecond is different.
	// Example:
	// lastCheckTimestamp: 2024-03-26 19:47:34.000
	// currentTimestamp:   2024-03-26 19:47:34.999
	case
		lastCheckTimestamp == currentTimestamp:
		lastCheckTimestamp += 1

	// Often happens, no latency.
	// Example:
	// lastCheckTimestamp: 2024-03-26 19:47:34.000
	// currentTimestamp:   2024-03-26 19:47:35.000
	case
		lastCheckTimestamp == currentTimestamp-1:
		lastCheckTimestamp = currentTimestamp

	// Latency in 3 seconds, which can be tolerant.
	// Example:
	// lastCheckTimestamp: 2024-03-26 19:47:31.000、2024-03-26 19:47:32.000
	// currentTimestamp:   2024-03-26 19:47:34.000
	case
		lastCheckTimestamp == currentTimestamp-2,
		lastCheckTimestamp == currentTimestamp-3:
		lastCheckTimestamp += 1

	// Too much latency, it ignores the fix, the cron job might not be triggered.
	default:
		// Too much delay, let's update the last timestamp to current one.
		intlog.Printf(
			ctx,
			`too much latency, last timestamp "%d", current "%d", latency "%d"`,
			lastCheckTimestamp, currentTimestamp, currentTimestamp-lastCheckTimestamp,
		)
		lastCheckTimestamp = currentTimestamp
	}
	s.lastCheckTimestamp.Set(lastCheckTimestamp)
	return lastCheckTimestamp
}
