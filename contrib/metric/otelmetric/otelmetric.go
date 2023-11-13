// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package otelmetric provides metric functionalities using OpenTelemetry metric.
package otelmetric

import (
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/gogf/gf/v2/os/gmetric"
)

// NewProvider creates and returns a metrics provider.
func NewProvider(option ...metric.Option) (gmetric.Provider, error) {
	return newLocalProvider(option...)
}

// MustProvider creates and returns a metrics provider.
// It panics if any error occurs.
func MustProvider(option ...metric.Option) gmetric.Provider {
	provider, err := newLocalProvider(option...)
	if err != nil {
		panic(err)
	}
	return provider
}
