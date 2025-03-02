// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

func builtInAnyConvertFuncForInt64(from any, to reflect.Value) error {
	v, err := doInt64(from)
	if err != nil {
		return err
	}
	to.SetInt(v)
	return nil
}

func builtInAnyConvertFuncForUint64(from any, to reflect.Value) error {
	v, err := doUint64(from)
	if err != nil {
		return err
	}
	to.SetUint(v)
	return nil
}

func builtInAnyConvertFuncForString(from any, to reflect.Value) error {
	v, err := doString(from)
	if err != nil {
		return err
	}
	to.SetString(v)
	return nil
}

func builtInAnyConvertFuncForFloat64(from any, to reflect.Value) error {
	v, err := doFloat64(from)
	if err != nil {
		return err
	}
	to.SetFloat(v)
	return nil
}

func builtInAnyConvertFuncForBool(from any, to reflect.Value) error {
	v, err := doBool(from)
	if err != nil {
		return err
	}
	to.SetBool(v)
	return nil
}

func builtInAnyConvertFuncForBytes(from any, to reflect.Value) error {
	v, err := doBytes(from)
	if err != nil {
		return err
	}
	to.SetBytes(v)
	return nil
}

func builtInAnyConvertFuncForTime(from any, to reflect.Value) error {
	*to.Addr().Interface().(*time.Time) = Time(from)
	return nil
}

func builtInAnyConvertFuncForGTime(from any, to reflect.Value) error {
	v := GTime(from)
	if v == nil {
		v = gtime.New()
	}
	*to.Addr().Interface().(*gtime.Time) = *v
	return nil
}
