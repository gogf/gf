// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// MetricConfig holds the basic options for creating a metric.
type MetricConfig struct {
	// Instrument is the OpenTelemetry instrumentation name to bind this Metric to a global MeterProvider.
	Instrument string

	// Name is the name of this metric.
	Name string

	// Help provides information about this Histogram.
	Help string

	// Unit is the unit for metric value.
	Unit string

	// Attributes holds the constant key-value pair description metadata for this metric.
	Attributes []Attribute
}

type localMetric struct {
	config MetricConfig
}

func newMetric(config MetricConfig) Metric {
	return &localMetric{
		config: config,
	}
}

func (l *localMetric) Inst() string {
	return l.config.Instrument
}

func (l *localMetric) Name() string {
	return l.config.Name
}

func (l *localMetric) Help() string {
	return l.config.Help
}

func (l *localMetric) Unit() string {
	return l.config.Unit
}

func (l *localMetric) Attrs() Attributes {
	return l.config.Attributes
}
