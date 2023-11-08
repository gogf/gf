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
	GetValue() float64
	SetValue(value float64)
	AddValue(value float64) float64
	GetObserveOptions() []metric.ObserveOption
	SetObserveOptionsByOption(option ...gmetric.Option)
	MergeToObserveOptions(option ...gmetric.Option) []metric.ObserveOption
}

// localBaseObservePerformer is a base struct to implement interface Performer.
type localBaseObservePerformer struct {
	config           gmetric.MetricConfig // Metric Configuration.
	value            *gtype.Float64       // Metric value in concurrent-safety.
	options          *gtype.Any           // []metric.ObserveOption
	attributesOption metric.ObserveOption // Converted attributes to key-value pairs.
}

func newBaseObservePerformer(config gmetric.MetricConfig) iBaseObservePerformer {
	return &localBaseObservePerformer{
		config:           config,
		value:            gtype.NewFloat64(),
		options:          gtype.NewAny([]metric.ObserveOption{}),
		attributesOption: metric.WithAttributes(attributesToKeyValues(config.Attributes)...),
	}
}

func (l *localBaseObservePerformer) GetValue() float64 {
	return l.value.Val()
}

func (l *localBaseObservePerformer) SetValue(value float64) {
	l.value.Set(value)
}

func (l *localBaseObservePerformer) AddValue(value float64) float64 {
	return l.value.Add(value)
}

func (l *localBaseObservePerformer) GetObserveOptions() []metric.ObserveOption {
	return l.options.Val().([]metric.ObserveOption)
}

func (l *localBaseObservePerformer) SetObserveOptionsByOption(option ...gmetric.Option) {
	l.options.Set(optionToObserveOptions(option...))
}

func (l *localBaseObservePerformer) MergeToObserveOptions(option ...gmetric.Option) []metric.ObserveOption {
	var (
		usedOption     gmetric.Option
		observeOptions = []metric.ObserveOption{l.attributesOption}
	)
	if len(option) > 0 {
		usedOption = option[0]
	}
	if len(usedOption.Attributes) > 0 {
		observeOptions = append(
			observeOptions,
			metric.WithAttributes(attributesToKeyValues(usedOption.Attributes)...),
		)
	}
	return observeOptions
}
