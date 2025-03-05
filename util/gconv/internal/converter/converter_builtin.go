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
	*to.Addr().Interface().(*time.Time) = t
	return nil
}

func (c *Converter) builtInAnyConvertFuncForGTime(from any, to reflect.Value) error {
	v, err := c.GTime(from)
	if err != nil {
		return err
	}
	if v == nil {
		v = gtime.New()
	}
	*to.Addr().Interface().(*gtime.Time) = *v
	return nil
}
