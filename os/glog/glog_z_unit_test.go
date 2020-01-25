// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"context"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func Test_SetConfigWithMap(t *testing.T) {
	gtest.Case(t, func() {
		l := New()
		m := map[string]interface{}{
			"path":     "/var/log",
			"level":    "all",
			"stdout":   false,
			"StStatus": 0,
		}
		err := l.SetConfigWithMap(m)
		gtest.Assert(err, nil)
		gtest.Assert(l.config.Path, m["path"])
		gtest.Assert(l.config.Level, m["level"])
		gtest.Assert(l.config.StdoutPrint, m["stdout"])
	})
}

func TestContextLog(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "name", "test")
	ctx = context.WithValue(ctx, "age", 18)

	Ctx(ctx, "age", "name").Println("Hello gf")

	log := New()
	log.Ctx(ctx, "name").Println("Hello gf")
}
