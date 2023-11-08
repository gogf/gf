// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otelmetric

import (
	"go.opentelemetry.io/otel/attribute"

	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/gogf/gf/v2/util/gconv"
)

func attributesToKeyValues(attrs gmetric.Attributes) []attribute.KeyValue {
	var keyValues = make([]attribute.KeyValue, 0)
	for _, attr := range attrs {
		keyValues = append(keyValues, attributeToKeyValue(attr))
	}
	return keyValues
}

func attributeToKeyValue(attr gmetric.Attribute) attribute.KeyValue {
	var (
		key   = attr.Key()
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
