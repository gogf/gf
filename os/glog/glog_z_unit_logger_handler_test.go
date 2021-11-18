// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

var arrayForHandlerTest1 = garray.NewStrArray()

func customHandler1(ctx context.Context, input *glog.HandlerInput) {
	arrayForHandlerTest1.Append(input.String(false))
	input.Next()
}

func TestLogger_SetHandlers1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		l.SetHandlers(customHandler1)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Print(ctx, 1, 2, 3)
		t.Assert(gstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 1)

		t.Assert(arrayForHandlerTest1.Len(), 1)
		t.Assert(gstr.Count(arrayForHandlerTest1.At(0), "1234567890"), 1)
		t.Assert(gstr.Count(arrayForHandlerTest1.At(0), "abcdefg"), 1)
		t.Assert(gstr.Count(arrayForHandlerTest1.At(0), "1 2 3"), 1)
	})
}

var arrayForHandlerTest2 = garray.NewStrArray()

func customHandler2(ctx context.Context, input *glog.HandlerInput) {
	arrayForHandlerTest2.Append(input.String(false))
}

func TestLogger_SetHandlers2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		l.SetHandlers(customHandler2)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Print(ctx, 1, 2, 3)
		t.Assert(gstr.Count(w.String(), "1234567890"), 0)
		t.Assert(gstr.Count(w.String(), "abcdefg"), 0)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 0)

		t.Assert(arrayForHandlerTest2.Len(), 1)
		t.Assert(gstr.Count(arrayForHandlerTest2.At(0), "1234567890"), 1)
		t.Assert(gstr.Count(arrayForHandlerTest2.At(0), "abcdefg"), 1)
		t.Assert(gstr.Count(arrayForHandlerTest2.At(0), "1 2 3"), 1)
	})
}
