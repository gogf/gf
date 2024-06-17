// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/gogf/gf/v2/util/gconv"
)

func generateAddOptions(
	meterOption gmetric.MeterOption, constOption metric.MeasurementOption, option ...gmetric.Option,
) []metric.AddOption {
	var (
		addOptions             = make([]metric.AddOption, 0)
		globalAttributesOption = getGlobalAttributesOption(gmetric.GetGlobalAttributesOption{
			Instrument:        meterOption.Instrument,
			InstrumentVersion: meterOption.InstrumentVersion,
		})
	)
	if constOption != nil {
		addOptions = append(addOptions, constOption)
	}
	if globalAttributesOption != nil {
		addOptions = append(addOptions, globalAttributesOption)
	}
	if len(option) > 0 {
		addOptions = append(
			addOptions,
			metric.WithAttributes(attributesToKeyValues(option[0].Attributes)...),
		)
	}
	return addOptions
}

func getGlobalAttributesOption(option gmetric.GetGlobalAttributesOption) metric.MeasurementOption {
	var (
		globalAttributesOption metric.MeasurementOption
		globalAttributes       = gmetric.GetGlobalAttributes(gmetric.GetGlobalAttributesOption{})
		instrumentAttributes   gmetric.Attributes
	)
	if option.Instrument != "" {
		instrumentAttributes = gmetric.GetGlobalAttributes(option)
	}
	if len(globalAttributes) > 0 {
		globalAttributesOption = metric.WithAttributes(attributesToKeyValues(globalAttributes)...)
	}
	if len(instrumentAttributes) > 0 {
		globalAttributesOption = metric.WithAttributes(attributesToKeyValues(instrumentAttributes)...)
	}
	return globalAttributesOption
}

func getDynamicOptionByMetricOption(option ...gmetric.Option) metric.MeasurementOption {
	var (
		usedOption    gmetric.Option
		dynamicOption metric.MeasurementOption
	)
	if len(option) > 0 {
		usedOption = option[0]
	}
	if len(usedOption.Attributes) > 0 {
		dynamicOption = metric.WithAttributes(attributesToKeyValues(usedOption.Attributes)...)
	}
	return dynamicOption
}

func genConstOptionForMetric(
	meterOption gmetric.MeterOption,
	metricOption gmetric.MetricOption,
) metric.MeasurementOption {
	return genConstOptionForMetricByAttributes(meterOption.Attributes, metricOption.Attributes)
}

func getConstOptionByMetric(meterOption gmetric.MeterOption, m gmetric.Metric) metric.MeasurementOption {
	return genConstOptionForMetricByAttributes(meterOption.Attributes, m.Info().Attributes())
}

func genConstOptionForMetricByAttributes(
	meterAttrs gmetric.Attributes,
	metricAttrs gmetric.Attributes,
) metric.MeasurementOption {
	var (
		constOption metric.MeasurementOption
		attributes  = make([]attribute.KeyValue, 0)
	)
	if len(meterAttrs) > 0 {
		attributes = append(attributes, attributesToKeyValues(meterAttrs)...)
	}
	if len(metricAttrs) > 0 {
		attributes = append(attributes, attributesToKeyValues(metricAttrs)...)
	}
	constOption = metric.WithAttributes(attributes...)
	return constOption
}

func metricToFloat64Observable(m gmetric.Metric) metric.Float64Observable {
	performer := m.(gmetric.PerformerExporter).Performer()
	switch m.Info().Type() {
	case gmetric.MetricTypeObservableCounter:
		return performer.(*localObservableCounterPerformer).Float64ObservableCounter

	case gmetric.MetricTypeObservableUpDownCounter:
		return performer.(*localObservableUpDownCounterPerformer).Float64ObservableUpDownCounter

	case gmetric.MetricTypeObservableGauge:
		return performer.(*localObservableGaugePerformer).Float64ObservableGauge

	default:
		panic(gerror.NewCode(
			gcode.CodeInvalidParameter,
			`Histogram is not support for converting to metric.Float64Observable`,
		))
	}
	return nil
}

// attributesToKeyValues converts attributes to OpenTelemetry key-value pair attributes.
func attributesToKeyValues(attrs gmetric.Attributes) []attribute.KeyValue {
	var keyValues = make([]attribute.KeyValue, 0)
	for _, attr := range attrs {
		keyValues = append(keyValues, attributeToKeyValue(attr))
	}
	return keyValues
}

// attributeToKeyValue converts attribute to OpenTelemetry key-value pair attribute.
func attributeToKeyValue(attr gmetric.Attribute) attribute.KeyValue {
	var (
		key   = string(attr.Key())
		value = attr.Value()
	)
	switch result := value.(type) {
	case bool:
		return attribute.Bool(key, result)
	case []bool:
		return attribute.BoolSlice(key, result)

	case int:
		return attribute.Int(key, result)
	case []int:
		return attribute.IntSlice(key, result)

	case int64:
		return attribute.Int64(key, result)
	case []int64:
		return attribute.Int64Slice(key, result)

	case float64:
		return attribute.Float64(key, result)
	case []float64:
		return attribute.Float64Slice(key, result)

	case string:
		return attribute.String(key, result)
	case []string:
		return attribute.StringSlice(key, result)

	default:
		return attribute.String(key, gconv.String(value))
	}
}
