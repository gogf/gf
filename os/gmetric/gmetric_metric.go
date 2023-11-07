// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

type localMetric struct {
	metricInfo MetricInfo
}

func newMetric(metricType MetricType, config MetricConfig) Metric {
	return &localMetric{
		metricInfo: newMetricInfo(metricType, config),
	}
}

// MetricInfo returns the basic information of a Metric.
func (l *localMetric) MetricInfo() MetricInfo {
	return l.metricInfo
}
