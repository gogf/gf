// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"go.opentelemetry.io/otel/metric"

	"github.com/gogf/gf/v2/os/gmetric"
)

type localCallbackSetter struct {
	observer metric.Observer
}

func newCallbackSetter(observer metric.Observer) gmetric.CallbackSetter {
	return &localCallbackSetter{
		observer: observer,
	}
}

func (l *localCallbackSetter) Set(m gmetric.Metric, value float64, option ...gmetric.Option) {
	l.observer.ObserveFloat64(metricToFloat64Observable(m), value)
}
