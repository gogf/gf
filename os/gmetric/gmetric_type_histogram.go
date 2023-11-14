// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import (
	"context"

	"github.com/gogf/gf/v2/internal/intlog"
)

// HistogramConfig bundles the configuration for creating a Histogram metric.
type HistogramConfig struct {
	MetricConfig

	// Buckets defines the buckets into which observations are counted.
	Buckets []float64
}

// localHistogram is the local implements for interface Histogram.
type localHistogram struct {
	Metric
	HistogramConfig
	HistogramPerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localHistogram)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localHistogram)(nil)
)

// NewHistogram creates and returns a new Histogram.
func NewHistogram(config HistogramConfig) Histogram {
	m := &localHistogram{
		Metric:             newMetric(MetricTypeHistogram, config.MetricConfig),
		HistogramConfig:    config,
		HistogramPerformer: newNoopHistogramPerformer(),
	}
	if globalProvider != nil {
		// Note that, if Histogram is created after Provider is creation,
		// it cannot customize its Buckets.
		m.Init(globalProvider)
		intlog.Printf(
			context.Background(),
			`Histogram "%s" is created after Provider creation, it cannot customize its Buckets`,
			config.MetricKey(),
		)
	}
	allMetrics = append(allMetrics, m)
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localHistogram) Init(provider Provider) {
	if _, ok := l.HistogramPerformer.(noopHistogramPerformer); !ok {
		return
	}
	l.HistogramPerformer = provider.Performer().Histogram(l.HistogramConfig)
}

// Buckets returns the bucket slice of the Histogram.
func (l *localHistogram) Buckets() []float64 {
	return l.HistogramConfig.Buckets
}

// Performer exports internal Performer.
func (l *localHistogram) Performer() any {
	return l.HistogramPerformer
}
