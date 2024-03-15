// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localObservableUpDownCounter is the local implements for interface Counter.
type localObservableUpDownCounter struct {
	Metric
	MetricConfig
	ObservableMetric
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localObservableUpDownCounter)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localObservableUpDownCounter)(nil)
)

// NewObservableUpDownCounter creates and returns a new ObservableUpDownCounter.
func NewObservableUpDownCounter(config MetricConfig) (ObservableCounter, error) {
	baseMetric, err := newMetric(MetricTypeObservableUpDownCounter, config)
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

// MustNewObservableUpDownCounter creates and returns a new ObservableUpDownCounter.
// It panics if any error occurs.
func MustNewObservableUpDownCounter(config MetricConfig) ObservableCounter {
	m, err := NewObservableCounter(config)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localObservableUpDownCounter) Init(provider Provider) (err error) {
	if _, ok := l.ObservableMetric.(noopObservableMetric); !ok {
		// already initialized.
		return
	}
	l.ObservableMetric, err = provider.Performer().ObservableUpDownCounter(l.MetricConfig)
	return
}

// Performer implements interface PerformerExporter, which exports internal Performer of Metric.
// This is usually used by metric implements.
func (l *localObservableUpDownCounter) Performer() any {
	return l.ObservableMetric
}
