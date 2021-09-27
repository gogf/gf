// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"github.com/gogf/gf/test/gtest"
	"strings"
	"testing"
)

func Test_SetConfigWithMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		m := map[string]interface{}{
			"path":     "/var/log",
			"level":    "all",
			"stdout":   false,
			"StStatus": 0,
		}
		err := l.SetConfigWithMap(m)
		t.Assert(err, nil)
		t.Assert(l.config.Path, m["path"])
		t.Assert(l.config.Level, LEVEL_ALL)
		t.Assert(l.config.StdoutPrint, m["stdout"])
	})
}

func Test_SetConfigWithMap_LevelStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		buffer := bytes.NewBuffer(nil)
		l := New()
		m := map[string]interface{}{
			"level": "all",
		}
		err := l.SetConfigWithMap(m)
		t.Assert(err, nil)

		l.SetWriter(buffer)

		l.Debug(ctx, "test")
		l.Warning(ctx, "test")
		t.Assert(strings.Contains(buffer.String(), "DEBU"), true)
		t.Assert(strings.Contains(buffer.String(), "WARN"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		buffer := bytes.NewBuffer(nil)
		l := New()
		m := map[string]interface{}{
			"level": "warn",
		}
		err := l.SetConfigWithMap(m)
		t.Assert(err, nil)
		l.SetWriter(buffer)
		l.Debug(ctx, "test")
		l.Warning(ctx, "test")
		t.Assert(strings.Contains(buffer.String(), "DEBU"), false)
		t.Assert(strings.Contains(buffer.String(), "WARN"), true)
	})
}
