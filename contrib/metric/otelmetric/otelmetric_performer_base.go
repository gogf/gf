// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"go.opentelemetry.io/otel/metric"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/os/gmetric"
)

// iBaseObservePerformer for observable metric only.
type iBaseObservePerformer interface {
	// GetValue returns the value of current observable performer.
	GetValue() float64

	// SetValue sets the value for current observable performer.
	SetValue(value float64)

	// AddValue adds given `delta` to the `value` of current observable performer.
	AddValue(delta float64) float64

	// GetObserveOptions returns the `currentOptions` that is merged with constant and dynamic currentOptions
	// of current observable performer.
	GetObserveOptions() []metric.ObserveOption

	// SetObserveOptionsByOption sets the `currentOptions` for current observable performer with metric
	// option.
	SetObserveOptionsByOption(option ...gmetric.Option)

	// MergeAttributesToObserveOptions merges constant and dynamic attributes and generates observe currentOptions.
	MergeAttributesToObserveOptions(attributes gmetric.Attributes) []metric.ObserveOption
}

// localBaseObservePerformer is a base struct to implement interface Performer.
type localBaseObservePerformer struct {
	config         gmetric.MetricConfig // Metric Configuration.
	value          *gtype.Float64       // Metric value in concurrent-safety.
	constOption    metric.ObserveOption // Converted attributes to key-value pairs.
	currentOptions *gtype.Any           // Merged const and dynamic []metric.ObserveOption.
}

// newBaseObservePerformer create and returns a base iBaseObservePerformer to implement interface Performer.
func newBaseObservePerformer(config gmetric.MetricConfig) iBaseObservePerformer {
	var (
		constOption    = getConstOptionByMetricConfig(config)
		currentOptions = []metric.ObserveOption{constOption} // Initialized with const option.
	)
	return &localBaseObservePerformer{
		config:         config,
		value:          gtype.NewFloat64(),
		currentOptions: gtype.NewAny(currentOptions),
		constOption:    constOption,
	}
}

// GetValue returns the value of current observable performer.
func (l *localBaseObservePerformer) GetValue() float64 {
	return l.value.Val()
}

// SetValue sets the value for current observable performer.
func (l *localBaseObservePerformer) SetValue(value float64) {
	l.value.Set(value)
}

// AddValue adds given `delta` to the `value` of current observable performer.
func (l *localBaseObservePerformer) AddValue(value float64) float64 {
	return l.value.Add(value)
}

// GetObserveOptions returns the observe currentOptions that is merged with constant and dynamic currentOptions
// of current observable performer.
func (l *localBaseObservePerformer) GetObserveOptions() []metric.ObserveOption {
	return l.currentOptions.Val().([]metric.ObserveOption)
}

// SetObserveOptionsByOption sets the `currentOptions` for current observable performer with metric
// option.
func (l *localBaseObservePerformer) SetObserveOptionsByOption(option ...gmetric.Option) {
	var (
		dynamicOption  = getDynamicOptionByMetricOption(option...)
		observeOptions = make([]metric.ObserveOption, 0)
	)
	if l.constOption != nil {
		observeOptions = append(observeOptions, l.constOption)
	}
	if dynamicOption != nil {
		observeOptions = append(observeOptions, dynamicOption)
	}
	l.currentOptions.Set(observeOptions)
}

// MergeAttributesToObserveOptions merges constant and dynamic attributes and generates observe currentOptions.
func (l *localBaseObservePerformer) MergeAttributesToObserveOptions(
	attributes gmetric.Attributes,
) []metric.ObserveOption {
	var observeOptions = make([]metric.ObserveOption, 0)
	if l.constOption != nil {
		observeOptions = append(observeOptions, l.constOption)
	}
	if len(attributes) > 0 {
		observeOptions = append(
			observeOptions,
			metric.WithAttributes(attributesToKeyValues(attributes)...),
		)
	}
	return observeOptions
}
