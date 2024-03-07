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
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localHistogram)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localHistogram)(nil)
)

// NewHistogram creates and returns a new Histogram.
func NewHistogram(config HistogramConfig) (Histogram, error) {
	baseMetric, err := newMetric(MetricTypeHistogram, config.MetricConfig)
	if err != nil {
		return nil, err
	}
	m := &localHistogram{
		Metric:             baseMetric,
		HistogramConfig:    config,
		HistogramPerformer: newNoopHistogramPerformer(),
	}
	if globalProvider != nil {
		// Note that, if Histogram is created after Provider is creation,
		// it cannot customize its Buckets.
		if err = m.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, m)
	return m, nil
}

// MustNewHistogram creates and returns a new Histogram.
// It panics if any error occurs.
func MustNewHistogram(config HistogramConfig) Histogram {
	m, err := NewHistogram(config)
	if err != nil {
		panic(err)
	}
	return m
}

// Init initializes the Metric in Provider creation.
func (l *localHistogram) Init(provider Provider) (err error) {
	if _, ok := l.HistogramPerformer.(noopHistogramPerformer); !ok {
		// already initialized.
		return
	}
	l.HistogramPerformer, err = provider.Performer().Histogram(l.HistogramConfig)
	return err
}

// Buckets returns the bucket slice of the Histogram.
func (l *localHistogram) Buckets() []float64 {
	return l.HistogramConfig.Buckets
}

// Performer exports internal Performer.
func (l *localHistogram) Performer() any {
	return l.HistogramPerformer
}
