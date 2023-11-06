// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// HistogramConfig bundles the configuration for creating a Histogram metric.
type HistogramConfig struct {
	MetricConfig

	// Buckets defines the buckets into which observations are counted.
	Buckets []float64
}

// localHistogram is the local implements for interface Histogram.
type localHistogram struct {
	metric     Metric
	histogram  metric.Float64Histogram
	attributes []attribute.KeyValue
}

// NewHistogram creates a new Histogram based on the provided HistogramOpts. It
// panics if the buckets in HistogramOpts are not in strictly increasing order.
//
// The returned implementation also implements ExemplarObserver. It is safe to
// perform the corresponding type assertion. Exemplars are tracked separately
// for each bucket.
func NewHistogram(config HistogramConfig) Histogram {
	histogram, err := otel.Meter(config.Instrument).
		Float64Histogram(
			config.Name,
			metric.WithDescription(config.Help),
			metric.WithUnit(config.Unit),
		)
	if err != nil {
		panic(gerror.WrapCodef(
			gcode.CodeInternalError,
			err,
			`error creating Histogram with HistogramConfig: %+v`,
			config,
		))
	}
	return &localHistogram{
		metric:     newMetric(config.MetricConfig),
		histogram:  histogram,
		attributes: attributesToKeyValues(config.Attributes),
	}
}

// Metric returns the basic information of a Metric.
func (l *localHistogram) Metric() Metric {
	return l.metric
}

// Record adds a single value to the histogram. The value is usually positive or zero.
func (l *localHistogram) Record(ctx context.Context, increment float64, option ...Option) {
	var (
		recordOpts  = make([]metric.RecordOption, 0)
		measureOpts = optionToMeasureOption(option...)
	)
	for _, opt := range measureOpts {
		recordOpts = append(recordOpts, opt)
	}
	l.histogram.Record(ctx, increment, recordOpts...)
}
