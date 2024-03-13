// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localObservableGauge is the local implements for interface Counter.
type localObservableGauge struct {
	Metric
	MetricConfig
	ObservableMetric
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localObservableGauge)(nil)
)

// NewObservableGauge creates and returns a new CounterObservable.
func NewObservableGauge(config MetricConfig) (ObservableCounter, error) {
	baseMetric, err := newMetric(MetricTypeObservableGauge, config)
	if err != nil {
		return nil, err
	}
	m := &localObservableGauge{
		Metric:           baseMetric,
		MetricConfig:     config,
		ObservableMetric: newNoopObservableMetric(),
	}
	if globalProvider != nil {
		if err = m.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, m)
	return m, nil
}

// MustNewObservableGauge creates and returns a new CounterObservable.
// It panics if any error occurs.
func MustNewObservableGauge(config MetricConfig) ObservableGauge {
	m, err := NewObservableGauge(config)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localObservableGauge) Init(provider Provider) (err error) {
	if _, ok := l.ObservableMetric.(noopObservableMetric); !ok {
		// already initialized.
		return
	}
	l.ObservableMetric, err = provider.Performer().ObservableGauge(l.MetricConfig)
	return
}
