// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

// Package gcron implements a cron pattern parser and job runner.
package gcron

import (
	"math"
	"time"

	"github.com/jin502437344/gf/os/gtimer"
)

const (
	STATUS_READY   = gtimer.STATUS_READY
	STATUS_RUNNING = gtimer.STATUS_RUNNING
	STATUS_STOPPED = gtimer.STATUS_STOPPED
	STATUS_CLOSED  = gtimer.STATUS_CLOSED

	gDEFAULT_TIMES = math.MaxInt32
)

var (
	// Default cron object.
	defaultCron = New()
)

// SetLogPath sets the logging folder path for default cron object.
func SetLogPath(path string) {
	defaultCron.SetLogPath(path)
}

// GetLogPath returns the logging folder path of default cron object.
func GetLogPath() string {
	return defaultCron.GetLogPath()
}

// SetLogLevel sets the logging level for default cron object.
func SetLogLevel(level int) {
	defaultCron.SetLogLevel(level)
}

// GetLogLevel returns the logging level for default cron object.
func GetLogLevel() int {
	return defaultCron.GetLogLevel()
}

// Add adds a timed task to default cron object.
// A unique <name> can be bound with the timed task.
// It returns and error if the <name> is already used.
func Add(pattern string, job func(), name ...string) (*Entry, error) {
	return defaultCron.Add(pattern, job, name...)
}

// AddSingleton adds a singleton timed task, to default cron object.
// A singleton timed task is that can only be running one single instance at the same time.
// A unique <name> can be bound with the timed task.
// It returns and error if the <name> is already used.
func AddSingleton(pattern string, job func(), name ...string) (*Entry, error) {
	return defaultCron.AddSingleton(pattern, job, name...)
}

// AddOnce adds a timed task which can be run only once, to default cron object.
// A unique <name> can be bound with the timed task.
// It returns and error if the <name> is already used.
func AddOnce(pattern string, job func(), name ...string) (*Entry, error) {
	return defaultCron.AddOnce(pattern, job, name...)
}

// AddTimes adds a timed task which can be run specified times, to default cron object.
// A unique <name> can be bound with the timed task.
// It returns and error if the <name> is already used.
func AddTimes(pattern string, times int, job func(), name ...string) (*Entry, error) {
	return defaultCron.AddTimes(pattern, times, job, name...)
}

// DelayAdd adds a timed task to default cron object after <delay> time.
func DelayAdd(delay time.Duration, pattern string, job func(), name ...string) {
	defaultCron.DelayAdd(delay, pattern, job, name...)
}

// DelayAddSingleton adds a singleton timed task after <delay> time to default cron object.
func DelayAddSingleton(delay time.Duration, pattern string, job func(), name ...string) {
	defaultCron.DelayAddSingleton(delay, pattern, job, name...)
}

// DelayAddOnce adds a timed task after <delay> time to default cron object.
// This timed task can be run only once.
func DelayAddOnce(delay time.Duration, pattern string, job func(), name ...string) {
	defaultCron.DelayAddOnce(delay, pattern, job, name...)
}

// DelayAddTimes adds a timed task after <delay> time to default cron object.
// This timed task can be run specified times.
func DelayAddTimes(delay time.Duration, pattern string, times int, job func(), name ...string) {
	defaultCron.DelayAddTimes(delay, pattern, times, job, name...)
}

// Search returns a scheduled task with the specified <name>.
// It returns nil if no found.
func Search(name string) *Entry {
	return defaultCron.Search(name)
}

// Remove deletes scheduled task which named <name>.
func Remove(name string) {
	defaultCron.Remove(name)
}

// Size returns the size of the timed tasks of default cron.
func Size() int {
	return defaultCron.Size()
}

// Entries return all timed tasks as slice.
func Entries() []*Entry {
	return defaultCron.Entries()
}

// Start starts running the specified timed task named <name>.
func Start(name string) {
	defaultCron.Start(name)
}

// Stop stops running the specified timed task named <name>.
func Stop(name string) {
	defaultCron.Stop(name)
}
