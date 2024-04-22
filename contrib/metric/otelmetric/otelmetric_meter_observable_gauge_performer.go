// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"context"

	"go.opentelemetry.io/otel/metric"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gmetric"
)

// localGaugePerformer is an implementer for interface GaugePerformer.
type localObservableGaugePerformer struct {
	gmetric.ObservableMetric
	metric.Float64ObservableGauge
}

// newObservableGaugePerformer creates and returns a GaugePerformer that truly takes action to implement Gauge.
func (l *localMeterPerformer) newObservableGaugePerformer(
	meter metric.Meter,
	metricName string,
	metricOption gmetric.MetricOption,
) (gmetric.ObservableGaugePerformer, error) {
	var (
		options = []metric.Float64ObservableGaugeOption{
			metric.WithDescription(metricOption.Help),
			metric.WithUnit(metricOption.Unit),
		}
	)
	if metricOption.Callback != nil {
		callback := metric.WithFloat64Callback(func(ctx context.Context, observer metric.Float64Observer) error {
			return metricOption.Callback(ctx, l.newMetricObserver(metricOption, observer))
		})
		options = append(options, callback)
	}
	gauge, err := meter.Float64ObservableGauge(metricName, options...)
	if err != nil {
		return nil, gerror.WrapCodef(
			gcode.CodeInternalError,
			err,
			`create Float64ObservableGauge "%s" failed with option: %s`,
			metricName, gjson.MustEncodeString(metricOption),
		)
	}
	return &localObservableGaugePerformer{
		Float64ObservableGauge: gauge,
	}, nil
}
