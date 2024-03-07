// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"context"

	"go.opentelemetry.io/otel/metric"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gmetric"
)

// localGaugePerformer is an implementer for interface GaugePerformer.
type localGaugePerformer struct {
	metric.Float64ObservableGauge
	baseObservePerformer iBaseObservePerformer
}

// newGaugePerformer creates and returns a GaugePerformer that truly takes action to implement Gauge.
func newGaugePerformer(meter metric.Meter, config gmetric.GaugeConfig) (gmetric.GaugePerformer, error) {
	var (
		baseObservePerformer = newBaseObservePerformer(config.MetricConfig)
		options              = []metric.Float64ObservableGaugeOption{
			metric.WithDescription(config.Help),
			metric.WithUnit(config.Unit),
		}
	)
	if !hasGlobalCallbackMetricSet.Contains(config.MetricKey()) {
		callback := metric.WithFloat64Callback(func(ctx context.Context, observer metric.Float64Observer) error {
			if config.Callback != nil {
				result, err := config.Callback(ctx)
				if err != nil {
					return gerror.WrapCodef(
						gcode.CodeOperationFailed,
						err,
						`callback failed for metric "%s"`, config.Name,
					)
				}
				if result != nil {
					observer.Observe(
						result.Value,
						baseObservePerformer.MergeAttributesToObserveOptions(result.Attributes)...,
					)
				}
				return nil
			}
			observer.Observe(
				baseObservePerformer.GetValue(),
				baseObservePerformer.GetObserveOptions()...,
			)
			return nil
		})
		options = append(options, callback)
	}
	gauge, err := meter.Float64ObservableGauge(config.Name, options...)
	if err != nil {
		return nil, gerror.WrapCodef(
			gcode.CodeInternalError,
			err,
			`create Float64ObservableGauge failed with config: %+v`,
			config,
		)
	}
	return &localGaugePerformer{
		Float64ObservableGauge: gauge,
		baseObservePerformer:   baseObservePerformer,
	}, nil
}

// Set sets the Gauge to an arbitrary value.
func (l *localGaugePerformer) Set(value float64, option ...gmetric.Option) {
	l.baseObservePerformer.SetValue(value)
	l.baseObservePerformer.SetObserveOptionsByOption(option...)
}

// Inc increments the Gauge by 1. Use Add to increment it by arbitrary values.
func (l *localGaugePerformer) Inc(option ...gmetric.Option) {
	l.Add(1, option...)
}

// Dec decrements the Gauge by 1. Use Sub to decrement it by arbitrary values.
func (l *localGaugePerformer) Dec(option ...gmetric.Option) {
	l.Sub(1, option...)
}

// Sub subtracts the given value from the Gauge. (The value can be
// negative, resulting in an increase of the Gauge.)
func (l *localGaugePerformer) Sub(decrement float64, option ...gmetric.Option) {
	l.Add(-decrement, option...)
}

// Add adds the given value to the Gauge. (The value can be negative,
// resulting in a decrease of the Gauge.)
func (l *localGaugePerformer) Add(increment float64, option ...gmetric.Option) {
	l.baseObservePerformer.AddValue(increment)
	l.baseObservePerformer.SetObserveOptionsByOption(option...)
}
