// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

type GlobalProvider interface {
	Meter(option MeterOption) Meter
}

type Meter interface {
	Counter(name string, option MetricOption) (Counter, error)
	UpDownCounter(name string, option MetricOption) (UpDownCounter, error)
	Histogram(name string, option MetricOption) (Histogram, error)
	ObservableCounter(name string, option MetricOption) (ObservableCounter, error)
	ObservableUpDownCounter(name string, option MetricOption) (ObservableUpDownCounter, error)
	ObservableGauge(name string, option MetricOption) (ObservableGauge, error)

	MustCounter(name string, option MetricOption) Counter
	MustUpDownCounter(name string, option MetricOption) UpDownCounter
	MustHistogram(name string, option MetricOption) Histogram
	MustObservableCounter(name string, option MetricOption) ObservableCounter
	MustObservableUpDownCounter(name string, option MetricOption) ObservableUpDownCounter
	MustObservableGauge(name string, option MetricOption) ObservableGauge

	// RegisterCallback registers callback on certain metrics.
	// A callback is bound to certain component and version, it is called when the associated metrics are read.
	// Multiple callbacks on the same component and version will be called by their registered sequence.
	RegisterCallback(callback Callback, canBeCallbackMetrics ...ObservableMetric) error

	MustRegisterCallback(callback Callback, canBeCallbackMetrics ...ObservableMetric)
}

type localGlobalProvider struct {
}

var (
	// globalProvider is the provider for global usage.
	globalProvider Provider
)

func GetGlobalProvider() GlobalProvider {
	return &localGlobalProvider{}
}

// SetGlobalProvider registers `provider` as the global Provider,
// which means the following metrics creating will be base on the global provider.
func SetGlobalProvider(provider Provider) {
	globalProvider = provider
}

func (l *localGlobalProvider) Meter(option MeterOption) Meter {
	return newMeter(option)
}
