// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localUpDownCounter is the local implements for interface Counter.
type localUpDownCounter struct {
	Metric
	MetricConfig
	UpDownCounterPerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localUpDownCounter)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localUpDownCounter)(nil)
)

// NewUpDownCounter creates and returns a new Counter.
func NewUpDownCounter(config MetricConfig) (UpDownCounter, error) {
	baseMetric, err := newMetric(MetricTypeUpDownCounter, config)
	if err != nil {
		return nil, err
	}
	m := &localUpDownCounter{
		Metric:                 baseMetric,
		MetricConfig:           config,
		UpDownCounterPerformer: newNoopUpDownCounterPerformer(),
	}
	if globalProvider != nil {
		if err = m.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, m)
	return m, nil
}

// MustNewUpDownCounter creates and returns a new Counter.
// It panics if any error occurs.
func MustNewUpDownCounter(config MetricConfig) UpDownCounter {
	m, err := NewUpDownCounter(config)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localUpDownCounter) Init(provider Provider) (err error) {
	if _, ok := l.UpDownCounterPerformer.(noopUpDownCounterPerformer); !ok {
		// already initialized.
		return
	}
	l.UpDownCounterPerformer, err = provider.Performer().UpDownCounter(l.MetricConfig)
	return
}

// Performer implements interface PerformerExporter, which exports internal Performer of Metric.
// This is usually used by metric implements.
func (l *localUpDownCounter) Performer() any {
	return l.UpDownCounterPerformer
}
