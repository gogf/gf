// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// MetricConfig holds the basic options for creating a metric.
type MetricConfig struct {
	// REQUIRED: Name is the name of this metric.
	Name string

	// Help provides information about this Histogram.
	Help string

	// Unit is the unit for metric value.
	Unit string

	// Attributes holds the constant key-value pair description metadata for this metric.
	Attributes Attributes

	// Instrument is the OpenTelemetry instrumentation name to bind this Metric to a global MeterProvider.
	Instrument string

	// InstrumentVersion is the OpenTelemetry instrumentation version to bind this Metric to a global MeterProvider.
	InstrumentVersion string
}

type localMetricInfo struct {
	config     MetricConfig
	metricType MetricType
}

func newMetricInfo(metricType MetricType, config MetricConfig) MetricInfo {
	if config.Name == "" {
		panic(gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`metric name cannot be empty, invalid metric config: %+v`,
			config,
		))
	}
	return &localMetricInfo{
		config:     config,
		metricType: metricType,
	}
}

func (l *localMetricInfo) Name() string {
	return l.config.Name
}

func (l *localMetricInfo) Help() string {
	return l.config.Help
}

func (l *localMetricInfo) Unit() string {
	return l.config.Unit
}

func (l *localMetricInfo) Type() MetricType {
	return l.metricType
}

func (l *localMetricInfo) Attrs() Attributes {
	return l.config.Attributes
}

func (l *localMetricInfo) Instrument() string {
	return l.config.Instrument
}

func (l *localMetricInfo) InstrumentVersion() string {
	return l.config.InstrumentVersion
}
