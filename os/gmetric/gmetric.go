// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmetric provides interface definitions and simple api for metric feature.
package gmetric

import (
	"context"
)

// MetricType is the type of metric.
type MetricType string

const (
	MetricTypeCounter   MetricType = `Counter`   // Counter.
	MetricTypeGauge     MetricType = `Gauge`     // Gauge.
	MetricTypeHistogram MetricType = `Histogram` // Histogram.
)

// Provider manages all Metric exporting.
// Be caution that the Histogram buckets could not be customized if the creation of the Histogram
// is before the creation of Provider.
type Provider interface {
	// SetAsGlobal sets current provider as global meter provider for current process.
	SetAsGlobal()

	// Performer creates and returns a Performer.
	// A Performer can produce types of Metric performer.
	Performer() Performer

	// ForceFlush flushes all pending metrics.
	//
	// This method honors the deadline or cancellation of ctx. An appropriate
	// error will be returned in these situations. There is no guaranteed that all
	// metrics be flushed or all resources have been released in these situations.
	ForceFlush(ctx context.Context) error

	// Shutdown shuts down the Provider flushing all pending metrics and
	// releasing any held computational resources.
	Shutdown(ctx context.Context) error
}

// Performer manages all Metric performer creating.
type Performer interface {
	// Counter creates and returns a CounterPerformer that performs
	// the operations for Counter metric.
	Counter(config CounterConfig) CounterPerformer

	// Gauge creates and returns a GaugePerformer that performs
	// the operations for Gauge metric.
	Gauge(config GaugeConfig) GaugePerformer

	// Histogram creates and returns a HistogramPerformer that performs
	// the operations for Histogram metric.
	Histogram(config HistogramConfig) HistogramPerformer
}

// Metric models a single sample value with its metadata being exported.
type Metric interface {
	// Info returns the basic information of a Metric.
	Info() MetricInfo
}

// MetricInfo exports information of the Metric.
type MetricInfo interface {
	Key() string                  // Key returns the unique string key of the metric.
	Name() string                 // Name returns the name of the metric.
	Help() string                 // Help returns the help description of the metric.
	Unit() string                 // Unit returns the unit name of the metric.
	Type() MetricType             // Type returns the type of the metric.
	Attributes() Attributes       // Attributes returns the constant attribute slice of the metric.
	Instrument() MetricInstrument // Instrument returns the instrument info of the metric.
}

// MetricInstrument exports the instrument information of a metric.
type MetricInstrument interface {
	Name() string    // Name returns the instrument name of the metric.
	Version() string // Version returns the instrument version of the metric.
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

// CounterPerformer performs operations for Counter metric.
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

// GaugePerformer performs operations for Gauge metric.
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

// HistogramPerformer performs operations for Histogram metric.
type HistogramPerformer interface {
	// Record adds a single value to the histogram.
	// The value is usually positive or zero.
	Record(increment float64, option ...Option)
}

// MetricInitializer manages the initialization for Metric.
// It is called internally in Provider creation.
type MetricInitializer interface {
	// Init initializes the Metric in Provider creation.
	// It sets the metric performer which really takes action.
	Init(provider Provider)
}

// PerformerExporter exports internal Performer of Metric.
// It is called internally in Provider creation.
type PerformerExporter interface {
	// Performer exports internal Performer of Metric.
	Performer() any
}

// CallbackResult is the result that a callback should return.
type CallbackResult struct {
	Value      float64    // New metric value after callback.
	Attributes Attributes // Dynamic attributes after callback.
}

// MetricCallback function for metric.
// A Callback is automatically called when metric reader starts reading the metric value.
type MetricCallback func(ctx context.Context) (*CallbackResult, error)

// GlobalCallback function for metric.
type GlobalCallback func(ctx context.Context, m CallbackSetter) error

// CallbackSetter sets the value for certain initialized Metric.
type CallbackSetter interface {
	// Set sets the value for certain initialized Metric.
	Set(m Metric, value float64, option ...Option)
}

var (
	// metrics stores all created Metric by current package.
	allMetrics = make([]Metric, 0)

	// globalProvider is the provider for global usage.
	globalProvider Provider
)

// SetGlobalProvider registers `provider` as the global Provider,
// which means the following metrics creating will be base on the global provider.
func SetGlobalProvider(provider Provider) {
	globalProvider = provider
}

// GetAllMetrics returns all Metric that created by current package.
func GetAllMetrics() []Metric {
	return allMetrics
}
