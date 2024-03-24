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

// localCounterPerformer is an implementer for interface gmetric.CounterPerformer.
type localCounterPerformer struct {
	gmetric.MeterOption
	gmetric.MetricOption
	metric.Float64Counter
	constOption metric.MeasurementOption
}

// newCounterPerformer creates and returns a CounterPerformer that truly takes action to implement Counter.
func (l *localMeterPerformer) newCounterPerformer(
	meter metric.Meter,
	metricName string,
	metricOption gmetric.MetricOption,
) (gmetric.CounterPerformer, error) {
	var (
		options = []metric.Float64CounterOption{
			metric.WithDescription(metricOption.Help),
			metric.WithUnit(metricOption.Unit),
		}
	)
	counter, err := meter.Float64Counter(metricName, options...)
	if err != nil {
		return nil, gerror.WrapCodef(
			gcode.CodeInternalError,
			err,
			`create Float64Counter "%s" failed with option: %s`,
			metricName, gjson.MustEncodeString(metricOption),
		)
	}
	return &localCounterPerformer{
		MetricOption:   metricOption,
		MeterOption:    l.MeterOption,
		Float64Counter: counter,
		constOption:    genConstOptionForMetric(l.MeterOption, metricOption),
	}, nil
}

// Inc increments the counter by 1.
func (l *localCounterPerformer) Inc(ctx context.Context, option ...gmetric.Option) {
	l.Add(ctx, 1, option...)
}

// Add adds the given value to the counter. It panics if the value is < 0.
func (l *localCounterPerformer) Add(ctx context.Context, increment float64, option ...gmetric.Option) {
	l.Float64Counter.Add(ctx, increment, generateAddOptions(l.MeterOption, l.constOption, option...)...)
}
