// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog_test

import (
	"bytes"
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/text/gstr"
	"testing"
)

func Test_Ctx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Ctx(ctx).Print(1, 2, 3)
		t.Assert(gstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 1)
	})
}

func Test_Ctx_Config(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		m := map[string]interface{}{
			"CtxKeys": g.SliceStr{"Trace-Id", "Span-Id", "Test"},
		}
		err := l.SetConfigWithMap(m)
		t.Assert(err, nil)
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Ctx(ctx).Print(1, 2, 3)
		t.Assert(gstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 1)
	})
}
