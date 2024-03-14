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
	MetricTypeCounter           MetricType = `Counter`           // Counter.
	MetricTypeHistogram         MetricType = `Histogram`         // Histogram.
	MetricTypeObservableCounter MetricType = `ObservableCounter` // ObservableCounter.
	MetricTypeObservableGauge   MetricType = `ObservableGauge`   // ObservableGauge.
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

	// RegisterCallback registers callback on certain metrics.
	// A callback is bound to certain component and version, it is called when the associated metrics are read.
	// Multiple callbacks on the same component and version will be called by their registered sequence.
	RegisterCallback(callback Callback, canBeCallbackMetrics ...ObservableMetric) error

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
	Counter(config MetricConfig) (CounterPerformer, error)

	// Histogram creates and returns a HistogramPerformer that performs
	// the operations for Histogram metric.
	Histogram(config MetricConfig) (HistogramPerformer, error)

	// ObservableCounter creates and returns an ObservableMetric that performs
	// the operations for ObservableCounter metric.
	ObservableCounter(config MetricConfig) (ObservableMetric, error)

	// ObservableGauge creates and returns an ObservableMetric that performs
	// the operations for ObservableGauge metric.
	ObservableGauge(config MetricConfig) (ObservableMetric, error)
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
	Key() AttributeKey // The key for this attribute.
	Value() any        // The value for this attribute.
}

// AttributeKey is the attribute key.
type AttributeKey string

// Option holds the option for perform a metric operation.
type Option struct {
	// Attributes holds the dynamic key-value pair metadata.
	Attributes Attributes
}

// Counter is a Metric that represents a single numerical value that can ever
// goes up or down.
type Counter interface {
	Metric
	CounterPerformer
}

// ObservableCounter is an instrument used to asynchronously
// record float64 measurements once per collection cycle. Observations are only
// made within a callback for this instrument. The value observed is assumed
// the to be the cumulative sum of the count.
type ObservableCounter interface {
	Metric
	ObservableMetric
}

// CounterPerformer performs operations for Counter metric.
type CounterPerformer interface {
	// Inc increments the counter by 1. Use Add to increment it by arbitrary
	// non-negative values.
	Inc(ctx context.Context, option ...Option)

	// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary values.
	Dec(ctx context.Context, option ...Option)

	// Add adds the given value to the counter. It panics if the value is < 0.
	Add(ctx context.Context, increment float64, option ...Option)
}

// ObservableGauge is an instrument used to asynchronously record
// instantaneous float64 measurements once per collection cycle. Observations
// are only made within a callback for this instrument.
type ObservableGauge interface {
	Metric
	ObservableMetric
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

	// globalProvider is the provider for global usage.
	globalProvider Provider
)

// IsEnabled returns whether the metrics feature is enabled.
func IsEnabled() bool {
	return globalProvider != nil
}

// SetGlobalProvider registers `provider` as the global Provider,
// which means the following metrics creating will be base on the global provider.
func SetGlobalProvider(provider Provider) {
	globalProvider = provider
}

// GetAllMetrics returns all Metric that created by current package.
func GetAllMetrics() []Metric {
	return allMetrics
}
