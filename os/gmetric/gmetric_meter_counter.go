// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localCounter is the local implements for interface Counter.
type localCounter struct {
	Metric
	MeterOption
	MetricOption
	CounterPerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localCounter)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localCounter)(nil)
)

// Counter creates and returns a new Counter.
func (meter *localMeter) Counter(name string, option MetricOption) (Counter, error) {
	m, err := meter.newMetric(MetricTypeCounter, name, option)
	if err != nil {
		return nil, err
	}
	counter := &localCounter{
		Metric:           m,
		MeterOption:      meter.MeterOption,
		MetricOption:     option,
		CounterPerformer: newNoopCounterPerformer(),
	}
	if globalProvider != nil {
		if err = counter.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, counter)
	return counter, nil
}

// MustCounter creates and returns a new Counter.
// It panics if any error occurs.
func (meter *localMeter) MustCounter(name string, option MetricOption) Counter {
	m, err := meter.Counter(name, option)
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
	l.CounterPerformer, err = provider.MeterPerformer(l.MeterOption).CounterPerformer(
		l.Info().Name(),
		l.MetricOption,
	)
	return
}

// Performer implements interface PerformerExporter, which exports internal Performer of Metric.
// This is usually used by metric implements.
func (l *localCounter) Performer() any {
	return l.CounterPerformer
}
