// Copyright 2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gtime

import (
	"time"
)

// TimeWrapper is a wrapper for stdlib struct time.Time.
// It's used for overwriting some functions of time.Time, for example: String.
type TimeWrapper struct {
	time.Time
}

// String overwrites the String function of time.Time.
func (t TimeWrapper) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
