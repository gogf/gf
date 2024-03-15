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

// localObservableUpDownCounterPerformer is an implementer for interface CounterPerformer.
type localObservableUpDownCounterPerformer struct {
	gmetric.ObservableMetric
	metric.Float64ObservableUpDownCounter
}

// newObservableUpDownCounterPerformer creates and returns a UpDownCounterPerformer that truly takes action to
// implement ObservableUpDownCounter.
func newObservableUpDownCounterPerformer(meter metric.Meter, config gmetric.MetricConfig) (gmetric.ObservableMetric, error) {
	var (
		options = []metric.Float64ObservableUpDownCounterOption{
			metric.WithDescription(config.Help),
			metric.WithUnit(config.Unit),
		}
	)
	if !hasGlobalCallbackMetricSet.Contains(config.MetricKey()) {
		callback := metric.WithFloat64Callback(func(ctx context.Context, observer metric.Float64Observer) error {
			if config.Callback == nil {
				return nil
			}
			return config.Callback(ctx, newMetricObserver(config, observer))
		})
		options = append(options, callback)
	}
	counter, err := meter.Float64ObservableUpDownCounter(config.Name, options...)
	if err != nil {
		return nil, gerror.WrapCodef(
			gcode.CodeInternalError,
			err,
			`create Float64ObservableUpDownCounter failed with config: %+v`,
			config,
		)
	}
	return &localObservableUpDownCounterPerformer{
		Float64ObservableUpDownCounter: counter,
	}, nil
}
