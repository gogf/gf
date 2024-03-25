// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localUpDownCounter is the local implements for interface UpDownCounter.
type localUpDownCounter struct {
	Metric
	MeterOption
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
		MeterOption:            meter.MeterOption,
		MetricOption:           option,
		UpDownCounterPerformer: newNoopUpDownCounterPerformer(),
	}
	if globalProvider != nil {
		if err = updownCounter.Init(globalProvider); err != nil {
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
func (l *localUpDownCounter) Init(provider Provider) (err error) {
	if _, ok := l.UpDownCounterPerformer.(noopUpDownCounterPerformer); !ok {
		// already initialized.
		return
	}
	l.UpDownCounterPerformer, err = provider.MeterPerformer(l.MeterOption).UpDownCounterPerformer(
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
