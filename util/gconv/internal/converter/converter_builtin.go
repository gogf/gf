// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

func (c *Converter) builtInAnyConvertFuncForInt64(from any, to reflect.Value) error {
	v, err := c.Int64(from)
	if err != nil {
		return err
	}
	to.SetInt(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForUint64(from any, to reflect.Value) error {
	v, err := c.Uint64(from)
	if err != nil {
		return err
	}
	to.SetUint(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForString(from any, to reflect.Value) error {
	v, err := c.String(from)
	if err != nil {
		return err
	}
	to.SetString(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForFloat64(from any, to reflect.Value) error {
	v, err := c.Float64(from)
	if err != nil {
		return err
	}
	to.SetFloat(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForBool(from any, to reflect.Value) error {
	v, err := c.Bool(from)
	if err != nil {
		return err
	}
	to.SetBool(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForBytes(from any, to reflect.Value) error {
	v, err := c.Bytes(from)
	if err != nil {
		return err
	}
	to.SetBytes(v)
	return nil
}

func (c *Converter) builtInAnyConvertFuncForTime(from any, to reflect.Value) error {
	t, err := c.Time(from)
	if err != nil {
		return err
	}
	return setConvertedTimeValue(to, t)
}

func (c *Converter) builtInAnyConvertFuncForGTime(from any, to reflect.Value) error {
	v, err := c.GTime(from)
	if err != nil {
		return err
	}
	if v == nil {
		v = gtime.New()
	}
	return setConvertedGTimeValue(to, *v)
}

// setConvertedTimeValue assigns t to to without requiring to.Addr().
// Struct conversion may pass an unaddressable time.Time value.
func setConvertedTimeValue(to reflect.Value, t time.Time) error {
	switch {
	case to.Kind() == reflect.Struct && to.CanSet():
		to.Set(reflect.ValueOf(t))
		return nil
	case to.Kind() == reflect.Pointer:
		if to.IsNil() {
			if !to.CanSet() {
				return nil
			}
			to.Set(reflect.New(to.Type().Elem()))
		}
		return setConvertedTimeValue(to.Elem(), t)
	case to.CanAddr():
		*to.Addr().Interface().(*time.Time) = t
		return nil
	default:
		return nil
	}
}

// setConvertedGTimeValue assigns v to to without requiring to.Addr().
func setConvertedGTimeValue(to reflect.Value, v gtime.Time) error {
	switch {
	case to.Kind() == reflect.Struct && to.CanSet():
		to.Set(reflect.ValueOf(v))
		return nil
	case to.Kind() == reflect.Pointer:
		if to.IsNil() {
			if !to.CanSet() {
				return nil
			}
			to.Set(reflect.New(to.Type().Elem()))
		}
		return setConvertedGTimeValue(to.Elem(), v)
	case to.CanAddr():
		*to.Addr().Interface().(*gtime.Time) = v
		return nil
	default:
		return nil
	}
}
