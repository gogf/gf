// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import "fmt"

// localMetricInfo implements interface MetricInfo.
type localMetricInfo struct {
	MetricType
	MetricOption
	InstrumentInfo
	MetricName string
}

// newMetricInfo creates and returns a MetricInfo.
func (meter *localMeter) newMetricInfo(
	metricType MetricType, metricName string, metricOption MetricOption,
) MetricInfo {
	return &localMetricInfo{
		MetricName:     metricName,
		MetricType:     metricType,
		MetricOption:   metricOption,
		InstrumentInfo: meter.newInstrumentInfo(),
	}
}

// Name returns the name of the metric.
func (l *localMetricInfo) Name() string {
	return l.MetricName
}

// Help returns the help description of the metric.
func (l *localMetricInfo) Help() string {
	return l.MetricOption.Help
}

// Unit returns the unit name of the metric.
func (l *localMetricInfo) Unit() string {
	return l.MetricOption.Unit
}

// Type returns the type of the metric.
func (l *localMetricInfo) Type() MetricType {
	return l.MetricType
}

// Attributes returns the constant attribute slice of the metric.
func (l *localMetricInfo) Attributes() Attributes {
	return l.MetricOption.Attributes
}

// Instrument returns the instrument info of the metric.
func (l *localMetricInfo) Instrument() InstrumentInfo {
	return l.InstrumentInfo
}

func (l *localMetricInfo) Key() string {
	if l.Instrument().Name() != "" && l.Instrument().Version() != "" {
		return fmt.Sprintf(
			`%s@%s:%s`,
			l.Instrument().Name(),
			l.Instrument().Version(),
			l.Name(),
		)
	}
	if l.Instrument().Name() != "" && l.Instrument().Version() == "" {
		return fmt.Sprintf(
			`%s:%s`,
			l.Instrument().Name(),
			l.Name(),
		)
	}
	return l.Name()
}
