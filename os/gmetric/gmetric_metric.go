// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
)

// localMetric implements interface Metric.
type localMetric struct {
	metricInfo MetricInfo
}

// newMetric creates and returns an object that implements interface Metric.
func newMetric(metricType MetricType, config MetricConfig) (Metric, error) {
	if config.Name == "" {
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`error creating %s metric while given name is empty, config: %s`,
			metricType, gjson.MustEncodeString(config),
		)
	}
	if !gregex.IsMatchString(MetricNamePattern, config.Name) {
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid metric name "%s", should match regular expression pattern "%s"`,
			config.Name, MetricNamePattern,
		)
	}
	return &localMetric{
		metricInfo: newMetricInfo(metricType, config),
	}, nil
}

// Info returns the basic information of a Metric.
func (l *localMetric) Info() MetricInfo {
	return l.metricInfo
}
