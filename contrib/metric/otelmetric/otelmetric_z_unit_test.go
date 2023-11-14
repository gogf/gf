// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric_test

import (
	"context"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"testing"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = gctx.New()
			counter = gmetric.NewCounter(gmetric.CounterConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.counter",
					Help: "This is a simple demo for Counter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_a", 1),
					},
					Instrument:        "github.com/gogf/gf/example/metric/basic",
					InstrumentVersion: "v1.0",
				},
			})
			gauge = gmetric.NewGauge(gmetric.GaugeConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.gauge",
					Help: "This is a simple demo for Gauge usage",
					Unit: "bytes",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_b", 2),
					},
					Instrument:        "github.com/gogf/gf/example/metric/basic",
					InstrumentVersion: "v1.1",
				},
			})
			histogram = gmetric.NewHistogram(gmetric.HistogramConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.histogram",
					Help: "This is a simple demo for histogram usage",
					Unit: "ms",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_c", 3),
					},
					Instrument:        "github.com/gogf/gf/example/metric/basic",
					InstrumentVersion: "v1.2",
				},
				Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
			})
		)
		reader := metric.NewManualReader()

		// OpenTelemetry provider.
		provider := otelmetric.MustProvider(metric.WithReader(reader))
		defer provider.Shutdown(ctx)

		// Add value for counter.
		counter.Inc()
		counter.Add(10)

		// Set value for gauge.
		gauge.Set(100)
		gauge.Inc()
		gauge.Sub(1)

		// Record values for histogram.
		histogram.Record(1)
		histogram.Record(20)
		histogram.Record(30)
		histogram.Record(101)
		histogram.Record(2000)
		histogram.Record(9000)
		histogram.Record(20000)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(ctx, &rm)
		t.AssertNil(err)

		content := gjson.MustEncodeString(rm)
		content, err = gregex.ReplaceString(`Time":".+?"`, `Time":""`, content)
		t.AssertNil(err)
		expectContent := `
{"Scope":{"Name":"github.com/gogf/gf/example/metric/basic","Version":"v1.2","SchemaURL":""},"Metrics":[{"Name":"goframe.metric.demo.histogram","Description":"This is a simple demo for histogram usage","Unit":"ms","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_c","Value":{"Type":"INT64","Value":3}}],"StartTime":"","Time":"","Count":7,"Bounds":[0,10,20,50,100,500,1000,2000,5000,10000],"BucketCounts":[0,1,1,1,0,1,0,1,0,1,1],"Min":{},"Max":{},"Sum":31152}],"Temporality":"CumulativeTemporality"}}]}
{"Scope":{"Name":"github.com/gogf/gf/example/metric/basic","Version":"v1.0","SchemaURL":""},"Metrics":[{"Name":"goframe.metric.demo.counter","Description":"This is a simple demo for Counter usage","Unit":"%","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_a","Value":{"Type":"INT64","Value":1}}],"StartTime":"","Time":"","Value":11}],"Temporality":"CumulativeTemporality","IsMonotonic":true}}]}
{"Scope":{"Name":"github.com/gogf/gf/example/metric/basic","Version":"v1.1","SchemaURL":""},"Metrics":[{"Name":"goframe.metric.demo.gauge","Description":"This is a simple demo for Gauge usage","Unit":"bytes","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_b","Value":{"Type":"INT64","Value":2}}],"StartTime":"","Time":"","Value":100}]}}]}
`
		for _, line := range gstr.SplitAndTrim(expectContent, "\n") {
			t.Assert(gstr.Contains(content, line), true)
		}
	})
}

func Test_DynamicAttributes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = gctx.New()
			counter = gmetric.NewCounter(gmetric.CounterConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.counter",
					Help: "This is a simple demo for dynamic attributes",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label", 1),
					},
					Instrument:        "github.com/gogf/gf/example/metric/dynamic_attributes",
					InstrumentVersion: "v1.0",
				},
			})
			dynamicAttributes = gmetric.Option{
				Attributes: gmetric.Attributes{
					gmetric.NewAttribute("dynamic_label", 2),
				},
			}
		)

		reader := metric.NewManualReader()

		// OpenTelemetry provider.
		provider := otelmetric.MustProvider(metric.WithReader(reader))
		defer provider.Shutdown(ctx)

		// Add value for counter.
		counter.Inc(dynamicAttributes)
		counter.Add(10, dynamicAttributes)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(ctx, &rm)
		t.AssertNil(err)

		content := gjson.MustEncodeString(rm)
		content, err = gregex.ReplaceString(`Time":".+?"`, `Time":""`, content)
		t.AssertNil(err)

		expectContent := `
{"Scope":{"Name":"github.com/gogf/gf/example/metric/dynamic_attributes","Version":"v1.0","SchemaURL":""},"Metrics":[{"Name":"goframe.metric.demo.counter","Description":"This is a simple demo for dynamic attributes","Unit":"%","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label","Value":{"Type":"INT64","Value":1}},{"Key":"dynamic_label","Value":{"Type":"INT64","Value":2}}],"StartTime":"","Time":"","Value":11}],"Temporality":"CumulativeTemporality","IsMonotonic":true}}]}
`
		for _, line := range gstr.SplitAndTrim(expectContent, "\n") {
			t.Assert(gstr.Contains(content, line), true)
		}
	})
}

func Test_HistogramBuckets(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx        = gctx.New()
			histogram1 = gmetric.NewHistogram(gmetric.HistogramConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.histogram1",
					Help: "This is a simple demo for histogram usage",
					Unit: "ms",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_a", 1),
					},
					Instrument:        "github.com/gogf/gf/example/metric/histogram_buckets",
					InstrumentVersion: "v1.0",
				},
				Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
			})
			histogram2 = gmetric.NewHistogram(gmetric.HistogramConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.histogram2",
					Help: "This demos we can specify custom buckets in Histogram creating",
					Unit: "",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_b", 2),
					},
					Instrument:        "github.com/gogf/gf/example/metric/histogram_buckets",
					InstrumentVersion: "v1.0",
				},
				Buckets: []float64{100, 200, 300, 400, 500},
			})
		)

		reader := metric.NewManualReader()

		// OpenTelemetry provider.
		provider := otelmetric.MustProvider(metric.WithReader(reader))
		defer provider.Shutdown(ctx)

		// Record values for histogram1.
		histogram1.Record(1)
		histogram1.Record(20)
		histogram1.Record(30)
		histogram1.Record(101)
		histogram1.Record(2000)
		histogram1.Record(9000)
		histogram1.Record(20000)

		// Record values for histogram2.
		histogram2.Record(1)
		histogram2.Record(10)
		histogram2.Record(199)
		histogram2.Record(299)
		histogram2.Record(399)
		histogram2.Record(499)
		histogram2.Record(501)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(ctx, &rm)
		t.AssertNil(err)

		content := gjson.MustEncodeString(rm)
		content, err = gregex.ReplaceString(`Time":".+?"`, `Time":""`, content)
		t.AssertNil(err)

		expectKeyContent := `
{"Scope":{"Name":"github.com/gogf/gf/example/metric/histogram_buckets","Version":"v1.0","SchemaURL":""}
{"Name":"goframe.metric.demo.histogram1","Description":"This is a simple demo for histogram usage","Unit":"ms","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_a","Value":{"Type":"INT64","Value":1}}],"StartTime":"","Time":"","Count":7,"Bounds":[0,10,20,50,100,500,1000,2000,5000,10000],"BucketCounts":[0,1,1,1,0,1,0,1,0,1,1],"Min":{},"Max":{},"Sum":31152}],"Temporality":"CumulativeTemporality"}}
{"Name":"goframe.metric.demo.histogram2","Description":"This demos we can specify custom buckets in Histogram creating","Unit":"","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_b","Value":{"Type":"INT64","Value":2}}],"StartTime":"","Time":"","Count":7,"Bounds":[100,200,300,400,500],"BucketCounts":[2,1,1,1,1,1],"Min":{},"Max":{},"Sum":1908}],"Temporality":"CumulativeTemporality"}}
`
		for _, line := range gstr.SplitAndTrim(expectKeyContent, "\n") {
			t.Assert(gstr.Contains(content, line), true)
		}
	})
}

func Test_MetricCallback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx = gctx.New()
			_   = gmetric.NewCounter(gmetric.CounterConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.counter",
					Help: "This is a simple demo for Counter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_a", 1),
					},
					Instrument:        "github.com/gogf/gf/example/metric/basic",
					InstrumentVersion: "v1.4",
				},
				Callback: func(ctx context.Context) (*gmetric.CallbackResult, error) {
					return &gmetric.CallbackResult{
						Value: 100,
						Attributes: gmetric.Attributes{
							gmetric.NewAttribute("const_label_b", 2),
						},
					}, nil
				},
			})
			_ = gmetric.NewGauge(gmetric.GaugeConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.gauge",
					Help: "This is a simple demo for Gauge usage",
					Unit: "bytes",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_c", 3),
					},
					Instrument:        "github.com/gogf/gf/example/metric/basic",
					InstrumentVersion: "v1.5",
				},
				Callback: func(ctx context.Context) (*gmetric.CallbackResult, error) {
					return &gmetric.CallbackResult{
						Value: 101,
						Attributes: gmetric.Attributes{
							gmetric.NewAttribute("const_label_d", 4),
						},
					}, nil
				},
			})
		)
		reader := metric.NewManualReader()

		// OpenTelemetry provider.
		provider := otelmetric.MustProvider(metric.WithReader(reader))
		defer provider.Shutdown(ctx)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(ctx, &rm)
		t.AssertNil(err)

		content := gjson.MustEncodeString(rm)
		content, err = gregex.ReplaceString(`Time":".+?"`, `Time":""`, content)
		t.AssertNil(err)
		expectContent := `
{"Scope":{"Name":"github.com/gogf/gf/example/metric/basic","Version":"v1.4","SchemaURL":""},"Metrics":[{"Name":"goframe.metric.demo.counter","Description":"This is a simple demo for Counter usage","Unit":"%","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_a","Value":{"Type":"INT64","Value":1}},{"Key":"const_label_b","Value":{"Type":"INT64","Value":2}}],"StartTime":"","Time":"","Value":100}],"Temporality":"CumulativeTemporality","IsMonotonic":true}}]}
{"Name":"goframe.metric.demo.gauge","Description":"This is a simple demo for Gauge usage","Unit":"bytes","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_c","Value":{"Type":"INT64","Value":3}},{"Key":"const_label_d","Value":{"Type":"INT64","Value":4}}
`

		for _, line := range gstr.SplitAndTrim(expectContent, "\n") {
			t.Assert(gstr.Contains(content, line), true)
		}
	})
}

func Test_GlobalCallback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = gctx.New()
			counter = gmetric.NewCounter(gmetric.CounterConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.counter",
					Help: "This is a simple demo for Counter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_a", 1),
					},
					Instrument: "github.com/gogf/gf/example/metric/basic",
				},
			})
			gauge = gmetric.NewGauge(gmetric.GaugeConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.gauge",
					Help: "This is a simple demo for Gauge usage",
					Unit: "bytes",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_b", 2),
					},
					Instrument: "github.com/gogf/gf/example/metric/basic",
				},
			})
		)
		// global callback.
		gmetric.RegisterCallback(
			func(ctx context.Context, setter gmetric.CallbackSetter) error {
				setter.Set(counter, 100)
				setter.Set(gauge, 101)
				return nil
			},
			counter,
			gauge,
		)

		reader := metric.NewManualReader()
		// OpenTelemetry provider.
		provider := otelmetric.MustProvider(metric.WithReader(reader))
		defer provider.Shutdown(ctx)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(ctx, &rm)
		t.AssertNil(err)

		content := gjson.MustEncodeString(rm)
		content, err = gregex.ReplaceString(`Time":".+?"`, `Time":""`, content)
		t.AssertNil(err)
		expectContent := `
{"Name":"goframe.metric.demo.counter","Description":"This is a simple demo for Counter usage","Unit":"%","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_a","Value":{"Type":"INT64","Value":1}}],"StartTime":"","Time":"","Value":100}],"Temporality":"CumulativeTemporality","IsMonotonic":true}}
{"Name":"goframe.metric.demo.gauge","Description":"This is a simple demo for Gauge usage","Unit":"bytes","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_b","Value":{"Type":"INT64","Value":2}}],"StartTime":"","Time":"","Value":101}]}}
`
		for _, line := range gstr.SplitAndTrim(expectContent, "\n") {
			t.Assert(gstr.Contains(content, line), true)
		}
	})
}

func Test_GlobalCallback_DynamicAttributes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = gctx.New()
			counter = gmetric.NewCounter(gmetric.CounterConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.counter",
					Help: "This is a simple demo for Counter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_a", 1),
					},
					Instrument: "github.com/gogf/gf/example/metric/basic",
				},
			})
		)
		// global callback.
		gmetric.RegisterCallback(
			func(ctx context.Context, setter gmetric.CallbackSetter) error {
				setter.Set(counter, 1000, gmetric.Option{
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("dynamic_label_b", 2),
					}},
				)
				return nil
			},
			counter,
		)

		reader := metric.NewManualReader()
		// OpenTelemetry provider.
		provider := otelmetric.MustProvider(metric.WithReader(reader))
		defer provider.Shutdown(ctx)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(ctx, &rm)
		t.AssertNil(err)

		content := gjson.MustEncodeString(rm)
		content, err = gregex.ReplaceString(`Time":".+?"`, `Time":""`, content)
		t.AssertNil(err)
		expectContent := `
{"Scope":{"Name":"github.com/gogf/gf/example/metric/basic","Version":"","SchemaURL":""},"Metrics":[{"Name":"goframe.metric.demo.counter","Description":"This is a simple demo for Counter usage","Unit":"%","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_a","Value":{"Type":"INT64","Value":1}},{"Key":"dynamic_label_b","Value":{"Type":"INT64","Value":2}}],"StartTime":"","Time":"","Value":1000}],"Temporality":"CumulativeTemporality","IsMonotonic":true}}]}
`
		for _, line := range gstr.SplitAndTrim(expectContent, "\n") {
			t.Assert(gstr.Contains(content, line), true)
		}
	})
}

func Test_GlobalCallback_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			counter = gmetric.NewCounter(gmetric.CounterConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.counter",
					Help: "This is a simple demo for Counter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_a", 1),
					},
					Instrument:        "github.com/gogf/gf/example/metric/basic",
					InstrumentVersion: "v1.4",
				},
			})
			gauge = gmetric.NewGauge(gmetric.GaugeConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.gauge",
					Help: "This is a simple demo for Gauge usage",
					Unit: "bytes",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_c", 3),
					},
					Instrument:        "github.com/gogf/gf/example/metric/basic",
					InstrumentVersion: "v1.5",
				},
			})
		)
		// global callback.
		gmetric.RegisterCallback(
			func(ctx context.Context, setter gmetric.CallbackSetter) error {
				setter.Set(counter, 100)
				setter.Set(gauge, 101)
				return nil
			},
			counter,
			gauge,
		)

		reader := metric.NewManualReader()
		// OpenTelemetry provider.
		_, err := otelmetric.NewProvider(metric.WithReader(reader))
		t.Assert(gstr.Contains(
			err.Error(),
			`multiple instrument or instrument version metrics used in the same callback`,
		), true)
	})
}
