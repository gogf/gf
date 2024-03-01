// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtrace_test

import (
	"strings"
	"testing"

	"github.com/gogf/gf/v2/encoding/gcompress"

	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSafeContent(t *testing.T) {
	var (
		defText   = "中"
		shortData = strings.Repeat(defText, gtrace.MaxContentLogSize()-1)
		standData = strings.Repeat(defText, gtrace.MaxContentLogSize())
		longData  = strings.Repeat(defText, gtrace.MaxContentLogSize()+1)
	)
	gtest.C(t, func(t *gtest.T) {

		compressedContent, err := gcompress.Gzip([]byte(defText))
		t.AssertNil(err)
		t.AssertEQ(gtrace.IsGzipped(compressedContent), true)
	})

	// safe content
	gtest.C(t, func(t *gtest.T) {

		t1, err := gtrace.SafeContent([]byte(shortData))
		t.AssertNil(err)
		t.Assert(t1, shortData)

		t2, err := gtrace.SafeContent([]byte(standData))
		t.AssertNil(err)
		t.Assert(t2, standData)

		t3, err := gtrace.SafeContent([]byte(longData))
		t.AssertNil(err)
		t.Assert(t3, standData+"...")
	})

	// compress content
	var (
		compressShortData, _ = gcompress.Gzip([]byte(shortData))
		compressStandData, _ = gcompress.Gzip([]byte(standData))
		compressLongData, _  = gcompress.Gzip([]byte(longData))
	)
	gtest.C(t, func(t *gtest.T) {

		t1, err := gtrace.SafeContent(compressShortData)
		t.AssertNil(err)
		t.Assert(t1, shortData)

		t2, err := gtrace.SafeContent(compressStandData)
		t.AssertNil(err)
		t.Assert(t2, standData)

		t3, err := gtrace.SafeContent(compressLongData)
		t.AssertNil(err)
		t.Assert(t3, standData+"...")
	})
}
