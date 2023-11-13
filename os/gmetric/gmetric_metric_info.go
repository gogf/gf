// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import (
	"fmt"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// MetricConfig holds the basic options for creating a metric.
type MetricConfig struct {
	// REQUIRED: Name is the name of this metric.
	// If no Name given, Metric creation panics.
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

// localMetricInfo implements interface MetricInfo.
type localMetricInfo struct {
	config     MetricConfig
	metricType MetricType
}

// newMetricInfo creates and returns a MetricInfo.
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

// Name returns the name of the metric.
func (l *localMetricInfo) Name() string {
	return l.config.Name
}

// Help returns the help description of the metric.
func (l *localMetricInfo) Help() string {
	return l.config.Help
}

// Unit returns the unit name of the metric.
func (l *localMetricInfo) Unit() string {
	return l.config.Unit
}

// Type returns the type of the metric.
func (l *localMetricInfo) Type() MetricType {
	return l.metricType
}

// Attributes returns the constant attribute slice of the metric.
func (l *localMetricInfo) Attributes() Attributes {
	return l.config.Attributes
}

// Instrument returns the instrument name of the metric.
func (l *localMetricInfo) Instrument() string {
	return l.config.Instrument
}

// InstrumentVersion returns the instrument version of the metric.
func (l *localMetricInfo) InstrumentVersion() string {
	return l.config.InstrumentVersion
}

// Key returns the unique string key for the metric.
func (l *localMetricInfo) Key() string {
	return l.config.MetricKey()
}

// MetricKey returns the unique string key for the metric.
func (c MetricConfig) MetricKey() string {
	if c.Instrument != "" && c.InstrumentVersion != "" {
		return fmt.Sprintf(
			`%s@%s:%s`,
			c.Instrument,
			c.InstrumentVersion,
			c.Name,
		)
	}
	if c.Instrument != "" && c.InstrumentVersion == "" {
		return fmt.Sprintf(
			`%s:%s`,
			c.Instrument,
			c.Name,
		)
	}
	return c.Name
}
