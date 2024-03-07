// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Counter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			config = gmetric.CounterConfig{
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
			}
			counter = gmetric.MustNewCounter(config)
		)
		t.Assert(counter.Info().Name(), config.Name)
		t.Assert(counter.Info().Help(), config.Help)
		t.Assert(counter.Info().Unit(), config.Unit)
		t.Assert(counter.Info().Attributes(), config.Attributes)
		t.Assert(counter.Info().Instrument().Name(), config.Instrument)
		t.Assert(counter.Info().Instrument().Version(), config.InstrumentVersion)
	})
}

func Test_Gauge(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			config = gmetric.GaugeConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.gauge",
					Help: "This is a simple demo for Gauge usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_a", 1),
					},
					Instrument:        "github.com/gogf/gf/example/metric/basic",
					InstrumentVersion: "v1.0",
				},
			}
			counter = gmetric.MustNewGauge(config)
		)
		t.Assert(counter.Info().Name(), config.Name)
		t.Assert(counter.Info().Help(), config.Help)
		t.Assert(counter.Info().Unit(), config.Unit)
		t.Assert(counter.Info().Attributes(), config.Attributes)
		t.Assert(counter.Info().Instrument().Name(), config.Instrument)
		t.Assert(counter.Info().Instrument().Version(), config.InstrumentVersion)
	})
}

func Test_Histogram(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			config = gmetric.HistogramConfig{
				MetricConfig: gmetric.MetricConfig{
					Name: "goframe.metric.demo.hist",
					Help: "This is a simple demo for Gauge usage",
					Unit: "%",
					Attributes: gmetric.Attributes{
						gmetric.NewAttribute("const_label_a", 1),
					},
					Instrument:        "github.com/gogf/gf/example/metric/basic",
					InstrumentVersion: "v1.0",
				},
				Buckets: []float64{0, 10, 20, 50, 100, 500, 1000, 2000, 5000, 10000},
			}
			counter = gmetric.MustNewHistogram(config)
		)
		t.Assert(counter.Info().Name(), config.Name)
		t.Assert(counter.Info().Help(), config.Help)
		t.Assert(counter.Info().Unit(), config.Unit)
		t.Assert(counter.Info().Attributes(), config.Attributes)
		t.Assert(counter.Info().Instrument().Name(), config.Instrument)
		t.Assert(counter.Info().Instrument().Version(), config.InstrumentVersion)
		t.Assert(counter.Buckets(), config.Buckets)
	})
}

func Test_CommonAttributes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		commonAttributes := gmetric.CommonAttributes()
		t.AssertGT(len(commonAttributes), 1)
		fmt.Println(commonAttributes)
	})
}
