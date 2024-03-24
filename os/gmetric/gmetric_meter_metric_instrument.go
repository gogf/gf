// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// localMetricInstrument implements interface MetricInstrument.
type localInstrumentInfo struct {
	name    string
	version string
}

// newInstrumentInfo creates and returns a MetricInstrument.
func (meter *localMeter) newInstrumentInfo() InstrumentInfo {
	return &localInstrumentInfo{
		name:    meter.Instrument,
		version: meter.InstrumentVersion,
	}
}

// Name returns the instrument name of the metric.
func (l *localInstrumentInfo) Name() string {
	return l.name
}

// Version returns the instrument version of the metric.
func (l *localInstrumentInfo) Version() string {
	return l.version
}
