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

// localObserver implements interface gmetric.Observer.
type localObserver struct {
	metric.Observer
	gmetric.MeterOption
}

// newObserver creates and returns gmetric.Observer.
func newObserver(observer metric.Observer, meterOption gmetric.MeterOption) gmetric.Observer {
	return &localObserver{
		Observer:    observer,
		MeterOption: meterOption,
	}
}

// Observe observes the value for certain initialized Metric.
// It adds the value to total result if the observed Metrics is type of Counter.
// It sets the value as the result if the observed Metrics is type of Gauge.
func (l *localObserver) Observe(om gmetric.ObservableMetric, value float64, option ...gmetric.Option) {
	var (
		m                      = om.(gmetric.Metric)
		constOption            = getConstOptionByMetric(l.MeterOption, m)
		dynamicOption          = getDynamicOptionByMetricOption(option...)
		globalAttributesOption = getGlobalAttributesOption(gmetric.GetGlobalAttributesOption{
			Instrument:        m.Info().Instrument().Name(),
			InstrumentVersion: m.Info().Instrument().Version(),
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
	l.Observer.ObserveFloat64(metricToFloat64Observable(m), value, observeOptions...)
}
