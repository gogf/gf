// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime

import (
	"time"
)

<<<<<<< HEAD
// TimeWrapper is a wrapper for stdlib struct time.Time.
// It's used for overwriting some functions of time.Time, for example: String.
type TimeWrapper struct {
=======
// wrapper is a wrapper for stdlib struct time.Time.
// It's used for overwriting some functions of time.Time, for example: String.
type wrapper struct {
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
	time.Time
}

// String overwrites the String function of time.Time.
<<<<<<< HEAD
func (t TimeWrapper) String() string {
=======
func (t wrapper) String() string {
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
