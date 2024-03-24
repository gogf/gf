// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localObservableCounter is the local implements for interface ObservableCounter.
type localObservableCounter struct {
	Metric
	MeterOption
	MetricOption
	ObservableCounterPerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localObservableCounter)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localObservableCounter)(nil)
)

// ObservableCounter creates and returns a new ObservableCounter.
func (meter *localMeter) ObservableCounter(name string, option MetricOption) (ObservableCounter, error) {
	m, err := meter.newMetric(MetricTypeObservableCounter, name, option)
	if err != nil {
		return nil, err
	}
	observableCounter := &localObservableCounter{
		Metric:                     m,
		MeterOption:                meter.MeterOption,
		MetricOption:               option,
		ObservableCounterPerformer: newNoopObservableCounterPerformer(),
	}
	if globalProvider != nil {
		if err = observableCounter.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, observableCounter)
	return observableCounter, nil
}

// MustObservableCounter creates and returns a new ObservableCounter.
// It panics if any error occurs.
func (meter *localMeter) MustObservableCounter(name string, option MetricOption) ObservableCounter {
	m, err := meter.ObservableCounter(name, option)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localObservableCounter) Init(provider Provider) (err error) {
	if _, ok := l.ObservableCounterPerformer.(noopObservableCounterPerformer); !ok {
		// already initialized.
		return
	}
	l.ObservableCounterPerformer, err = provider.MeterPerformer(l.MeterOption).ObservableCounterPerformer(
		l.Info().Name(),
		l.MetricOption,
	)
	return err
}

// Performer implements interface PerformerExporter, which exports internal Performer of Metric.
// This is usually used by metric implements.
func (l *localObservableCounter) Performer() any {
	return l.ObservableCounterPerformer
}
