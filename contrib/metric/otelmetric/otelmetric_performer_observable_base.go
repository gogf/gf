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

// localObservableBase is a base struct to implement interface observable metric.
type localObservableBase struct {
	config      gmetric.MetricConfig // Metric Configuration.
	constOption metric.ObserveOption // Converted attributes to key-value pairs.
}

// newObservableBase create and returns a base iObservablePerformer to implement interface Performer.
func newObservableBase(config gmetric.MetricConfig) *localObservableBase {
	return &localObservableBase{
		config:      config,
		constOption: getConstOptionByMetricConfig(config),
	}
}

// MergeAttributesToObserveOptions merges constant and dynamic attributes and generates observe currentOptions.
func (l *localObservableBase) MergeAttributesToObserveOptions(attributes gmetric.Attributes) []metric.ObserveOption {
	var (
		observeOptions         = make([]metric.ObserveOption, 0)
		globalAttributesOption = getGlobalAttributesOption(gmetric.GetGlobalAttributesOption{
			Instrument:        l.config.Instrument,
			InstrumentVersion: l.config.InstrumentVersion,
		})
	)
	if l.constOption != nil {
		observeOptions = append(observeOptions, l.constOption)
	}
	if globalAttributesOption != nil {
		observeOptions = append(observeOptions, globalAttributesOption)
	}
	if len(attributes) > 0 {
		observeOptions = append(
			observeOptions,
			metric.WithAttributes(attributesToKeyValues(attributes)...),
		)
	}
	return observeOptions
}
