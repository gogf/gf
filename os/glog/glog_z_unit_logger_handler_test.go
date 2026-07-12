// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

var arrayForHandlerTest1 = garray.NewStrArray()

func customHandler1(ctx context.Context, input *glog.HandlerInput) {
	arrayForHandlerTest1.Append(input.String(false))
	input.Next(ctx)
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

func TestLogger_SetHandlers_HandlerJson(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		l.SetHandlers(glog.HandlerJson)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Debug(ctx, g.Map{"Uid": 100, "Name": "john"})
		output := readJsonHandlerOutput(t, w)
		t.AssertNE(output["Time"], nil)
		t.Assert(output["Trace-Id"], "1234567890")
		t.Assert(output["Span-Id"], "abcdefg")
		t.Assert(output["Level"], "DEBU")
		t.Assert(output["Name"], "john")
		t.Assert(output["Uid"], float64(100))
		t.Assert(output["CtxStr"], nil)
		t.Assert(output["Content"], nil)
		t.Assert(gstr.Contains(w.String(), `"{\"`), false)
	})
}

func TestLogger_SetHandlers_HandlerJson_Content(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		l.SetHandlers(glog.HandlerJson)

		l.Debug(context.Background(), "plain content")
		output := readJsonHandlerOutput(t, w)
		t.Assert(output["Content"], "plain content")
		t.Assert(output["Level"], "DEBU")
	})
}

func TestLogger_SetHandlers_HandlerJson_ContentCollision(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		l.SetHandlers(glog.HandlerJson)
		l.SetCtxKeys("Trace-Id")
		ctx := context.WithValue(context.Background(), "Trace-Id", "ctx-trace")

		l.Debug(ctx, g.Map{
			"Level":    "payload-level",
			"Trace-Id": "payload-trace",
			"Name":     "john",
		})
		output := readJsonHandlerOutput(t, w)
		t.Assert(output["Level"], "DEBU")
		t.Assert(output["Trace-Id"], "ctx-trace")
		t.Assert(output["Name"], "john")
	})
}

func readJsonHandlerOutput(t *gtest.T, w *bytes.Buffer) map[string]any {
	var output map[string]any
	err := json.Unmarshal(bytes.TrimSpace(w.Bytes()), &output)
	t.AssertNil(err)
	return output
}

func TestLogger_SetHandlers_HandlerStructure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		l.SetHandlers(glog.HandlerStructure)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Debug(ctx, "debug", "uid", 1000)
		l.Info(ctx, "info", "' '", `"\n`)

		t.Assert(gstr.Count(w.String(), "uid=1000"), 1)
		t.Assert(gstr.Count(w.String(), "Content=debug"), 1)
		t.Assert(gstr.Count(w.String(), `"' '"="\"\\n"`), 1)
		t.Assert(gstr.Count(w.String(), `CtxStr="1234567890, abcdefg"`), 2)
	})
}

func Test_SetDefaultHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		oldHandler := glog.GetDefaultHandler()
		glog.SetDefaultHandler(func(ctx context.Context, in *glog.HandlerInput) {
			glog.HandlerJson(ctx, in)
		})
		defer glog.SetDefaultHandler(oldHandler)

		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Debug(ctx, 1, 2, 3)
		t.Assert(gstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gstr.Count(w.String(), `"1 2 3"`), 1)
		t.Assert(gstr.Count(w.String(), `"DEBU"`), 1)
	})
}
