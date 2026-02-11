// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package otel provides OpenTelemetry configurations and utilities.
package otel

// Config holds OpenTelemetry configuration options.
type Config struct {
	// TraceSQLEnabled enables OpenTelemetry tracing for SQL operations.
	TraceSQLEnabled bool `json:"traceSQLEnabled"`
	// TraceRequestEnabled enables tracing of HTTP request parameters.
	TraceRequestEnabled bool `json:"traceRequestEnabled"`
	// TraceResponseEnabled enables tracing of HTTP response parameters.
	TraceResponseEnabled bool `json:"traceResponseEnabled"`
}

// NewConfig creates and returns a new OTEL configuration with default settings.
func NewConfig() *Config {
	return &Config{
		TraceSQLEnabled:      false,
		TraceRequestEnabled:  false,
		TraceResponseEnabled: false,
	}
}

// IsTracingSQLEnabled returns whether SQL tracing is enabled.
func (c *Config) IsTracingSQLEnabled() bool {
	return c.TraceSQLEnabled
}

// IsTracingRequestEnabled returns whether HTTP request tracing is enabled.
func (c *Config) IsTracingRequestEnabled() bool {
	return c.TraceRequestEnabled
}

// IsTracingResponseEnabled returns whether HTTP response tracing is enabled.
func (c *Config) IsTracingResponseEnabled() bool {
	return c.TraceResponseEnabled
}
