// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localMetric implements interface Metric.
type localMetric struct {
	metricInfo MetricInfo
}

// newMetric creates and returns an object that implements interface Metric.
func newMetric(metricType MetricType, config MetricConfig) Metric {
	return &localMetric{
		metricInfo: newMetricInfo(metricType, config),
	}
}

// MetricInfo returns the basic information of a Metric.
func (l *localMetric) MetricInfo() MetricInfo {
	return l.metricInfo
}
