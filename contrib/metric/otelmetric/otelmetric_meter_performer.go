// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"context"
	"fmt"
	"reflect"

	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gmetric"
)

// localMeterPerformer implements interface gmetric.Performer.
type localMeterPerformer struct {
	gmetric.MeterOption
	*metric.MeterProvider
}

// newMeterPerformer creates and returns gmetric.Meter.
func newMeterPerformer(provider *metric.MeterProvider, option gmetric.MeterOption) gmetric.MeterPerformer {
	meterPerformer := &localMeterPerformer{
		MeterOption:   option,
		MeterProvider: provider,
	}
	return meterPerformer
}

// CounterPerformer creates and returns a CounterPerformer that performs
// the operations for Counter metric.
func (l *localMeterPerformer) CounterPerformer(name string, option gmetric.MetricOption) (gmetric.CounterPerformer, error) {
	return l.newCounterPerformer(l.createMeter(), name, option)
}

// UpDownCounterPerformer creates and returns a UpDownCounterPerformer that performs
// the operations for UpDownCounter metric.
func (l *localMeterPerformer) UpDownCounterPerformer(name string, option gmetric.MetricOption) (gmetric.UpDownCounterPerformer, error) {
	return l.newUpDownCounterPerformer(l.createMeter(), name, option)
}

// HistogramPerformer creates and returns a HistogramPerformer that performs
// the operations for Histogram metric.
func (l *localMeterPerformer) HistogramPerformer(name string, option gmetric.MetricOption) (gmetric.HistogramPerformer, error) {
	return l.newHistogramPerformer(l.createMeter(), name, option)
}

// ObservableCounterPerformer creates and returns an ObservableMetric that performs
// the operations for ObservableCounter metric.
func (l *localMeterPerformer) ObservableCounterPerformer(name string, option gmetric.MetricOption) (gmetric.ObservableMetric, error) {
	return l.newObservableCounterPerformer(l.createMeter(), name, option)
}

// ObservableUpDownCounterPerformer creates and returns an ObservableMetric that performs
// the operations for ObservableUpDownCounter metric.
func (l *localMeterPerformer) ObservableUpDownCounterPerformer(name string, option gmetric.MetricOption) (gmetric.ObservableMetric, error) {
	return l.newObservableUpDownCounterPerformer(l.createMeter(), name, option)
}

// ObservableGaugePerformer creates and returns an ObservableMetric that performs
// the operations for ObservableGauge metric.
func (l *localMeterPerformer) ObservableGaugePerformer(name string, option gmetric.MetricOption) (gmetric.ObservableMetric, error) {
	return l.newObservableGaugePerformer(l.createMeter(), name, option)
}

// RegisterCallback registers callback on certain metrics.
// A callback is bound to certain component and version, it is called when the associated metrics are read.
// Multiple callbacks on the same component and version will be called by their registered sequence.
func (l *localMeterPerformer) RegisterCallback(
	callback gmetric.Callback, observableMetrics ...gmetric.ObservableMetric,
) error {
	var metrics = make([]gmetric.Metric, 0)
	for _, v := range observableMetrics {
		m, ok := v.(gmetric.Metric)
		if !ok {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid metric parameter "%s" for RegisterCallback, which does not implement interface Metric`,
				reflect.TypeOf(v).String(),
			)
		}
		metrics = append(metrics, m)
	}
	// group the metric by instrument and instrument version.
	var (
		instrumentSet      = gset.NewStrSet()
		underlyingMeterMap = map[otelmetric.Meter][]otelmetric.Observable{}
	)
	for _, m := range metrics {
		var meter = l.Meter(
			m.Info().Instrument().Name(),
			otelmetric.WithInstrumentationVersion(m.Info().Instrument().Version()),
		)
		instrumentSet.Add(fmt.Sprintf(
			`%s@%s`,
			m.Info().Instrument().Name(),
			m.Info().Instrument().Version(),
		))
		if _, ok := underlyingMeterMap[meter]; !ok {
			underlyingMeterMap[meter] = make([]otelmetric.Observable, 0)
		}
		underlyingMeterMap[meter] = append(underlyingMeterMap[meter], metricToFloat64Observable(m))
	}
	if len(underlyingMeterMap) > 1 {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`multiple instrument or instrument version metrics used in the same callback: %s`,
			instrumentSet.Join(","),
		)
	}
	// do callback registering.
	for meter, observables := range underlyingMeterMap {
		_, err := meter.RegisterCallback(
			func(ctx context.Context, observer otelmetric.Observer) error {
				return callback(ctx, newObserver(observer, l.MeterOption))
			},
			observables...,
		)
		if err != nil {
			return gerror.WrapCode(
				gcode.CodeInternalError, err,
				`RegisterCallback failed`,
			)
		}
	}
	return nil
}

// createMeter creates and returns an OpenTelemetry Meter.
func (l *localMeterPerformer) createMeter() otelmetric.Meter {
	return l.Meter(
		l.Instrument,
		otelmetric.WithInstrumentationVersion(l.InstrumentVersion),
	)
}
