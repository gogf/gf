// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcfg_test

import (
    "github.com/gogf/gf/g/os/gcfg"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
)

func Test_Basic(t *testing.T) {
    config := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
    gtest.Case(t, func() {
        path := gcfg.DEFAULT_CONFIG_FILE
        err  := gfile.PutContents(path, config)
        gtest.Assert(err, nil)
        defer gfile.Remove(path)

        c := gcfg.New()
        gtest.Assert(c.Get("v1"),           1)
        gtest.AssertEQ(c.GetInt("v1"),      1)
        gtest.AssertEQ(c.GetInt8("v1"),     int8(1))
        gtest.AssertEQ(c.GetInt16("v1"),    int16(1))
        gtest.AssertEQ(c.GetInt32("v1"),    int32(1))
        gtest.AssertEQ(c.GetInt64("v1"),    int64(1))
        gtest.AssertEQ(c.GetUint("v1"),     uint(1))
        gtest.AssertEQ(c.GetUint8("v1"),    uint8(1))
        gtest.AssertEQ(c.GetUint16("v1"),   uint16(1))
        gtest.AssertEQ(c.GetUint32("v1"),   uint32(1))
        gtest.AssertEQ(c.GetUint64("v1"),   uint64(1))

        gtest.AssertEQ(c.GetVar("v1").String(), "1")
        gtest.AssertEQ(c.GetVar("v1").Bool(),   true)
        gtest.AssertEQ(c.GetVar("v2").String(), "true")
        gtest.AssertEQ(c.GetVar("v2").Bool(),   true)

        gtest.AssertEQ(c.GetString("v1"),  "1")
        gtest.AssertEQ(c.GetFloat32("v4"), float32(1.23))
        gtest.AssertEQ(c.GetFloat64("v4"), float64(1.23))
        gtest.AssertEQ(c.GetString("v2"),  "true")
        gtest.AssertEQ(c.GetBool("v2"),    true)
        gtest.AssertEQ(c.GetBool("v3"),    false)

        gtest.AssertEQ(c.Contains("v1"),    true)
        gtest.AssertEQ(c.Contains("v2"),    true)
        gtest.AssertEQ(c.Contains("v3"),    true)
        gtest.AssertEQ(c.Contains("v4"),    true)
        gtest.AssertEQ(c.Contains("v5"),    false)

        gtest.AssertEQ(c.GetInts("array"),        []int{1,2,3})
        gtest.AssertEQ(c.GetStrings("array"),     []string{"1","2","3"})
        gtest.AssertEQ(c.GetArray("array"),       []interface{}{"1","2","3"})
        gtest.AssertEQ(c.GetInterfaces("array"),  []interface{}{"1","2","3"})
        gtest.AssertEQ(c.GetMap("redis"),         map[string]interface{}{
            "disk"  : "127.0.0.1:6379,0",
            "cache" : "127.0.0.1:6379,1",
        })
        gtest.AssertEQ(c.FilePath(),    gfile.Pwd() + gfile.Separator + path)

    })
}

func Test_Content(t *testing.T) {
    content := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
    gcfg.SetContent(content)
    defer gcfg.ClearContent()

    gtest.Case(t, func() {
        c := gcfg.New()
        gtest.Assert(c.Get("v1"),           1)
        gtest.AssertEQ(c.GetInt("v1"),      1)
        gtest.AssertEQ(c.GetInt8("v1"),     int8(1))
        gtest.AssertEQ(c.GetInt16("v1"),    int16(1))
        gtest.AssertEQ(c.GetInt32("v1"),    int32(1))
        gtest.AssertEQ(c.GetInt64("v1"),    int64(1))
        gtest.AssertEQ(c.GetUint("v1"),     uint(1))
        gtest.AssertEQ(c.GetUint8("v1"),    uint8(1))
        gtest.AssertEQ(c.GetUint16("v1"),   uint16(1))
        gtest.AssertEQ(c.GetUint32("v1"),   uint32(1))
        gtest.AssertEQ(c.GetUint64("v1"),   uint64(1))

        gtest.AssertEQ(c.GetVar("v1").String(), "1")
        gtest.AssertEQ(c.GetVar("v1").Bool(),   true)
        gtest.AssertEQ(c.GetVar("v2").String(), "true")
        gtest.AssertEQ(c.GetVar("v2").Bool(),   true)

        gtest.AssertEQ(c.GetString("v1"),  "1")
        gtest.AssertEQ(c.GetFloat32("v4"), float32(1.23))
        gtest.AssertEQ(c.GetFloat64("v4"), float64(1.23))
        gtest.AssertEQ(c.GetString("v2"),  "true")
        gtest.AssertEQ(c.GetBool("v2"),    true)
        gtest.AssertEQ(c.GetBool("v3"),    false)

        gtest.AssertEQ(c.Contains("v1"),    true)
        gtest.AssertEQ(c.Contains("v2"),    true)
        gtest.AssertEQ(c.Contains("v3"),    true)
        gtest.AssertEQ(c.Contains("v4"),    true)
        gtest.AssertEQ(c.Contains("v5"),    false)

        gtest.AssertEQ(c.GetInts("array"),        []int{1,2,3})
        gtest.AssertEQ(c.GetStrings("array"),     []string{"1","2","3"})
        gtest.AssertEQ(c.GetArray("array"),       []interface{}{"1","2","3"})
        gtest.AssertEQ(c.GetInterfaces("array"),  []interface{}{"1","2","3"})
        gtest.AssertEQ(c.GetMap("redis"),         map[string]interface{}{
            "disk"  : "127.0.0.1:6379,0",
            "cache" : "127.0.0.1:6379,1",
        })
    })
}

func Test_SetFileName(t *testing.T) {
    config := `
{
	"array": [
		1,
		2,
		3
	],
	"redis": {
		"cache": "127.0.0.1:6379,1",
		"disk": "127.0.0.1:6379,0"
	},
	"v1": 1,
	"v2": "true",
	"v3": "off",
	"v4": "1.234"
}
`
    gtest.Case(t, func() {
        path := "config.json"
        err  := gfile.PutContents(path, config)
        gtest.Assert(err, nil)
        defer gfile.Remove(path)

        c := gcfg.New()
        c.SetFileName(path)
        gtest.Assert(c.Get("v1"),           1)
        gtest.AssertEQ(c.GetInt("v1"),      1)
        gtest.AssertEQ(c.GetInt8("v1"),     int8(1))
        gtest.AssertEQ(c.GetInt16("v1"),    int16(1))
        gtest.AssertEQ(c.GetInt32("v1"),    int32(1))
        gtest.AssertEQ(c.GetInt64("v1"),    int64(1))
        gtest.AssertEQ(c.GetUint("v1"),     uint(1))
        gtest.AssertEQ(c.GetUint8("v1"),    uint8(1))
        gtest.AssertEQ(c.GetUint16("v1"),   uint16(1))
        gtest.AssertEQ(c.GetUint32("v1"),   uint32(1))
        gtest.AssertEQ(c.GetUint64("v1"),   uint64(1))

        gtest.AssertEQ(c.GetVar("v1").String(), "1")
        gtest.AssertEQ(c.GetVar("v1").Bool(),   true)
        gtest.AssertEQ(c.GetVar("v2").String(), "true")
        gtest.AssertEQ(c.GetVar("v2").Bool(),   true)

        gtest.AssertEQ(c.GetString("v1"),  "1")
        gtest.AssertEQ(c.GetFloat32("v4"), float32(1.234))
        gtest.AssertEQ(c.GetFloat64("v4"), float64(1.234))
        gtest.AssertEQ(c.GetString("v2"),  "true")
        gtest.AssertEQ(c.GetBool("v2"),    true)
        gtest.AssertEQ(c.GetBool("v3"),    false)

        gtest.AssertEQ(c.Contains("v1"),    true)
        gtest.AssertEQ(c.Contains("v2"),    true)
        gtest.AssertEQ(c.Contains("v3"),    true)
        gtest.AssertEQ(c.Contains("v4"),    true)
        gtest.AssertEQ(c.Contains("v5"),    false)

        gtest.AssertEQ(c.GetInts("array"),        []int{1,2,3})
        gtest.AssertEQ(c.GetStrings("array"),     []string{"1","2","3"})
        gtest.AssertEQ(c.GetArray("array"),       []interface{}{"1","2","3"})
        gtest.AssertEQ(c.GetInterfaces("array"),  []interface{}{"1","2","3"})
        gtest.AssertEQ(c.GetMap("redis"),         map[string]interface{}{
            "disk"  : "127.0.0.1:6379,0",
            "cache" : "127.0.0.1:6379,1",
        })
        gtest.AssertEQ(c.FilePath(),    gfile.Pwd() + gfile.Separator + path)

    })
}

func Test_Instance(t *testing.T) {
    config := `
{
	"array": [
		1,
		2,
		3
	],
	"redis": {
		"cache": "127.0.0.1:6379,1",
		"disk": "127.0.0.1:6379,0"
	},
	"v1": 1,
	"v2": "true",
	"v3": "off",
	"v4": "1.234"
}
`
    gtest.Case(t, func() {
        path := gcfg.DEFAULT_CONFIG_FILE
        err  := gfile.PutContents(path, config)
        gtest.Assert(err, nil)
        defer gfile.Remove(path)

        c := gcfg.Instance()
        gtest.Assert(c.Get("v1"),           1)
        gtest.AssertEQ(c.GetInt("v1"),      1)
        gtest.AssertEQ(c.GetInt8("v1"),     int8(1))
        gtest.AssertEQ(c.GetInt16("v1"),    int16(1))
        gtest.AssertEQ(c.GetInt32("v1"),    int32(1))
        gtest.AssertEQ(c.GetInt64("v1"),    int64(1))
        gtest.AssertEQ(c.GetUint("v1"),     uint(1))
        gtest.AssertEQ(c.GetUint8("v1"),    uint8(1))
        gtest.AssertEQ(c.GetUint16("v1"),   uint16(1))
        gtest.AssertEQ(c.GetUint32("v1"),   uint32(1))
        gtest.AssertEQ(c.GetUint64("v1"),   uint64(1))

        gtest.AssertEQ(c.GetVar("v1").String(), "1")
        gtest.AssertEQ(c.GetVar("v1").Bool(),   true)
        gtest.AssertEQ(c.GetVar("v2").String(), "true")
        gtest.AssertEQ(c.GetVar("v2").Bool(),   true)

        gtest.AssertEQ(c.GetString("v1"),  "1")
        gtest.AssertEQ(c.GetFloat32("v4"), float32(1.234))
        gtest.AssertEQ(c.GetFloat64("v4"), float64(1.234))
        gtest.AssertEQ(c.GetString("v2"),  "true")
        gtest.AssertEQ(c.GetBool("v2"),    true)
        gtest.AssertEQ(c.GetBool("v3"),    false)

        gtest.AssertEQ(c.Contains("v1"),    true)
        gtest.AssertEQ(c.Contains("v2"),    true)
        gtest.AssertEQ(c.Contains("v3"),    true)
        gtest.AssertEQ(c.Contains("v4"),    true)
        gtest.AssertEQ(c.Contains("v5"),    false)

        gtest.AssertEQ(c.GetInts("array"),        []int{1,2,3})
        gtest.AssertEQ(c.GetStrings("array"),     []string{"1","2","3"})
        gtest.AssertEQ(c.GetArray("array"),       []interface{}{"1","2","3"})
        gtest.AssertEQ(c.GetInterfaces("array"),  []interface{}{"1","2","3"})
        gtest.AssertEQ(c.GetMap("redis"),         map[string]interface{}{
            "disk"  : "127.0.0.1:6379,0",
            "cache" : "127.0.0.1:6379,1",
        })
        gtest.AssertEQ(c.FilePath(),    gfile.Pwd() + gfile.Separator + path)

    })
}