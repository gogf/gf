// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/gogf/gf/v2/os/gmetric"
)

type localMeter struct {
	provider *metric.MeterProvider
}

func newMeter(provider *metric.MeterProvider) gmetric.Meter {
	meter := &localMeter{
		provider: provider,
	}
	return meter
}

func (l *localMeter) CounterPerformer(config gmetric.CounterConfig) gmetric.CounterPerformer {
	var (
		meter = l.provider.Meter(
			config.Instrument,
			otelmetric.WithInstrumentationVersion(config.InstrumentVersion),
		)
		performer = newCounterPerformer(meter, config)
	)
	return performer
}

func (l *localMeter) GaugePerformer(config gmetric.GaugeConfig) gmetric.GaugePerformer {
	var (
		meter     = l.provider.Meter(config.Instrument)
		performer = newGaugePerformer(meter, config)
	)
	return performer
}

func (l *localMeter) HistogramPerformer(config gmetric.HistogramConfig) gmetric.HistogramPerformer {
	var (
		meter     = l.provider.Meter(config.Instrument)
		performer = newHistogramPerformer(meter, config)
	)
	return performer
}
