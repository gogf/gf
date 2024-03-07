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

// localPerformer implements interface gmetric.Performer.
type localPerformer struct {
	*metric.MeterProvider
}

// newPerformer creates and returns gmetric.Meter.
func newPerformer(provider *metric.MeterProvider) gmetric.Performer {
	performer := &localPerformer{
		MeterProvider: provider,
	}
	return performer
}

// Counter creates and returns a CounterPerformer that performs
// the operations for Counter metric.
func (l *localPerformer) Counter(config gmetric.CounterConfig) (gmetric.CounterPerformer, error) {
	var meter = l.createMeter(config.Instrument, config.InstrumentVersion)
	return newCounterPerformer(meter, config)
}

// Gauge creates and returns a GaugePerformer that performs
// the operations for Gauge metric.
func (l *localPerformer) Gauge(config gmetric.GaugeConfig) (gmetric.GaugePerformer, error) {
	var meter = l.createMeter(config.Instrument, config.InstrumentVersion)
	return newGaugePerformer(meter, config)
}

// Histogram creates and returns a HistogramPerformer that performs
// the operations for Histogram metric.
func (l *localPerformer) Histogram(config gmetric.HistogramConfig) (gmetric.HistogramPerformer, error) {
	var meter = l.createMeter(config.Instrument, config.InstrumentVersion)
	return newHistogramPerformer(meter, config)
}

// createMeter creates and returns an OpenTelemetry Meter.
func (l *localPerformer) createMeter(instrument, instrumentVersion string) otelmetric.Meter {
	return l.Meter(
		instrument,
		otelmetric.WithInstrumentationVersion(instrumentVersion),
	)
}
