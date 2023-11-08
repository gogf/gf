// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

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
	// Check the implements for interface Initializer.
	_ Initializer = (*localHistogram)(nil)
)

// NewHistogram creates and returns a new Histogram.
func NewHistogram(config HistogramConfig) Histogram {
	m := &localHistogram{
		Metric:             newMetric(MetricTypeHistogram, config.MetricConfig),
		HistogramConfig:    config,
		HistogramPerformer: newNoopHistogramPerformer(),
	}
	metrics = append(metrics, m)
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localHistogram) Init(provider Provider) {
	if _, ok := l.HistogramPerformer.(noopHistogramPerformer); !ok {
		return
	}
	l.HistogramPerformer = provider.Meter().HistogramPerformer(l.HistogramConfig)
}

// Buckets returns the bucket slice of the Histogram.
func (l *localHistogram) Buckets() []float64 {
	return l.HistogramConfig.Buckets
}
