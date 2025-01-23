// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime

import (
	"sync"
	"time"
)

var (
	setStringLayoutMu sync.Mutex
	setStringLayout   = "2006-01-02T15:04:05-07:00"
)

// SetStringLayout sets the default string layout for function String.
// The default string layout is "2006-01-02 15:04:05".
// It's used for overwriting the default string layout for time.Time.String.
//
// You can set it to time.RFC3339 for example.
func SetStringLayout(layout string) {
	setStringLayoutMu.Lock()
	defer setStringLayoutMu.Unlock()
	setStringLayout = layout
}

// SetStringISO8601 sets the default string layout to ISO8601 for function String.
func SetStringISO8601() {
	SetStringLayout("2006-01-02T15:04:05-07:00")
}

// SetStringRFC822 sets the default string layout to RFC822 for function String.
func SetStringRFC822() {
	SetStringLayout("Mon, 02 Jan 06 15:04 MST")
}

// wrapper is a wrapper for stdlib struct time.Time.
// It's used for overwriting some functions of time.Time, for example: String.
type wrapper struct {
	time.Time
}

// String overwrites the String function of time.Time.
func (t wrapper) String() string {
	if t.IsZero() {
		return ""
	}
	if t.Year() == 0 {
		// Only time.
		return t.Format("15:04:05")
	}
	return t.Format(setStringLayout)
}
