// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// CounterConfig bundles the configuration for creating a Counter metric.
type CounterConfig struct {
	MetricConfig
	Callback MetricCallback // Callback function for metric.
}

// localCounter is the local implements for interface Counter.
type localCounter struct {
	Metric
	CounterConfig
	CounterPerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localCounter)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localCounter)(nil)
)

// NewCounter creates and returns a new Counter.
func NewCounter(config CounterConfig) Counter {
	m := &localCounter{
		Metric:           newMetric(MetricTypeCounter, config.MetricConfig),
		CounterConfig:    config,
		CounterPerformer: newNoopCounterPerformer(),
	}
	if globalProvider != nil {
		m.Init(globalProvider)
	}
	allMetrics = append(allMetrics, m)
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localCounter) Init(provider Provider) {
	if _, ok := l.CounterPerformer.(noopCounterPerformer); !ok {
		return
	}
	l.CounterPerformer = provider.Performer().Counter(l.CounterConfig)
}

// Performer exports internal Performer.
func (l *localCounter) Performer() any {
	return l.CounterPerformer
}
