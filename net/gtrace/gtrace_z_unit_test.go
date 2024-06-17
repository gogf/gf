// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtrace_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/encoding/gcompress"

	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func TestWithTraceID(t *testing.T) {
	var (
		ctx  = context.Background()
		uuid = `a323f910-f690-11ec-963d-79c0b7fcf119`
	)
	gtest.C(t, func(t *gtest.T) {
		newCtx, err := gtrace.WithTraceID(ctx, uuid)
		t.AssertNE(err, nil)
		t.Assert(newCtx, ctx)
	})
	gtest.C(t, func(t *gtest.T) {
		var traceId = gstr.Replace(uuid, "-", "")
		newCtx, err := gtrace.WithTraceID(ctx, traceId)
		t.AssertNil(err)
		t.AssertNE(newCtx, ctx)
		t.Assert(gtrace.GetTraceID(ctx), "")
		t.Assert(gtrace.GetTraceID(newCtx), traceId)
	})
}

func TestWithUUID(t *testing.T) {
	var (
		ctx  = context.Background()
		uuid = `a323f910-f690-11ec-963d-79c0b7fcf119`
	)
	gtest.C(t, func(t *gtest.T) {
		newCtx, err := gtrace.WithTraceID(ctx, uuid)
		t.AssertNE(err, nil)
		t.Assert(newCtx, ctx)
	})
	gtest.C(t, func(t *gtest.T) {
		newCtx, err := gtrace.WithUUID(ctx, uuid)
		t.AssertNil(err)
		t.AssertNE(newCtx, ctx)
		t.Assert(gtrace.GetTraceID(ctx), "")
		t.Assert(gtrace.GetTraceID(newCtx), gstr.Replace(uuid, "-", ""))
	})
}

func TestSafeContent(t *testing.T) {
	var (
		defText    = "中"
		shortData  = strings.Repeat(defText, gtrace.MaxContentLogSize()-1)
		standData  = strings.Repeat(defText, gtrace.MaxContentLogSize())
		longData   = strings.Repeat(defText, gtrace.MaxContentLogSize()+1)
		header     = http.Header{}
		gzipHeader = http.Header{
			"Content-Encoding": []string{"gzip"},
		}
	)

	// safe content
	gtest.C(t, func(t *gtest.T) {

		t1, err := gtrace.SafeContentForHttp([]byte(shortData), header)
		t.AssertNil(err)
		t.Assert(t1, shortData)
		t.Assert(gtrace.SafeContent([]byte(shortData)), shortData)

		t2, err := gtrace.SafeContentForHttp([]byte(standData), header)
		t.AssertNil(err)
		t.Assert(t2, standData)
		t.Assert(gtrace.SafeContent([]byte(standData)), standData)

		t3, err := gtrace.SafeContentForHttp([]byte(longData), header)
		t.AssertNil(err)
		t.Assert(t3, standData+"...")
		t.Assert(gtrace.SafeContent([]byte(longData)), standData+"...")
	})

	// compress content
	var (
		compressShortData, _ = gcompress.Gzip([]byte(shortData))
		compressStandData, _ = gcompress.Gzip([]byte(standData))
		compressLongData, _  = gcompress.Gzip([]byte(longData))
	)
	gtest.C(t, func(t *gtest.T) {

		t1, err := gtrace.SafeContentForHttp(compressShortData, gzipHeader)
		t.AssertNil(err)
		t.Assert(t1, shortData)

		t2, err := gtrace.SafeContentForHttp(compressStandData, gzipHeader)
		t.AssertNil(err)
		t.Assert(t2, standData)

		t3, err := gtrace.SafeContentForHttp(compressLongData, gzipHeader)
		t.AssertNil(err)
		t.Assert(t3, standData+"...")
	})
}
