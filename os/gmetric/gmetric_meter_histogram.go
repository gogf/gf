// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localHistogram is the local implements for interface Histogram.
type localHistogram struct {
	Metric
	MeterOption
	MetricOption
	HistogramPerformer
}

var (
	// Check the implements for interface MetricInitializer.
	_ MetricInitializer = (*localHistogram)(nil)
	// Check the implements for interface PerformerExporter.
	_ PerformerExporter = (*localHistogram)(nil)
)

// Histogram creates and returns a new Histogram.
func (meter *localMeter) Histogram(name string, option MetricOption) (Histogram, error) {
	m, err := meter.newMetric(MetricTypeHistogram, name, option)
	if err != nil {
		return nil, err
	}
	histogram := &localHistogram{
		Metric:             m,
		MeterOption:        meter.MeterOption,
		MetricOption:       option,
		HistogramPerformer: newNoopHistogramPerformer(),
	}
	if globalProvider != nil {
		if err = histogram.Init(globalProvider); err != nil {
			return nil, err
		}
	}
	allMetrics = append(allMetrics, histogram)
	return histogram, nil
}

// MustHistogram creates and returns a new Histogram.
// It panics if any error occurs.
func (meter *localMeter) MustHistogram(name string, option MetricOption) Histogram {
	m, err := meter.Histogram(name, option)
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
	l.HistogramPerformer, err = provider.MeterPerformer(l.MeterOption).HistogramPerformer(
		l.Info().Name(),
		l.MetricOption,
	)
	return err
}

// Buckets returns the bucket slice of the Histogram.
func (l *localHistogram) Buckets() []float64 {
	return l.MetricOption.Buckets
}

// Performer implements interface PerformerExporter, which exports internal Performer of Metric.
// This is usually used by metric implements.
func (l *localHistogram) Performer() any {
	return l.HistogramPerformer
}
