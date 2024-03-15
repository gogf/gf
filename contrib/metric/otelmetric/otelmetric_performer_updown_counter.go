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

// localUpDownCounterPerformer is an implementer for interface gmetric.UpDownCounterPerformer.
type localUpDownCounterPerformer struct {
	metric.Float64UpDownCounter
	config      gmetric.MetricConfig
	constOption metric.MeasurementOption
}

// newUpDownCounterPerformer creates and returns a CounterPerformer that truly takes action to implement Counter.
func newUpDownCounterPerformer(meter metric.Meter, config gmetric.MetricConfig) (gmetric.UpDownCounterPerformer, error) {
	var (
		options = []metric.Float64UpDownCounterOption{
			metric.WithDescription(config.Help),
			metric.WithUnit(config.Unit),
		}
	)
	counter, err := meter.Float64UpDownCounter(config.Name, options...)
	if err != nil {
		return nil, gerror.WrapCodef(
			gcode.CodeInternalError,
			err,
			`create Float64Counter failed with config: %+v`,
			config,
		)
	}
	return &localUpDownCounterPerformer{
		Float64UpDownCounter: counter,
		config:               config,
		constOption:          getConstOptionByMetricConfig(config),
	}, nil
}

// Inc increments the counter by 1.
func (l *localUpDownCounterPerformer) Inc(ctx context.Context, option ...gmetric.Option) {
	l.Add(ctx, 1, option...)
}

// Dec decrements the counter by 1.
func (l *localUpDownCounterPerformer) Dec(ctx context.Context, option ...gmetric.Option) {
	l.Add(ctx, -1, option...)
}

// Add adds the given value to the counter.
func (l *localUpDownCounterPerformer) Add(ctx context.Context, increment float64, option ...gmetric.Option) {
	l.Float64UpDownCounter.Add(ctx, increment, generateAddOptions(l.config, l.constOption, option...)...)
}
