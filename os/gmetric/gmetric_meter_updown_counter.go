// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localUpDownCounter is the local implements for interface Counter.
type localUpDownCounter struct {
	Metric
	MetricOption
	UpDownCounterPerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localUpDownCounter)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localUpDownCounter)(nil)
)

// UpDownCounter creates and returns a new Counter.
func (meter *localMeter) UpDownCounter(name string, option MetricOption) (UpDownCounter, error) {
	m, err := meter.newMetric(MetricTypeUpDownCounter, name, option)
	if err != nil {
		return nil, err
	}
	updownCounter := &localUpDownCounter{
		Metric:                 m,
		MetricOption:           option,
		UpDownCounterPerformer: newNoopUpDownCounterPerformer(),
	}
	if globalProvider != nil {
		if err = updownCounter.Init(meter.Performer()); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, updownCounter)
	return updownCounter, nil
}

// MustUpDownCounter creates and returns a new Counter.
// It panics if any error occurs.
func (meter *localMeter) MustUpDownCounter(name string, option MetricOption) UpDownCounter {
	m, err := meter.UpDownCounter(name, option)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localUpDownCounter) Init(performer MeterPerformer) (err error) {
	if _, ok := l.UpDownCounterPerformer.(noopUpDownCounterPerformer); !ok {
		// already initialized.
		return
	}
	l.UpDownCounterPerformer, err = performer.UpDownCounterPerformer(
		l.Info().Name(),
		l.MetricOption,
	)
	return
}

// Performer implements interface PerformerExporter, which exports internal Performer of Metric.
// This is usually used by metric implements.
func (l *localUpDownCounter) Performer() any {
	return l.UpDownCounterPerformer
}
