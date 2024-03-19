// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genv_test

import (
	"os"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestGEnvAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(os.Environ(), genv.All())
	})
}

func TestGEnvMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.AssertNil(err)
		t.Assert(genv.Map()[key], "TEST")
	})
}

func TestGEnvGet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.AssertNil(err)
		t.AssertEQ(genv.Get(key).String(), "TEST")
	})
}

func TestGEnvGetVar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.AssertNil(err)
		t.AssertEQ(genv.Get(key).String(), "TEST")
	})
}

func TestGEnvContains(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.AssertNil(err)
		t.AssertEQ(genv.Contains(key), true)
		t.AssertEQ(genv.Contains("none"), false)
	})
}

func TestGEnvSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := genv.Set(key, "TEST")
		t.AssertNil(err)
		t.AssertEQ(os.Getenv(key), "TEST")
	})
}

func TestGEnvSetMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := genv.SetMap(g.MapStrStr{
			"K1": "TEST1",
			"K2": "TEST2",
		})
		t.AssertNil(err)
		t.AssertEQ(os.Getenv("K1"), "TEST1")
		t.AssertEQ(os.Getenv("K2"), "TEST2")
	})
}

func TestGEnvBuild(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := genv.Build(map[string]string{
			"k1": "v1",
			"k2": "v2",
		})
		t.AssertIN("k1=v1", s)
		t.AssertIN("k2=v2", s)
	})
}

func TestGEnvRemove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := gconv.String(gtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.AssertNil(err)
		err = genv.Remove(key)
		t.AssertNil(err)
		t.AssertEQ(os.Getenv(key), "")
	})
}

func TestGetWithCmd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gcmd.Init("-test", "2")
		t.Assert(genv.GetWithCmd("TEST"), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		genv.Set("TEST", "1")
		defer genv.Remove("TEST")
		gcmd.Init("-test", "2")
		t.Assert(genv.GetWithCmd("test"), 1)
	})
}

func TestMapFromEnv(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := genv.MapFromEnv([]string{"a=1", "b=2"})
		t.Assert(m, g.Map{"a": 1, "b": 2})
	})
}

func TestMapToEnv(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := genv.MapToEnv(g.MapStrStr{"a": "1"})
		t.Assert(s, []string{"a=1"})
	})
}

func TestFilter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := genv.Filter([]string{"a=1", "a=3"})
		t.Assert(s, []string{"a=3"})
	})
}
