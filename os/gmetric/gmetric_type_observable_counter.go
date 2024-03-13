// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localObservableCounter is the local implements for interface Counter.
type localObservableCounter struct {
	Metric
	MetricConfig
	ObservableMetric
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localObservableCounter)(nil)
)

// NewObservableCounter creates and returns a new CounterObservable.
func NewObservableCounter(config MetricConfig) (ObservableCounter, error) {
	baseMetric, err := newMetric(MetricTypeObservableCounter, config)
	if err != nil {
		return nil, err
	}
	m := &localObservableCounter{
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

// MustNewObservableCounter creates and returns a new CounterObservable.
// It panics if any error occurs.
func MustNewObservableCounter(config MetricConfig) ObservableCounter {
	m, err := NewObservableCounter(config)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localObservableCounter) Init(provider Provider) (err error) {
	if _, ok := l.ObservableMetric.(noopObservableMetric); !ok {
		// already initialized.
		return
	}
	l.ObservableMetric, err = provider.Performer().ObservableCounter(l.MetricConfig)
	return
}
