// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Time converts `any` to time.Time.
func (c *Converter) Time(anyInput any, format ...string) (time.Time, error) {
	// Unwrap IVal / single-entry maps first so a typed time value is not
	// serialized through gtime.String(), which drops the location.
	anyInput = unwrapTimeSource(anyInput)
	// It's already this type.
	if len(format) == 0 {
		switch v := anyInput.(type) {
		case time.Time:
			return v, nil
		case *time.Time:
			if v == nil {
				return time.Time{}, nil
			}
			return *v, nil
		case gtime.Time:
			return v.Time, nil
		case *gtime.Time:
			if v == nil {
				return time.Time{}, nil
			}
			return v.Time, nil
		}
	}
	t, err := c.GTime(anyInput, format...)
	if err != nil {
		return time.Time{}, err
	}
	if t != nil {
		return t.Time, nil
	}
	return time.Time{}, nil
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func (c *Converter) Duration(anyInput any) (time.Duration, error) {
	// It's already this type.
	if v, ok := anyInput.(time.Duration); ok {
		return v, nil
	}
	s, err := c.String(anyInput)
	if err != nil {
		return 0, err
	}
	if !utils.IsNumeric(s) {
		return gtime.ParseDuration(s)
	}
	i, err := c.Int64(anyInput)
	if err != nil {
		return 0, err
	}
	return time.Duration(i), nil
}

// GTime converts `any` to *gtime.Time.
// The parameter `format` can be used to specify the format of `any`.
// It returns the converted value that matched the first format of the formats slice.
// If no `format` given, it converts `any` using gtime.NewFromTimeStamp if `any` is numeric,
// or using gtime.StrToTime if `any` is string.
func (c *Converter) GTime(anyInput any, format ...string) (*gtime.Time, error) {
	if empty.IsNil(anyInput) {
		return nil, nil
	}
	// Unwrap IVal / single-entry maps before interface and type matching so
	// ORM values such as gvar.Var and {"now": *gtime.Time} keep their location.
	anyInput = unwrapTimeSource(anyInput)
	if empty.IsNil(anyInput) {
		return nil, nil
	}
	if v, ok := anyInput.(localinterface.IGTime); ok {
		return v.GTime(format...), nil
	}
	// It's already this type.
	if len(format) == 0 {
		switch v := anyInput.(type) {
		case *gtime.Time:
			return v, nil
		case gtime.Time:
			return &v, nil
		case time.Time:
			return gtime.NewFromTime(v), nil
		case *time.Time:
			if v == nil {
				return nil, nil
			}
			return gtime.NewFromTime(*v), nil
		}
	}
	s, err := c.String(anyInput)
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return gtime.New(), nil
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		for _, item := range format {
			t, err := gtime.StrToTimeFormat(s, item)
			if err != nil {
				return nil, err
			}
			if t != nil {
				return t, nil
			}
		}
		return nil, nil
	}
	if utils.IsNumeric(s) {
		i, err := c.Int64(s)
		if err != nil {
			return nil, err
		}
		return gtime.NewFromTimeStamp(i), nil
	} else {
		return gtime.StrToTime(s)
	}
}

// unwrapTimeSource extracts a scalar time source from wrappers that would
// otherwise be serialized through a timezone-less string.
//
// It unwraps:
//  1. localinterface.IVal implementations such as *gvar.Var
//  2. a map with exactly one entry, which is the ORM SELECT NOW() / Structs
//     into []time.Time shape: map[string]any{"now": *gtime.Time}
func unwrapTimeSource(anyInput any) any {
	const maxUnwrapDepth = 8
	for i := 0; i < maxUnwrapDepth; i++ {
		if empty.IsNil(anyInput) {
			return nil
		}
		if v, ok := anyInput.(localinterface.IVal); ok {
			if empty.IsNil(v) {
				return nil
			}
			next := v.Val()
			if !isSameTimeSource(next, anyInput) {
				anyInput = next
				continue
			}
			// IVal returned itself; still try single-value map unwrap below.
		}
		if next, ok := unwrapSingleValueMap(anyInput); ok {
			anyInput = next
			continue
		}
		return anyInput
	}
	return anyInput
}

// isSameTimeSource reports whether next still refers to the same source as
// current. Interface comparison is skipped for non-comparable dynamic types
// such as a named map that implements IVal and returns itself.
func isSameTimeSource(next, current any) bool {
	if next == nil || current == nil {
		return next == nil && current == nil
	}
	nextValue := reflect.ValueOf(next)
	currentValue := reflect.ValueOf(current)
	if !nextValue.IsValid() || !currentValue.IsValid() {
		return !nextValue.IsValid() && !currentValue.IsValid()
	}
	if !nextValue.Type().Comparable() || !currentValue.Type().Comparable() {
		return isSameNonComparableTimeSource(nextValue, currentValue)
	}
	return next == current
}

// isSameNonComparableTimeSource compares values that cannot use interface
// equality, without calling reflect.Value.Pointer on unsupported kinds.
func isSameNonComparableTimeSource(nextValue, currentValue reflect.Value) bool {
	if nextValue.Type() != currentValue.Type() {
		return false
	}
	switch nextValue.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return nextValue.Pointer() == currentValue.Pointer()
	default:
		// Stop unwrapping to avoid an interface-comparison panic.
		return true
	}
}

// unwrapSingleValueMap returns the only value of a one-entry map.
// Multi-entry maps are left unchanged so Structs into []time.Time does not
// pick an arbitrary field.
func unwrapSingleValueMap(anyInput any) (any, bool) {
	switch m := anyInput.(type) {
	case map[string]any:
		if len(m) != 1 {
			return nil, false
		}
		for _, value := range m {
			return value, true
		}
	case map[any]any:
		if len(m) != 1 {
			return nil, false
		}
		for _, value := range m {
			return value, true
		}
	default:
		rv := reflect.ValueOf(anyInput)
		for rv.Kind() == reflect.Pointer {
			if rv.IsNil() {
				return nil, false
			}
			rv = rv.Elem()
		}
		if rv.Kind() != reflect.Map || rv.Len() != 1 {
			return nil, false
		}
		iter := rv.MapRange()
		if iter.Next() {
			return iter.Value().Interface(), true
		}
	}
	return nil, false
}

// assignIfTimeDestination converts params into dest when dest is time.Time or
// gtime.Time (value or pointer). It preserves the source location.
func (c *Converter) assignIfTimeDestination(dest reflect.Value, params any) (bool, error) {
	if !dest.IsValid() {
		return false, nil
	}
	destType := dest.Type()
	for destType.Kind() == reflect.Pointer {
		destType = destType.Elem()
	}
	switch destType {
	case reflect.TypeOf(time.Time{}):
		t, err := c.Time(params)
		if err != nil {
			return true, err
		}
		return true, setConvertedTimeValue(dest, t)
	case reflect.TypeOf(gtime.Time{}):
		t, err := c.GTime(params)
		if err != nil {
			return true, err
		}
		if t == nil {
			t = gtime.New()
		}
		return true, setConvertedGTimeValue(dest, *t)
	default:
		return false, nil
	}
}
