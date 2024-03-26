// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localMeter for Meter implements.
type localMeter struct {
	MeterOption
}

// MeterOption holds the creation option for a Meter.
type MeterOption struct {
	// Instrument is the instrumentation name to bind this Metric to a global MeterProvider.
	// This is an optional configuration for a metric.
	Instrument string

	// InstrumentVersion is the instrumentation version to bind this Metric to a global MeterProvider.
	// This is an optional configuration for a metric.
	InstrumentVersion string

	// Attributes holds the constant key-value pair description metadata for all metrics of Meter.
	// This is an optional configuration for a meter.
	Attributes Attributes
}

// newMeter creates and returns a Meter implementer.
func newMeter(option MeterOption) Meter {
	return &localMeter{
		MeterOption: option,
	}
}

// Performer creates and returns the Performer of the Meter.
func (meter *localMeter) Performer() MeterPerformer {
	if globalProvider == nil {
		return nil
	}
	return globalProvider.MeterPerformer(meter.MeterOption)
}
