// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localObservableUpDownCounter is the local implements for interface ObservableUpDownCounter.
type localObservableUpDownCounter struct {
	Metric
	MeterOption
	MetricOption
	ObservableUpDownCounterPerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localObservableUpDownCounter)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localObservableUpDownCounter)(nil)
)

// ObservableUpDownCounter creates and returns a new ObservableUpDownCounter.
func (meter *localMeter) ObservableUpDownCounter(name string, option MetricOption) (ObservableUpDownCounter, error) {
	m, err := meter.newMetric(MetricTypeObservableUpDownCounter, name, option)
	if err != nil {
		return nil, err
	}
	observableUpDownCounter := &localObservableUpDownCounter{
		Metric:                           m,
		MeterOption:                      meter.MeterOption,
		MetricOption:                     option,
		ObservableUpDownCounterPerformer: newNoopObservableUpDownCounterPerformer(),
	}
	if globalProvider != nil {
		if err = observableUpDownCounter.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, observableUpDownCounter)
	return observableUpDownCounter, nil
}

// MustObservableUpDownCounter creates and returns a new ObservableUpDownCounter.
// It panics if any error occurs.
func (meter *localMeter) MustObservableUpDownCounter(name string, option MetricOption) ObservableUpDownCounter {
	m, err := meter.ObservableCounter(name, option)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localObservableUpDownCounter) Init(provider Provider) (err error) {
	if _, ok := l.ObservableUpDownCounterPerformer.(noopObservableUpDownCounterPerformer); !ok {
		// already initialized.
		return
	}
	l.ObservableUpDownCounterPerformer, err = provider.MeterPerformer(l.MeterOption).ObservableUpDownCounterPerformer(
		l.Info().Name(),
		l.MetricOption,
	)
	return err
}

// Performer implements interface PerformerExporter, which exports internal Performer of Metric.
// This is usually used by metric implements.
func (l *localObservableUpDownCounter) Performer() any {
	return l.ObservableUpDownCounterPerformer
}
