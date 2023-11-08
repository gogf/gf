// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmetric interface definitions and basic functionalities for metric feature.
package gmetric

import "context"

type MetricType string

const (
	MetricTypeCounter   MetricType = `Counter`
	MetricTypeGauge     MetricType = `Gauge`
	MetricTypeHistogram MetricType = `Histogram`
)

// Provider manages all Metric exporting.
type Provider interface {
	// SetAsGlobal sets current provider as global meter provider.
	SetAsGlobal()

	// Meter creates and returns a Meter.
	Meter(instrument string) Meter

	// ForceFlush flushes all pending telemetry.
	//
	// This method honors the deadline or cancellation of ctx. An appropriate
	// error will be returned in these situations. There is no guaranteed that all
	// telemetry be flushed or all resources have been released in these
	// situations.
	ForceFlush(ctx context.Context) error

	// Shutdown shuts down the Provider flushing all pending telemetry and
	// releasing any held computational resources.
	Shutdown(ctx context.Context) error
}

// Meter manages all Metric performer creating.
type Meter interface {
	CounterPerformer(config CounterConfig) CounterPerformer
	GaugePerformer(config GaugeConfig) GaugePerformer
	HistogramPerformer(config HistogramConfig) HistogramPerformer
}

// Metric models a single sample value with its metadata being exported.
type Metric interface {
	// MetricInfo returns the basic information of a Metric.
	MetricInfo() MetricInfo
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
	Metric
	CounterPerformer
}

type CounterPerformer interface {
	// Inc increments the counter by 1. Use Add to increment it by arbitrary
	// non-negative values.
	Inc(option ...Option)

	// Add adds the given value to the counter. It panics if the value is < 0.
	Add(increment float64, option ...Option)
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
	Metric
	GaugePerformer
}

type GaugePerformer interface {
	// Set sets the Gauge to an arbitrary value.
	Set(value float64, option ...Option)

	// Inc increments the Gauge by 1. Use Add to increment it by arbitrary values.
	Inc(option ...Option)

	// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary values.
	Dec(option ...Option)

	// Add adds the given value to the Gauge. (The value can be negative,
	// resulting in a decrease of the Gauge.)
	Add(increment float64, option ...Option)

	// Sub subtracts the given value from the Gauge. (The value can be
	// negative, resulting in an increase of the Gauge.)
	Sub(decrement float64, option ...Option)
}

// Histogram counts individual observations from an event or sample stream in
// configurable static buckets (or in dynamic sparse buckets as part of the
// experimental Native Histograms, see below for more details). Similar to a
// Summary, it also provides a sum of observations and an observation count.
type Histogram interface {
	Metric
	HistogramPerformer

	// Buckets returns the bucket slice of the Histogram.
	Buckets() []float64
}

type HistogramPerformer interface {
	// Record adds a single value to the histogram.
	// The value is usually positive or zero.
	Record(increment float64, option ...Option)
}

type MetricInfo interface {
	Name() string      // Name returns the name of the metric.
	Help() string      // Help returns the help description of the metric.
	Unit() string      // Unit returns the unit name of the metric.
	Inst() string      // Inst returns the instrument name of the metric.
	Type() MetricType  // Type returns the type of the metric.
	Attrs() Attributes // Attrs returns the constant attribute slice of the metric.
}

// Initializer manages the initialization for Metric.
// It is called internally in Provider creation.
type Initializer interface {
	// Init initializes the Metric in Provider creation.
	Init(provider Provider)
}

var (
	// metrics stores all created Metric.
	metrics = make([]Metric, 0)
)

func GetAllMetrics() []Metric {
	return metrics
}
