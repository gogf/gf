// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric_test

import (
	"context"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

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
			counter = gmetric.MustNewCounter(gmetric.MetricConfig{
				Name: "goframe.metric.demo.counter",
				Help: "This is a simple demo for Counter usage",
				Unit: "%",
				Attributes: gmetric.Attributes{
					gmetric.NewAttribute("const_label_a", 1),
				},
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.1",
			})

			histogram = gmetric.MustNewHistogram(gmetric.MetricConfig{
				Name: "goframe.metric.demo.histogram",
				Help: "This is a simple demo for histogram usage",
				Unit: "ms",
				Attributes: gmetric.Attributes{
					gmetric.NewAttribute("const_label_b", 3),
				},
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.2",
				Buckets:           []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
			})

			observableCounter = gmetric.MustNewObservableCounter(gmetric.MetricConfig{
				Name: "goframe.metric.demo.observable_counter",
				Help: "This is a simple demo for ObservableCounter usage",
				Unit: "%",
				Attributes: gmetric.Attributes{
					gmetric.NewAttribute("const_label_c", 1),
				},
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.3",
			})

			observableGauge = gmetric.MustNewObservableGauge(gmetric.MetricConfig{
				Name: "goframe.metric.demo.observable_gauge",
				Help: "This is a simple demo for ObservableGauge usage",
				Unit: "%",
				Attributes: gmetric.Attributes{
					gmetric.NewAttribute("const_label_d", 1),
				},
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.3",
			})
		)

		gmetric.MustRegisterCallback(func(ctx context.Context, obs gmetric.Observer) error {
			obs.Observe(observableCounter, 10)
			obs.Observe(observableGauge, 10)
			return nil
		}, observableCounter, observableGauge)

		reader := metric.NewManualReader()

		// OpenTelemetry provider.
		provider := otelmetric.MustProvider(metric.WithReader(reader))
		defer provider.Shutdown(ctx)

		// Add value for counter.
		counter.Inc(ctx)
		counter.Add(ctx, 10)
		counter.Dec(ctx, gmetric.Option{
			Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_a", "1")},
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
			Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_d", "4")},
		}
		histogram.Record(100, histogramOption)
		histogram.Record(200, histogramOption)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(ctx, &rm)
		t.AssertNil(err)

		//var (
		//	sm1 = instrumentation.Scope{
		//		Name:      "github.com/gogf/gf/example/metric/basic",
		//		Version:   "v1.1",
		//		SchemaURL: "",
		//	}
		//	sm2 = instrumentation.Scope{
		//		Name:      "github.com/gogf/gf/example/metric/basic",
		//		Version:   "v1.1",
		//		SchemaURL: "",
		//	}
		//	sm3 = instrumentation.Scope{
		//		Name:      "github.com/gogf/gf/example/metric/basic",
		//		Version:   "v1.1",
		//		SchemaURL: "",
		//	}
		//)
		//t.Assert(len(rm.ScopeMetrics), 4)
		//for _, sm := range rm.ScopeMetrics {
		//	switch sm.Scope {
		//
		//	}
		//}
		//t.Assert(rm.ScopeMetrics[0].Scope.Version, "v1.2")
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
			ctx     = gctx.New()
			counter = gmetric.MustNewCounter(gmetric.MetricConfig{
				Name: "goframe.metric.demo.counter",
				Help: "This is a simple demo for Counter usage",
				Unit: "%",
				Attributes: gmetric.Attributes{
					gmetric.NewAttribute("const_label_a", 1),
				},
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.1",
			})

			histogram = gmetric.MustNewHistogram(gmetric.MetricConfig{
				Name: "goframe.metric.demo.histogram",
				Help: "This is a simple demo for histogram usage",
				Unit: "ms",
				Attributes: gmetric.Attributes{
					gmetric.NewAttribute("const_label_b", 3),
				},
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.2",
				Buckets:           []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
			})

			observableCounter = gmetric.MustNewObservableCounter(gmetric.MetricConfig{
				Name: "goframe.metric.demo.observable_counter",
				Help: "This is a simple demo for ObservableCounter usage",
				Unit: "%",
				Attributes: gmetric.Attributes{
					gmetric.NewAttribute("const_label_c", 1),
				},
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.3",
			})

			observableGauge = gmetric.MustNewObservableGauge(gmetric.MetricConfig{
				Name: "goframe.metric.demo.observable_gauge",
				Help: "This is a simple demo for ObservableGauge usage",
				Unit: "%",
				Attributes: gmetric.Attributes{
					gmetric.NewAttribute("const_label_d", 1),
				},
				Instrument:        "github.com/gogf/gf/example/metric/basic",
				InstrumentVersion: "v1.3",
			})
		)

		gmetric.MustRegisterCallback(func(ctx context.Context, obs gmetric.Observer) error {
			obs.Observe(observableCounter, 10)
			obs.Observe(observableGauge, 10)
			return nil
		}, observableCounter, observableGauge)

		reader := metric.NewManualReader()

		// OpenTelemetry provider.
		provider := otelmetric.MustProvider(metric.WithReader(reader))
		defer provider.Shutdown(ctx)

		// Add value for counter.
		counter.Inc(ctx)
		counter.Add(ctx, 10)
		counter.Dec(ctx, gmetric.Option{
			Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_a", "1")},
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
			Attributes: gmetric.Attributes{gmetric.NewAttribute("dynamic_label_d", "4")},
		}
		histogram.Record(100, histogramOption)
		histogram.Record(200, histogramOption)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(ctx, &rm)
		t.AssertNil(err)

		content := gjson.MustEncodeString(rm)
		content, err = gregex.ReplaceString(`Time":".+?"`, `Time":""`, content)
		t.AssertNil(err)
		expectContent := `
{"Name":"goframe.metric.demo.counter","Description":"This is a simple demo for Counter usage","Unit":"%","Data":{"DataPoints":[{"Attributes":[{"Key":"const_label_a","Value":{"Type":"INT64","Value":1}},{"Key":"dynamic_label_b","Value":{"Type":"INT64","Value":2}},{"Key":"g1","Value":{"Type":"INT64","Value":1}}],"StartTime":"","Time":"","Value":1000}],"Temporality":"CumulativeTemporality","IsMonotonic":true}}
`
		fmt.Println(content)
		for _, line := range gstr.SplitAndTrim(expectContent, "\n") {
			fmt.Println(line)
			t.Assert(gstr.Contains(content, line), true)
		}
	})
}
