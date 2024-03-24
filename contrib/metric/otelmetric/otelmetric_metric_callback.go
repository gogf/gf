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

// localMetricObserver implements interface gmetric.CallbackObserver.
type localMetricObserver struct {
	gmetric.MeterOption
	gmetric.MetricOption
	metric.Float64Observer
}

func (l *localMeterPerformer) newMetricObserver(
	metricOption gmetric.MetricOption,
	float64Observer metric.Float64Observer,
) gmetric.MetricObserver {
	return &localMetricObserver{
		MeterOption:     l.MeterOption,
		MetricOption:    metricOption,
		Float64Observer: float64Observer,
	}
}

// Observe observes the value for certain initialized Metric.
// It adds the value to total result if the observed Metrics is type of Counter.
// It sets the value as the result if the observed Metrics is type of Gauge.
func (l *localMetricObserver) Observe(value float64, option ...gmetric.Option) {
	var (
		constOption            = genConstOptionForMetric(l.MeterOption, l.MetricOption)
		dynamicOption          = getDynamicOptionByMetricOption(option...)
		globalAttributesOption = getGlobalAttributesOption(gmetric.GetGlobalAttributesOption{
			Instrument:        l.Instrument,
			InstrumentVersion: l.InstrumentVersion,
		})
		observeOptions = make([]metric.ObserveOption, 0)
	)
	if globalAttributesOption != nil {
		observeOptions = append(observeOptions, globalAttributesOption)
	}
	if constOption != nil {
		observeOptions = append(observeOptions, constOption)
	}
	if dynamicOption != nil {
		observeOptions = append(observeOptions, dynamicOption)
	}
	l.Float64Observer.Observe(value, observeOptions...)
}
