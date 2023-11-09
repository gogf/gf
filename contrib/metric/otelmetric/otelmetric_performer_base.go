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

	// GetObserveOptions returns the observe options that is merged with constant and dynamic options
	// of current observable performer.
	GetObserveOptions() []metric.ObserveOption

	// SetObserveOptionsByOption sets the observe options for current observable performer with metric
	// option.
	SetObserveOptionsByOption(option ...gmetric.Option)

	// MergeAttributesToObserveOptions merges constant and dynamic attributes and generates observe options.
	MergeAttributesToObserveOptions(attributes gmetric.Attributes) []metric.ObserveOption
}

// localBaseObservePerformer is a base struct to implement interface Performer.
type localBaseObservePerformer struct {
	config           gmetric.MetricConfig // Metric Configuration.
	value            *gtype.Float64       // Metric value in concurrent-safety.
	options          *gtype.Any           // []metric.ObserveOption
	attributesOption metric.ObserveOption // Converted attributes to key-value pairs.
}

// newBaseObservePerformer create and returns a base iBaseObservePerformer to implement interface Performer.
func newBaseObservePerformer(config gmetric.MetricConfig) iBaseObservePerformer {
	var attributesOption metric.ObserveOption
	if len(config.Attributes) > 0 {
		attributesOption = metric.WithAttributes(attributesToKeyValues(config.Attributes)...)
	}
	return &localBaseObservePerformer{
		config:           config,
		value:            gtype.NewFloat64(),
		options:          gtype.NewAny([]metric.ObserveOption{}),
		attributesOption: attributesOption,
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

// GetObserveOptions returns the observe options that is merged with constant and dynamic options
// of current observable performer.
func (l *localBaseObservePerformer) GetObserveOptions() []metric.ObserveOption {
	return l.options.Val().([]metric.ObserveOption)
}

// SetObserveOptionsByOption sets the observe options for current observable performer with metric
// option.
func (l *localBaseObservePerformer) SetObserveOptionsByOption(option ...gmetric.Option) {
	var (
		usedOption     gmetric.Option
		observeOptions = make([]metric.ObserveOption, 0)
	)
	if l.attributesOption != nil {
		observeOptions = append(observeOptions, l.attributesOption)
	}
	if len(option) > 0 {
		usedOption = option[0]
	}
	if len(usedOption.Attributes) > 0 {
		observeOptions = append(
			observeOptions,
			metric.WithAttributes(attributesToKeyValues(usedOption.Attributes)...),
		)
	}
	l.options.Set(observeOptions)
}

// MergeAttributesToObserveOptions merges constant and dynamic attributes and generates observe options.
func (l *localBaseObservePerformer) MergeAttributesToObserveOptions(
	attributes gmetric.Attributes,
) []metric.ObserveOption {
	var observeOptions = make([]metric.ObserveOption, 0)
	if l.attributesOption != nil {
		observeOptions = append(observeOptions, l.attributesOption)
	}
	if len(attributes) > 0 {
		observeOptions = append(
			observeOptions,
			metric.WithAttributes(attributesToKeyValues(attributes)...),
		)
	}
	return observeOptions
}
