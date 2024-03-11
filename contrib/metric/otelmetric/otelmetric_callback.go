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

// localCallbackSetter implements interface gmetric.CallbackObserver.
type localCallbackSetter struct {
	observer metric.Observer
}

// newCallbackObserver creates and returns gmetric.CallbackObserver.
func newCallbackObserver(observer metric.Observer) gmetric.CallbackObserver {
	return &localCallbackSetter{
		observer: observer,
	}
}

// Observe observes the value for certain initialized Metric.
// It adds the value to total result if the observed Metrics is type of Counter.
// It sets the value as the result if the observed Metrics is type of Gauge.
func (l *localCallbackSetter) Observe(m gmetric.CanBeCallbackMetric, value float64, option ...gmetric.Option) {
	var (
		constOption            = getConstOptionByMetric(m)
		dynamicOption          = getDynamicOptionByMetricOption(option...)
		globalAttributesOption = getGlobalAttributesOption(gmetric.GetGlobalAttributesOption{
			Instrument:        m.Info().Instrument().Name(),
			InstrumentVersion: m.Info().Instrument().Version(),
		})
		observeOptions = make([]metric.ObserveOption, 0)
	)
	if constOption != nil {
		observeOptions = append(observeOptions, constOption)
	}
	if globalAttributesOption != nil {
		observeOptions = append(observeOptions, globalAttributesOption)
	}
	if dynamicOption != nil {
		observeOptions = append(observeOptions, dynamicOption)
	}
	l.observer.ObserveFloat64(metricToFloat64Observable(m), value, observeOptions...)
}
