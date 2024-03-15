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
func (l *localPerformer) Counter(config gmetric.MetricConfig) (gmetric.CounterPerformer, error) {
	var meter = l.createMeter(config.Instrument, config.InstrumentVersion)
	return newCounterPerformer(meter, config)
}

// UpDownCounter creates and returns a UpDownCounterPerformer that performs
// the operations for UpDownCounter metric.
func (l *localPerformer) UpDownCounter(config gmetric.MetricConfig) (gmetric.UpDownCounterPerformer, error) {
	var meter = l.createMeter(config.Instrument, config.InstrumentVersion)
	return newUpDownCounterPerformer(meter, config)
}

// Histogram creates and returns a HistogramPerformer that performs
// the operations for Histogram metric.
func (l *localPerformer) Histogram(config gmetric.MetricConfig) (gmetric.HistogramPerformer, error) {
	var meter = l.createMeter(config.Instrument, config.InstrumentVersion)
	return newHistogramPerformer(meter, config)
}

// ObservableCounter creates and returns an ObservableMetric that performs
// the operations for ObservableCounter metric.
func (l *localPerformer) ObservableCounter(config gmetric.MetricConfig) (gmetric.ObservableMetric, error) {
	var meter = l.createMeter(config.Instrument, config.InstrumentVersion)
	return newObservableCounterPerformer(meter, config)
}

// ObservableUpDownCounter creates and returns an ObservableMetric that performs
// the operations for ObservableUpDownCounter metric.
func (l *localPerformer) ObservableUpDownCounter(config gmetric.MetricConfig) (gmetric.ObservableMetric, error) {
	var meter = l.createMeter(config.Instrument, config.InstrumentVersion)
	return newObservableUpDownCounterPerformer(meter, config)
}

// ObservableGauge creates and returns an ObservableMetric that performs
// the operations for ObservableGauge metric.
func (l *localPerformer) ObservableGauge(config gmetric.MetricConfig) (gmetric.ObservableMetric, error) {
	var meter = l.createMeter(config.Instrument, config.InstrumentVersion)
	return newObservableGaugePerformer(meter, config)
}

// createMeter creates and returns an OpenTelemetry Meter.
func (l *localPerformer) createMeter(instrument, instrumentVersion string) otelmetric.Meter {
	return l.Meter(
		instrument,
		otelmetric.WithInstrumentationVersion(instrumentVersion),
	)
}
