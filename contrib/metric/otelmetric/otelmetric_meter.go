// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
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
	l.provider.Meter(config.Instrument).
}

func (l *localMeter) GaugePerformer(config gmetric.GaugeConfig) gmetric.GaugePerformer {

}

func (l *localMeter) HistogramPerformer(config gmetric.HistogramConfig) gmetric.HistogramPerformer {

}
