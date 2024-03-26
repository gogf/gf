// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localObservableGauge is the local implements for interface ObservableGauge.
type localObservableGauge struct {
	Metric
	MeterOption
	MetricOption
	ObservableGaugePerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localObservableGauge)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localObservableGauge)(nil)
)

// ObservableGauge creates and returns a new ObservableGauge.
func (meter *localMeter) ObservableGauge(name string, option MetricOption) (ObservableGauge, error) {
	m, err := meter.newMetric(MetricTypeObservableGauge, name, option)
	if err != nil {
		return nil, err
	}
	observableGauge := &localObservableGauge{
		Metric:                   m,
		MeterOption:              meter.MeterOption,
		MetricOption:             option,
		ObservableGaugePerformer: newNoopObservableGaugePerformer(),
	}
	if globalProvider != nil {
		if err = observableGauge.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, observableGauge)
	return observableGauge, nil
}

// MustObservableGauge creates and returns a new ObservableGauge.
// It panics if any error occurs.
func (meter *localMeter) MustObservableGauge(name string, option MetricOption) ObservableGauge {
	m, err := meter.ObservableGauge(name, option)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localObservableGauge) Init(provider Provider) (err error) {
	if _, ok := l.ObservableGaugePerformer.(noopObservableGaugePerformer); !ok {
		// already initialized.
		return
	}
	l.ObservableGaugePerformer, err = provider.MeterPerformer(l.MeterOption).ObservableGaugePerformer(
		l.Info().Name(),
		l.MetricOption,
	)
	return err
}

// Performer implements interface PerformerExporter, which exports internal Performer of Metric.
// This is usually used by metric implements.
func (l *localObservableGauge) Performer() any {
	return l.ObservableGaugePerformer
}
