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
		compressShortData []byte
		compressStandData []byte
		compressLongData  []byte
	)
	gtest.C(t, func(t *gtest.T) {
		var err error
		compressShortData, err = gcompress.Gzip([]byte(shortData))
		t.AssertNil(err)
		compressStandData, err = gcompress.Gzip([]byte(standData))
		t.AssertNil(err)
		compressLongData, err = gcompress.Gzip([]byte(longData))
		t.AssertNil(err)
	})
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

func Test_SafeContentForHttp_GzipEncoded(t *testing.T) {
	var (
		content        = strings.Repeat("gzip content ", 10)
		compressedData []byte
	)
	gtest.C(t, func(t *gtest.T) {
		var err error
		compressedData, err = gcompress.Gzip([]byte(content))
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		for _, c := range []struct {
			header http.Header
			data   []byte
			expect string
		}{
			// Content codings are case-insensitive.
			{http.Header{"Content-Encoding": {"gzip"}}, compressedData, content},
			{http.Header{"Content-Encoding": {"GZip"}}, compressedData, content},
			{http.Header{"Content-Encoding": {"GZIP"}}, compressedData, content},
			// Parameters and surrounding blanks of the content coding are ignored.
			{http.Header{"Content-Encoding": {" gzip "}}, compressedData, content},
			{http.Header{"Content-Encoding": {"gzip; level=1"}}, compressedData, content},
			{http.Header{"Content-Encoding": {"GZip ; level=1"}}, compressedData, content},
			{http.Header{"Content-Encoding": {"gzip; foo"}}, compressedData, content},
			// Every coding of a comma separated value is checked.
			{http.Header{"Content-Encoding": {"br, gzip"}}, compressedData, content},
			// Every Content-Encoding header line is checked.
			{http.Header{"Content-Encoding": {"br", "gzip"}}, compressedData, content},
			// Other content codings are not decompressed.
			{http.Header{}, []byte(content), content},
			{http.Header{"Content-Encoding": {""}}, []byte(content), content},
			{http.Header{"Content-Encoding": {"identity"}}, []byte(content), content},
			{http.Header{"Content-Encoding": {"br"}}, []byte(content), content},
			// `x-gzip` is an alias of `gzip` in RFC 9110, but it is not handled on purpose.
			{http.Header{"Content-Encoding": {"x-gzip"}}, []byte(content), content},
			// The coding name is compared as a whole, not by prefix.
			{http.Header{"Content-Encoding": {"gzipx"}}, []byte(content), content},
		} {
			s, err := gtrace.SafeContentForHttp(c.data, c.header)
			t.AssertNil(err)
			t.Assert(s, c.expect)
		}
	})

	// The error is returned if the content is declared gzip-encoded but is not.
	gtest.C(t, func(t *gtest.T) {
		s, err := gtrace.SafeContentForHttp([]byte(content), http.Header{"Content-Encoding": {"gzip"}})
		t.AssertNE(err, nil)
		t.Assert(s, "")
	})
}
