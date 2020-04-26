// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcfg_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gogf/gf/os/genv"

	"github.com/gogf/gf/os/gcfg"

	"github.com/gogf/gf/test/gtest"
)

func init() {
	os.Setenv("GF_GCFG_ERRORPRINT", "false")
}

func Test_ExpandValueEnv(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gcfg.ExpandValueEnv("${GOGF||/usr/local/go}"), "/usr/local/go")
		gtest.Assert(gcfg.ExpandValueEnv("gogf"), "gogf")
	})
}

func Test_ExpandValueEnvForStr(t *testing.T) {
	content := `
v1    = "${V1}"
v2    = "${V2||1}"
v3    = "${V3}"
v4    = "${V4||1.23}"
array = "${ARRAY}"
[redis]
    disk  = "${REDIS_DISK}"
    cache = "${REDIS_CACHE||127.0.0.1:6379,1}"
`
	genv.Set("V1", "1")
	genv.Set("V3", "1.23")
	genv.Set("REDIS_DISK", "127.0.0.1:6379,0")
	genv.Set("ARRAY", "[1,2,3]")
	gcfg.SetContent(content)
	defer func() {
		gcfg.ClearContent()
		genv.Remove("V1")
		genv.Remove("V3")
		genv.Remove("REDIS_DISK")
		genv.Remove("ARRAY")
	}()

	gtest.Case(t, func() {
		c := gcfg.New()
		fmt.Println(c.GetString(content))
		gtest.AssertEQ(c.GetInt("v1"), 1)
		gtest.AssertEQ(c.GetString("v2"), "1")
		gtest.AssertEQ(c.GetFloat32("v3"), float32(1.23))
		gtest.AssertEQ(c.GetString("v4"), "1.23")
		gtest.AssertEQ(c.GetStrings("array"), []string{"1", "2", "3"})
		gtest.AssertEQ(c.GetMap("redis"), map[string]interface{}{
			"disk":  "127.0.0.1:6379,0",
			"cache": "127.0.0.1:6379,1",
		})
	})
}

func Test_envStrParse(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gcfg.EnvStrParse("\"strtest\""), "\"strtest\"")
		gtest.Assert(gcfg.EnvStrParse("true"), "true")
		gtest.Assert(gcfg.EnvStrParse("TRUE"), "true")
		gtest.Assert(gcfg.EnvStrParse("True"), "true")
		gtest.Assert(gcfg.EnvStrParse("false"), "false")
		gtest.Assert(gcfg.EnvStrParse("FALSE"), "false")
		gtest.Assert(gcfg.EnvStrParse("False"), "false")
		gtest.Assert(gcfg.EnvStrParse("1"), "1")
		gtest.Assert(gcfg.EnvStrParse("-1"), "-1")
		gtest.Assert(gcfg.EnvStrParse("1.1"), "1.1")
		gtest.Assert(gcfg.EnvStrParse("-1.1"), "-1.1")
		gtest.Assert(gcfg.EnvStrParse("[1,2,3]"), "[1,2,3]")
		gtest.Assert(gcfg.EnvStrParse("$$"), "\"$$\"")
		gtest.Assert(gcfg.EnvStrParse("test"), "\"test\"")
	})
}
