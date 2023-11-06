// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmetric interface definitions and basic functionalities for metric feature.
package gmetric

import "context"

// Metric models a single sample value with its metadata being exported.
type Metric interface {
	Name() string      // Name returns the name of the metric.
	Help() string      // Help returns the help description of the metric.
	Unit() string      // Unit returns the unit name of the metric.
	Inst() string      // Inst returns the instrument name of the metric.
	Attrs() Attributes // Attrs returns the constant attribute slice of the metric.
}

// Attributes is a slice of Attribute.
type Attributes []Attribute

// Attribute is the key-value pair item for Metric.
type Attribute interface {
	Key() string // The key for this attribute.
	Value() any  // The value for this attribute.
}

// Counter is a Metric that represents a single numerical value that only ever
// goes up. That implies that it cannot be used to count items whose number can
// also go down, e.g. the number of currently running goroutines. Those
// "counters" are represented by Gauges.
//
// A Counter is typically used to count requests served, tasks completed, errors
// occurred, etc.
//
// To create Counter instances, use NewCounter.
type Counter interface {
	// Metric returns the basic information of a Metric.
	Metric() Metric

	// Inc increments the counter by 1. Use Add to increment it by arbitrary
	// non-negative values.
	Inc(ctx context.Context)
	// Add adds the given value to the counter. It panics if the value is <
	// 0.
	Add(ctx context.Context, increment float64)
}

// Gauge is a Metric that represents a single numerical value that can
// arbitrarily go up and down.
//
// A Gauge is typically used for measured values like temperatures or current
// memory usage, but also "counts" that can go up and down, like the number of
// running goroutines.
//
// To create Gauge instances, use NewGauge.
type Gauge interface {
	// Metric returns the basic information of a Metric.
	Metric() Metric

	// Set sets the Gauge to an arbitrary value.
	Set(ctx context.Context, value float64)
	// Inc increments the Gauge by 1. Use Add to increment it by arbitrary
	// values.
	Inc(ctx context.Context)
	// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary
	// values.
	Dec(ctx context.Context)
	// Add adds the given value to the Gauge. (The value can be negative,
	// resulting in a decrease of the Gauge.)
	Add(ctx context.Context, increment float64)
	// Sub subtracts the given value from the Gauge. (The value can be
	// negative, resulting in an increase of the Gauge.)
	Sub(ctx context.Context, decrement float64)
}

// Histogram counts individual observations from an event or sample stream in
// configurable static buckets (or in dynamic sparse buckets as part of the
// experimental Native Histograms, see below for more details). Similar to a
// Summary, it also provides a sum of observations and an observation count.
type Histogram interface {
	// Metric returns the basic information of a Metric.
	Metric() Metric

	// Record adds a single value to the histogram.
	// The value is usually positive or zero.
	Record(ctx context.Context, increment float64, option ...Option)
}
