// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// GaugeConfig bundles the configuration for creating a Gauge metric.
type GaugeConfig struct {
	MetricConfig
}

// localGauge is the local implements for interface Gauge.
type localGauge struct {
	Metric
	GaugeConfig
	GaugePerformer
}

var (
	// Check the implements for interface Initializer.
	_ Initializer = (*localGauge)(nil)
)

// NewGauge creates and returns a new Gauge.
func NewGauge(config GaugeConfig) Gauge {
	m := &localGauge{
		Metric:         newMetric(MetricTypeGauge, config.MetricConfig),
		GaugeConfig:    config,
		GaugePerformer: newNoopGaugePerformer(),
	}
	metrics = append(metrics, m)
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localGauge) Init(provider Provider) {
	if _, ok := l.GaugePerformer.(noopGaugePerformer); !ok {
		return
	}
	l.GaugePerformer = provider.Meter().GaugePerformer(l.GaugeConfig)
}
