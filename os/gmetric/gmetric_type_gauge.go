// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// GaugeConfig bundles the configuration for creating a Gauge metric.
type GaugeConfig struct {
	MetricConfig
	Callback MetricCallback // Callback function for metric.
}

// localGauge is the local implements for interface Gauge.
type localGauge struct {
	Metric
	GaugeConfig
	GaugePerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localGauge)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localGauge)(nil)
)

// NewGauge creates and returns a new Gauge.
func NewGauge(config GaugeConfig) (Gauge, error) {
	baseMetric, err := newMetric(MetricTypeGauge, config.MetricConfig)
	if err != nil {
		return nil, err
	}
	m := &localGauge{
		Metric:         baseMetric,
		GaugeConfig:    config,
		GaugePerformer: newNoopGaugePerformer(),
	}
	if globalProvider != nil {
		if err = m.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, m)
	return m, nil
}

// MustNewGauge creates and returns a new Gauge.
// It panics if any error occurs.
func MustNewGauge(config GaugeConfig) Gauge {
	m, err := NewGauge(config)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localGauge) Init(provider Provider) (err error) {
	if _, ok := l.GaugePerformer.(noopGaugePerformer); !ok {
		// already initialized.
		return
	}
	l.GaugePerformer, err = provider.Performer().Gauge(l.GaugeConfig)
	return err
}

// Performer exports internal Performer.
func (l *localGauge) Performer() any {
	return l.GaugePerformer
}

func (*localGauge) canBeCallback() {}
