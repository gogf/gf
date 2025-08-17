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
	MetricTypeCounter                 MetricType = `Counter`                 // Counter.
	MetricTypeUpDownCounter           MetricType = `UpDownCounter`           // UpDownCounter.
	MetricTypeHistogram               MetricType = `Histogram`               // Histogram.
	MetricTypeObservableCounter       MetricType = `ObservableCounter`       // ObservableCounter.
	MetricTypeObservableUpDownCounter MetricType = `ObservableUpDownCounter` // ObservableUpDownCounter.
	MetricTypeObservableGauge         MetricType = `ObservableGauge`         // ObservableGauge.
)

const (
	// MetricNamePattern is the regular expression pattern for validating metric name.
	MetricNamePattern = `[\w\.\-\/]`
)

// Provider manages all Metric exporting.
// Be caution that the Histogram buckets could not be customized if the creation of the Histogram
// is before the creation of Provider.
type Provider interface {
	// SetAsGlobal sets current provider as global meter provider for current process,
	// which makes the following metrics creating on this Provider, especially the metrics created in runtime.
	SetAsGlobal()

	// MeterPerformer creates and returns the MeterPerformer that can produce kinds of metric Performer.
	MeterPerformer(config MeterOption) MeterPerformer

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

// MeterPerformer manages all Metric performers creating.
type MeterPerformer interface {
	// CounterPerformer creates and returns a CounterPerformer that performs
	// the operations for Counter metric.
	CounterPerformer(name string, option MetricOption) (CounterPerformer, error)

	// UpDownCounterPerformer creates and returns a UpDownCounterPerformer that performs
	// the operations for UpDownCounter metric.
	UpDownCounterPerformer(name string, option MetricOption) (UpDownCounterPerformer, error)

	// HistogramPerformer creates and returns a HistogramPerformer that performs
	// the operations for Histogram metric.
	HistogramPerformer(name string, option MetricOption) (HistogramPerformer, error)

	// ObservableCounterPerformer creates and returns an ObservableCounterPerformer that performs
	// the operations for ObservableCounter metric.
	ObservableCounterPerformer(name string, option MetricOption) (ObservableCounterPerformer, error)

	// ObservableUpDownCounterPerformer creates and returns an ObservableUpDownCounterPerformer that performs
	// the operations for ObservableUpDownCounter metric.
	ObservableUpDownCounterPerformer(name string, option MetricOption) (ObservableUpDownCounterPerformer, error)

	// ObservableGaugePerformer creates and returns an ObservableGaugePerformer that performs
	// the operations for ObservableGauge metric.
	ObservableGaugePerformer(name string, option MetricOption) (ObservableGaugePerformer, error)

	// RegisterCallback registers callback on certain metrics.
	// A callback is bound to certain component and version, it is called when the associated metrics are read.
	// Multiple callbacks on the same component and version will be called by their registered sequence.
	RegisterCallback(callback Callback, canBeCallbackMetrics ...ObservableMetric) error
}

// MetricOption holds the basic options for creating a metric.
type MetricOption struct {
	// Help provides information about this Histogram.
	// This is an optional configuration for a metric.
	Help string

	// Unit is the unit for metric value.
	// This is an optional configuration for a metric.
	Unit string

	// Attributes holds the constant key-value pair description metadata for this metric.
	// This is an optional configuration for a metric.
	Attributes Attributes

	// Buckets defines the buckets into which observations are counted.
	// For Histogram metric only.
	// A histogram metric uses default buckets if no explicit buckets configured.
	Buckets []float64

	// Callback function for metric, which is called when metric value changes.
	// For observable metric only.
	// If an observable metric has either Callback attribute nor global callback configured, it does nothing.
	Callback MetricCallback
}

// Metric models a single sample value with its metadata being exported.
type Metric interface {
	// Info returns the basic information of a Metric.
	Info() MetricInfo
}

// MetricInfo exports information of the Metric.
type MetricInfo interface {
	Key() string                // Key returns the unique string key of the metric.
	Name() string               // Name returns the name of the metric.
	Help() string               // Help returns the help description of the metric.
	Unit() string               // Unit returns the unit name of the metric.
	Type() MetricType           // Type returns the type of the metric.
	Attributes() Attributes     // Attributes returns the constant attribute slice of the metric.
	Instrument() InstrumentInfo // InstrumentInfo returns the instrument info of the metric.
}

// InstrumentInfo exports the instrument information of a metric.
type InstrumentInfo interface {
	Name() string    // Name returns the instrument name of the metric.
	Version() string // Version returns the instrument version of the metric.
}

// Counter is a Metric that represents a single numerical value that can ever
// goes up.
type Counter interface {
	Metric
	CounterPerformer
}

// CounterPerformer performs operations for Counter metric.
type CounterPerformer interface {
	// Inc increments the counter by 1. Use Add to increment it by arbitrary
	// non-negative values.
	Inc(ctx context.Context, option ...Option)

	// Add adds the given value to the counter. It panics if the value is < 0.
	Add(ctx context.Context, increment float64, option ...Option)
}

// UpDownCounter is a Metric that represents a single numerical value that can ever
// goes up or down.
type UpDownCounter interface {
	Metric
	UpDownCounterPerformer
}

// UpDownCounterPerformer performs operations for UpDownCounter metric.
type UpDownCounterPerformer interface {
	// Inc increments the counter by 1. Use Add to increment it by arbitrary
	// non-negative values.
	Inc(ctx context.Context, option ...Option)

	// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary values.
	Dec(ctx context.Context, option ...Option)

	// Add adds the given value to the counter. It panics if the value is < 0.
	Add(ctx context.Context, increment float64, option ...Option)
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

// ObservableCounter is an instrument used to asynchronously
// record float64 measurements once per collection cycle. Observations are only
// made within a callback for this instrument. The value observed is assumed
// the to be the cumulative sum of the count.
type ObservableCounter interface {
	Metric
	ObservableCounterPerformer
}

// ObservableUpDownCounter is used to synchronously record float64 measurements during a computational
// operation.
type ObservableUpDownCounter interface {
	Metric
	ObservableUpDownCounterPerformer
}

// ObservableGauge is an instrument used to asynchronously record
// instantaneous float64 measurements once per collection cycle. Observations
// are only made within a callback for this instrument.
type ObservableGauge interface {
	Metric
	ObservableGaugePerformer
}

type (
	// ObservableCounterPerformer is performer for observable ObservableCounter.
	ObservableCounterPerformer = ObservableMetric

	// ObservableUpDownCounterPerformer is performer for observable ObservableUpDownCounter.
	ObservableUpDownCounterPerformer = ObservableMetric

	// ObservableGaugePerformer is performer for observable ObservableGauge.
	ObservableGaugePerformer = ObservableMetric
)

// ObservableMetric is an instrument used to asynchronously record
// instantaneous float64 measurements once per collection cycle.
type ObservableMetric interface {
	observable()
}

// MetricInitializer manages the initialization for Metric.
// It is called internally in metric interface implements.
type MetricInitializer interface {
	// Init initializes the Metric in Provider creation.
	// It sets the metric performer which really takes action.
	Init(provider Provider) error
}

// PerformerExporter exports internal Performer of Metric.
// It is called internally in metric interface implements.
type PerformerExporter interface {
	// Performer exports internal Performer of Metric.
	// This is usually used by metric implements.
	Performer() any
}

// MetricCallback is automatically called when metric reader starts reading the metric value.
type MetricCallback func(ctx context.Context, obs MetricObserver) error

// Callback is a function registered with a Meter that makes observations for
// the set of instruments it is registered with. The Observer parameter is used
// to record measurement observations for these instruments.
type Callback func(ctx context.Context, obs Observer) error

// Observer sets the value for certain initialized Metric.
type Observer interface {
	// Observe observes the value for certain initialized Metric.
	// It adds the value to total result if the observed Metrics is type of Counter.
	// It sets the value as the result if the observed Metrics is type of Gauge.
	Observe(m ObservableMetric, value float64, option ...Option)
}

// MetricObserver sets the value for bound Metric.
type MetricObserver interface {
	// Observe observes the value for certain initialized Metric.
	// It adds the value to total result if the observed Metrics is type of Counter.
	// It sets the value as the result if the observed Metrics is type of Gauge.
	Observe(value float64, option ...Option)
}

var (
	// metrics stores all created Metric by current package.
	allMetrics = make([]Metric, 0)
)

// IsEnabled returns whether the metrics feature is enabled.
func IsEnabled() bool {
	return globalProvider != nil
}

// GetAllMetrics returns all Metric that created by current package.
func GetAllMetrics() []Metric {
	return allMetrics
}
