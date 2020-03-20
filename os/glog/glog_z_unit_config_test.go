// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"github.com/gogf/gf/test/gtest"
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
		t.Assert(l.config.Level, m["level"])
		t.Assert(l.config.StdoutPrint, m["stdout"])
	})
}
