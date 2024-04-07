// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// newProviderConfigByOptions returns a config configured with options.
func newProviderConfigByOptions(options []Option) providerConfig {
	conf := providerConfig{}
	for _, o := range options {
		conf = o.apply(conf)
	}
	return conf
}

// Option applies a configuration option value to a MeterProvider.
type Option interface {
	apply(providerConfig) providerConfig
}

// optionFunc applies a set of options to a config.
type optionFunc func(providerConfig) providerConfig

// apply returns a config with option(s) applied.
func (o optionFunc) apply(conf providerConfig) providerConfig {
	return o(conf)
}

// providerConfig is the configuration for Provider.
type providerConfig struct {
	viewOption            metric.Option
	readerOption          metric.Option
	resourceOption        metric.Option
	enabledBuiltInMetrics bool
}

// IsBuiltInMetricsEnabled returns whether the builtin metrics is enabled.
func (cfg providerConfig) IsBuiltInMetricsEnabled() bool {
	return cfg.enabledBuiltInMetrics
}

// MetricOptions converts and returns the providerConfig as metrics options.
func (cfg providerConfig) MetricOptions() []metric.Option {
	var metricOptions = make([]metric.Option, 0)
	if cfg.viewOption != nil {
		metricOptions = append(metricOptions, cfg.viewOption)
	}
	if cfg.readerOption != nil {
		metricOptions = append(metricOptions, cfg.readerOption)
	}
	if cfg.resourceOption != nil {
		metricOptions = append(metricOptions, cfg.resourceOption)
	}
	return metricOptions
}

// WithBuiltInMetrics enables builtin metrics.
func WithBuiltInMetrics() Option {
	return optionFunc(func(cfg providerConfig) providerConfig {
		cfg.enabledBuiltInMetrics = true
		return cfg
	})
}

// WithResource associates a Resource with a MeterProvider. This Resource
// represents the entity producing telemetry and is associated with all Meters
// the MeterProvider will create.
func WithResource(res *resource.Resource) Option {
	return optionFunc(func(cfg providerConfig) providerConfig {
		cfg.resourceOption = metric.WithResource(res)
		return cfg
	})
}

// WithReader associates Reader r with a MeterProvider.
//
// By default, if this option is not used, the MeterProvider will perform no
// operations; no data will be exported without a Reader.
func WithReader(reader metric.Reader) Option {
	return optionFunc(func(cfg providerConfig) providerConfig {
		if reader == nil {
			return cfg
		}
		cfg.readerOption = metric.WithReader(reader)
		return cfg
	})
}

// WithView associates views a MeterProvider.
//
// Views are appended to existing ones in a MeterProvider if this option is
// used multiple times.
//
// By default, if this option is not used, the MeterProvider will use the
// default view.
func WithView(views ...metric.View) Option {
	return optionFunc(func(cfg providerConfig) providerConfig {
		cfg.viewOption = metric.WithView(views...)
		return cfg
	})
}
