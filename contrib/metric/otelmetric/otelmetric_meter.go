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

// localMeter implements interface gmetric.Meter.
type localMeter struct {
	provider *metric.MeterProvider
}

// newMeter creates and returns gmetric.Meter.
func newMeter(provider *metric.MeterProvider) gmetric.Meter {
	meter := &localMeter{
		provider: provider,
	}
	return meter
}

// CounterPerformer creates and returns a CounterPerformer that performs
// the operations for Counter metric.
func (l *localMeter) CounterPerformer(config gmetric.CounterConfig) gmetric.CounterPerformer {
	var (
		meter     = l.createMeter(config.Instrument, config.InstrumentVersion)
		performer = newCounterPerformer(meter, config)
	)
	return performer
}

// GaugePerformer creates and returns a GaugePerformer that performs
// the operations for Gauge metric.
func (l *localMeter) GaugePerformer(config gmetric.GaugeConfig) gmetric.GaugePerformer {
	var (
		meter     = l.createMeter(config.Instrument, config.InstrumentVersion)
		performer = newGaugePerformer(meter, config)
	)
	return performer
}

// HistogramPerformer creates and returns a HistogramPerformer that performs
// the operations for Histogram metric.
func (l *localMeter) HistogramPerformer(config gmetric.HistogramConfig) gmetric.HistogramPerformer {
	var (
		meter     = l.createMeter(config.Instrument, config.InstrumentVersion)
		performer = newHistogramPerformer(meter, config)
	)
	return performer
}

// createMeter creates and returns an OpenTelemetry Meter.
func (l *localMeter) createMeter(instrument, instrumentVersion string) otelmetric.Meter {
	return l.provider.Meter(
		instrument,
		otelmetric.WithInstrumentationVersion(instrumentVersion),
	)
}
