// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"go.opentelemetry.io/otel/metric"

	"github.com/gogf/gf/v2/os/gmetric"
)

// localCallbackSetter implements interface gmetric.CallbackSetter.
type localCallbackSetter struct {
	observer metric.Observer
}

// newCallbackSetter creates and returns gmetric.CallbackSetter.
func newCallbackSetter(observer metric.Observer) gmetric.CallbackSetter {
	return &localCallbackSetter{
		observer: observer,
	}
}

// Set sets the value and option for current metric in global callback.
func (l *localCallbackSetter) Set(m gmetric.Metric, value float64, option ...gmetric.Option) {
	var (
		constOption    = getConstOptionByMetric(m)
		dynamicOption  = getDynamicOptionByMetricOption(option...)
		observeOptions = make([]metric.ObserveOption, 0)
	)
	if constOption != nil {
		observeOptions = append(observeOptions, constOption)
	}
	if dynamicOption != nil {
		observeOptions = append(observeOptions, dynamicOption)
	}
	l.observer.ObserveFloat64(metricToFloat64Observable(m), value, observeOptions...)
}
