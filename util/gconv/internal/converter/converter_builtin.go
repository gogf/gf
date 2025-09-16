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
	// Enhanced timezone preservation: handle gtime.Time types directly first
	// before going through the general GTime converter to prevent timezone loss
	switch v := from.(type) {
	case *gtime.Time:
		if v == nil {
			v = gtime.New()
		}
		*to.Addr().Interface().(*gtime.Time) = *v
		return nil
	case gtime.Time:
		// Direct assignment to preserve timezone information
		*to.Addr().Interface().(*gtime.Time) = v
		return nil
	case map[string]interface{}:
		// Handle map inputs by extracting the first value and converting it directly
		// This prevents timezone loss that occurs when map is converted to JSON string
		if len(v) > 0 {
			for _, value := range v {
				// Convert the extracted value directly using c.GTime to preserve timezone
				gtimeResult, err := c.GTime(value)
				if err != nil {
					return err
				}
				if gtimeResult == nil {
					gtimeResult = gtime.New()
				}
				*to.Addr().Interface().(*gtime.Time) = *gtimeResult
				return nil
			}
		}
		// Empty map case
		*to.Addr().Interface().(*gtime.Time) = *gtime.New()
		return nil
	}

	// For other types, use the general GTime converter
	// The c.GTime method already handles timezone preservation for known types
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
