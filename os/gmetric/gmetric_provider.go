// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// GlobalProvider hold the entry for creating Meter and Metric.
// The GlobalProvider has only one function for Meter creating, which is designed for convenient usage.
type GlobalProvider interface {
	// Meter creates and returns the Meter by given MeterOption.
	Meter(option MeterOption) Meter
}

// Meter hold the functions for kinds of Metric creating.
type Meter interface {
	// Counter creates and returns a new Counter.
	Counter(name string, option MetricOption) (Counter, error)

	// UpDownCounter creates and returns a new UpDownCounter.
	UpDownCounter(name string, option MetricOption) (UpDownCounter, error)

	// Histogram creates and returns a new Histogram.
	Histogram(name string, option MetricOption) (Histogram, error)

	// ObservableCounter creates and returns a new ObservableCounter.
	ObservableCounter(name string, option MetricOption) (ObservableCounter, error)

	// ObservableUpDownCounter creates and returns a new ObservableUpDownCounter.
	ObservableUpDownCounter(name string, option MetricOption) (ObservableUpDownCounter, error)

	// ObservableGauge creates and returns a new ObservableGauge.
	ObservableGauge(name string, option MetricOption) (ObservableGauge, error)

	// MustCounter creates and returns a new Counter.
	// It panics if any error occurs.
	MustCounter(name string, option MetricOption) Counter

	// MustUpDownCounter creates and returns a new UpDownCounter.
	// It panics if any error occurs.
	MustUpDownCounter(name string, option MetricOption) UpDownCounter

	// MustHistogram creates and returns a new Histogram.
	// It panics if any error occurs.
	MustHistogram(name string, option MetricOption) Histogram

	// MustObservableCounter creates and returns a new ObservableCounter.
	// It panics if any error occurs.
	MustObservableCounter(name string, option MetricOption) ObservableCounter

	// MustObservableUpDownCounter creates and returns a new ObservableUpDownCounter.
	// It panics if any error occurs.
	MustObservableUpDownCounter(name string, option MetricOption) ObservableUpDownCounter

	// MustObservableGauge creates and returns a new ObservableGauge.
	// It panics if any error occurs.
	MustObservableGauge(name string, option MetricOption) ObservableGauge

	// RegisterCallback registers callback on certain metrics.
	// A callback is bound to certain component and version, it is called when the associated metrics are read.
	// Multiple callbacks on the same component and version will be called by their registered sequence.
	RegisterCallback(callback Callback, canBeCallbackMetrics ...ObservableMetric) error

	// MustRegisterCallback performs as RegisterCallback, but it panics if any error occurs.
	MustRegisterCallback(callback Callback, canBeCallbackMetrics ...ObservableMetric)
}

type localGlobalProvider struct {
}

var (
	// globalProvider is the provider for global usage.
	globalProvider Provider
)

// GetGlobalProvider retrieves the GetGlobalProvider instance.
func GetGlobalProvider() GlobalProvider {
	return &localGlobalProvider{}
}

// SetGlobalProvider registers `provider` as the global Provider,
// which means the following metrics creating will be base on the global provider.
func SetGlobalProvider(provider Provider) {
	globalProvider = provider
}

// Meter creates and returns the Meter by given MeterOption.
func (l *localGlobalProvider) Meter(option MeterOption) Meter {
	return newMeter(option)
}
