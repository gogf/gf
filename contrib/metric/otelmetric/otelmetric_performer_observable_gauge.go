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
type localObservableGaugePerformer struct {
	gmetric.ObservableMetric
	metric.Float64ObservableGauge
}

// newGaugePerformer creates and returns a GaugePerformer that truly takes action to implement Gauge.
func newObservableGaugePerformer(meter metric.Meter, config gmetric.MetricConfig) (gmetric.ObservableMetric, error) {
	var (
		observableBase = newObservableBase(config)
		options        = []metric.Float64ObservableGaugeOption{
			metric.WithDescription(config.Help),
			metric.WithUnit(config.Unit),
		}
	)
	if !hasGlobalCallbackMetricSet.Contains(config.MetricKey()) {
		callback := metric.WithFloat64Callback(func(ctx context.Context, observer metric.Float64Observer) error {
			if config.Callback == nil {
				return nil
			}
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
					observableBase.MergeAttributesToObserveOptions(result.Attributes)...,
				)
			}
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
	return &localObservableGaugePerformer{
		Float64ObservableGauge: gauge,
	}, nil
}
