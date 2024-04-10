// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx      = gctx.New()
			meterV11 = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.1",
			})
			meterV12 = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.2",
			})
			meterV13 = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.3",
			})
			meterV14 = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.4",
			})
			counter = meterV11.MustCounter(
				"goframe.metric.demo.counter",
				gmetric.MetricOption{
					Help: "This is a simple demo for Counter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_1", 1),
					},
				},
			)
			upDownCounter = meterV12.MustUpDownCounter(
				"goframe.metric.demo.updown_counter",
				gmetric.MetricOption{
					Help: "This is a simple demo for UpDownCounter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_2", 2),
					},
				},
			)
			histogram = meterV13.MustHistogram(
				"goframe.metric.demo.histogram",
				gmetric.MetricOption{
					Help: "This is a simple demo for histogram usage",
					Unit: "ms",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_3", 3),
					},
					Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
				},
			)
			observableCounter = meterV14.MustObservableCounter(
				"goframe.metric.demo.observable_counter",
				gmetric.MetricOption{
					Help: "This is a simple demo for ObservableCounter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_4", 4),
					},
				},
			)
			observableUpDownCounter = meterV14.MustObservableUpDownCounter(
				"goframe.metric.demo.observable_updown_counter",
				gmetric.MetricOption{
					Help: "This is a simple demo for ObservableUpDownCounter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_5", 5),
					},
				},
			)
			observableGauge = meterV14.MustObservableGauge(
				"goframe.metric.demo.observable_gauge",
				gmetric.MetricOption{
					Help: "This is a simple demo for ObservableGauge usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_6", 6),
					},
				},
			)
		)

		meterV14.MustRegisterCallback(func(ctx context.Context, obs gmetric.Observer) error {
			obs.Observe(observableCounter, 10, gmetric.Option{
				Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_4", "4")},
			})
			obs.Observe(observableUpDownCounter, 20, gmetric.Option{
				Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_5", "5")},
			})
			obs.Observe(observableGauge, 30, gmetric.Option{
				Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_6", "6")},
			})
			return nil
		}, observableCounter, observableUpDownCounter, observableGauge)

		reader := metric.NewManualReader()

		// OpenTelemetry provider.
		provider := otelmetric.MustProvider(otelmetric.WithReader(reader))
		defer provider.Shutdown(ctx)

		// Counter.
		counter.Inc(ctx)
		counter.Add(ctx, 10, gmetric.Option{
			Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_1", "1")},
		})

		upDownCounter.Add(ctx, 10)
		upDownCounter.Dec(ctx, gmetric.Option{
			Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_2", "2")},
		})

		// Record values for histogram.
		histogram.Record(1)
		histogram.Record(20)
		histogram.Record(30)
		histogram.Record(101)
		histogram.Record(2000)
		histogram.Record(9000)
		histogram.Record(20000)

		histogramOption := gmetric.Option{
			Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_3", "3")},
		}
		histogram.Record(100, histogramOption)
		histogram.Record(200, histogramOption)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(ctx, &rm)
		t.AssertNil(err)

		metricsJsonContent := gjson.MustEncodeString(rm)

		t.Assert(len(rm.ScopeMetrics), 4)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.counter`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.updown_counter`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.histogram`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.observable_counter`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.observable_updown_counter"`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.observable_gauge`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_2","Value":{"Type":"INT64","Value":2}}`), 2)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"dynamic_label_2","Value":{"Type":"STRING","Value":"2"}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_3","Value":{"Type":"INT64","Value":3}}`), 2)
		t.Assert(gstr.Count(metricsJsonContent, `"Count":7,"Bounds":[0,10,20,50,100,500,1000,2000,5000,10000],"BucketCounts":[0,1,1,1,0,1,0,1,0,1,1],"Min":1,"Max":20000,"Sum":31152`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_3","Value":{"Type":"INT64","Value":3}}`), 2)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"dynamic_label_3","Value":{"Type":"STRING","Value":"3"}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `"Count":2,"Bounds":[0,10,20,50,100,500,1000,2000,5000,10000],"BucketCounts":[0,0,0,0,1,1,0,0,0,0,0],"Min":100,"Max":200,"Sum":300`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_4","Value":{"Type":"INT64","Value":4}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"dynamic_label_4","Value":{"Type":"STRING","Value":"4"}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_5","Value":{"Type":"INT64","Value":5}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"dynamic_label_5","Value":{"Type":"STRING","Value":"5"}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_6","Value":{"Type":"INT64","Value":6}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"dynamic_label_6","Value":{"Type":"STRING","Value":"6"}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_1","Value":{"Type":"INT64","Value":1}}`), 2)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"dynamic_label_1","Value":{"Type":"STRING","Value":"1"}}`), 1)
	})
}

func Test_GlobalAttributes(t *testing.T) {
	gmetric.SetGlobalAttributes(gmetric.Attributes{
		gmetric.NewAttribute("g1", 1),
	}, gmetric.SetGlobalAttributesOption{
		Instrument:        "github.com/gogf/gf/example/metric/basic",
		InstrumentVersion: "v1.1",
		InstrumentPattern: "",
	})
	gmetric.SetGlobalAttributes(gmetric.Attributes{
		gmetric.NewAttribute("g2", 2),
	}, gmetric.SetGlobalAttributesOption{
		Instrument:        "github.com/gogf/gf/example/metric/basic",
		InstrumentVersion: "v1.3",
		InstrumentPattern: "",
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx      = gctx.New()
			meterV11 = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.1",
			})
			meterV12 = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.2",
			})
			meterV13 = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.3",
			})
			counter = meterV11.MustCounter(
				"goframe.metric.demo.counter",
				gmetric.MetricOption{
					Help: "This is a simple demo for Counter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_1", 1),
					},
				},
			)

			histogram = meterV12.MustHistogram(
				"goframe.metric.demo.histogram",
				gmetric.MetricOption{
					Help: "This is a simple demo for histogram usage",
					Unit: "ms",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_2", 2),
					},
					Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
				},
			)

			observableCounter = meterV13.MustObservableCounter(
				"goframe.metric.demo.observable_counter",
				gmetric.MetricOption{
					Help: "This is a simple demo for ObservableCounter usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_3", 3),
					},
				},
			)

			observableGauge = meterV13.MustObservableGauge(
				"goframe.metric.demo.observable_gauge",
				gmetric.MetricOption{
					Help: "This is a simple demo for ObservableGauge usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_4", 4),
					},
				},
			)
		)

		meterV13.MustRegisterCallback(func(ctx context.Context, obs gmetric.Observer) error {
			obs.Observe(observableCounter, 10, gmetric.Option{
				Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_3", "3")},
			})
			obs.Observe(observableGauge, 10, gmetric.Option{
				Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_4", "4")},
			})
			return nil
		}, observableCounter, observableGauge)

		reader := metric.NewManualReader()

		// OpenTelemetry provider.
		provider := otelmetric.MustProvider(otelmetric.WithReader(reader))
		defer provider.Shutdown(ctx)

		// Add value for counter.
		counter.Inc(ctx)
		counter.Add(ctx, 10, gmetric.Option{
			Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_1", "1")},
		})

		// Record values for histogram.
		histogram.Record(1)
		histogram.Record(20)
		histogram.Record(30)
		histogram.Record(101)
		histogram.Record(2000)
		histogram.Record(9000)
		histogram.Record(20000)

		histogramOption := gmetric.Option{
			Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_2", "2")},
		}
		histogram.Record(100, histogramOption)
		histogram.Record(200, histogramOption)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(ctx, &rm)
		t.AssertNil(err)

		metricsJsonContent := gjson.MustEncodeString(rm)
		t.Assert(len(rm.ScopeMetrics), 3)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.counter`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.histogram`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.observable_counter`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.observable_gauge`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `goframe.metric.demo.observable_gauge`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_1","Value":{"Type":"INT64","Value":1}}`), 2)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"g1","Value":{"Type":"INT64","Value":1}}`), 2)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"dynamic_label_1","Value":{"Type":"STRING","Value":"1"}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_2","Value":{"Type":"INT64","Value":2}}`), 2)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"dynamic_label_2","Value":{"Type":"STRING","Value":"2"}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `"Count":2,"Bounds":[0,10,20,50,100,500,1000,2000,5000,10000],"BucketCounts":[0,0,0,0,1,1,0,0,0,0,0],"Min":100,"Max":200,"Sum":300`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `"Count":7,"Bounds":[0,10,20,50,100,500,1000,2000,5000,10000],"BucketCounts":[0,1,1,1,0,1,0,1,0,1,1],"Min":1,"Max":20000,"Sum":31152`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_3","Value":{"Type":"INT64","Value":3}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"dynamic_label_3","Value":{"Type":"STRING","Value":"3"}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"g2","Value":{"Type":"INT64","Value":2}}`), 2)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"const_label_4","Value":{"Type":"INT64","Value":4}}`), 1)
		t.Assert(gstr.Count(metricsJsonContent, `{"Key":"dynamic_label_4","Value":{"Type":"STRING","Value":"4"}}`), 1)
	})
}
