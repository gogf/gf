// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localCounter is the local implements for interface Counter.
type localCounter struct {
	Metric
	MetricConfig
	CounterPerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localCounter)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localCounter)(nil)
)

// NewCounter creates and returns a new Counter.
func NewCounter(config MetricConfig) (Counter, error) {
	baseMetric, err := newMetric(MetricTypeCounter, config)
	if err != nil {
		return nil, err
	}
	m := &localCounter{
		Metric:           baseMetric,
		MetricConfig:     config,
		CounterPerformer: newNoopCounterPerformer(),
	}
	if globalProvider != nil {
		if err = m.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, m)
	return m, nil
}

// MustNewCounter creates and returns a new Counter.
// It panics if any error occurs.
func MustNewCounter(config MetricConfig) Counter {
	m, err := NewCounter(config)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localCounter) Init(provider Provider) (err error) {
	if _, ok := l.CounterPerformer.(noopCounterPerformer); !ok {
		// already initialized.
		return
	}
	l.CounterPerformer, err = provider.Performer().Counter(l.MetricConfig)
	return
}

// Performer exports internal Performer.
func (l *localCounter) Performer() any {
	return l.CounterPerformer
}
