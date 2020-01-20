// Copyright 2017-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genv_test

import (
	"os"
	"testing"

	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
)

func Test_GEnv_All(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(os.Environ(), genv.All())
	})
}

func Test_GEnv_Map(t *testing.T) {
	gtest.Case(t, func() {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		gtest.Assert(err, nil)
		gtest.Assert(genv.Map()[key], "TEST")
	})
}

func Test_GEnv_Get(t *testing.T) {
	gtest.Case(t, func() {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		gtest.Assert(err, nil)
		gtest.AssertEQ(genv.Get(key), "TEST")
	})
}

func Test_GEnv_Contains(t *testing.T) {
	gtest.Case(t, func() {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		gtest.Assert(err, nil)
		gtest.AssertEQ(genv.Contains(key), true)
		gtest.AssertEQ(genv.Contains("none"), false)
	})
}

func Test_GEnv_Set(t *testing.T) {
	gtest.Case(t, func() {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := genv.Set(key, "TEST")
		gtest.Assert(err, nil)
		gtest.AssertEQ(os.Getenv(key), "TEST")
	})
}

func Test_GEnv_Build(t *testing.T) {
	gtest.Case(t, func() {
		s := genv.Build(map[string]string{
			"k1": "v1",
			"k2": "v2",
		})
		gtest.AssertIN("k1=v1", s)
		gtest.AssertIN("k2=v2", s)
	})
}

func Test_GEnv_Remove(t *testing.T) {
	gtest.Case(t, func() {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		gtest.Assert(err, nil)
		err = genv.Remove(key)
		gtest.Assert(err, nil)
		gtest.AssertEQ(os.Getenv(key), "")
	})
}
