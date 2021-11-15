// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_LevelPrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := New()
		t.Assert(l.GetLevelPrefix(LEVEL_DEBU), defaultLevelPrefixes[LEVEL_DEBU])
		t.Assert(l.GetLevelPrefix(LEVEL_INFO), defaultLevelPrefixes[LEVEL_INFO])
		t.Assert(l.GetLevelPrefix(LEVEL_NOTI), defaultLevelPrefixes[LEVEL_NOTI])
		t.Assert(l.GetLevelPrefix(LEVEL_WARN), defaultLevelPrefixes[LEVEL_WARN])
		t.Assert(l.GetLevelPrefix(LEVEL_ERRO), defaultLevelPrefixes[LEVEL_ERRO])
		t.Assert(l.GetLevelPrefix(LEVEL_CRIT), defaultLevelPrefixes[LEVEL_CRIT])
		l.SetLevelPrefix(LEVEL_DEBU, "debug")
		t.Assert(l.GetLevelPrefix(LEVEL_DEBU), "debug")
		l.SetLevelPrefixes(map[int]string{
			LEVEL_CRIT: "critical",
		})
		t.Assert(l.GetLevelPrefix(LEVEL_DEBU), "debug")
		t.Assert(l.GetLevelPrefix(LEVEL_INFO), defaultLevelPrefixes[LEVEL_INFO])
		t.Assert(l.GetLevelPrefix(LEVEL_NOTI), defaultLevelPrefixes[LEVEL_NOTI])
		t.Assert(l.GetLevelPrefix(LEVEL_WARN), defaultLevelPrefixes[LEVEL_WARN])
		t.Assert(l.GetLevelPrefix(LEVEL_ERRO), defaultLevelPrefixes[LEVEL_ERRO])
		t.Assert(l.GetLevelPrefix(LEVEL_CRIT), "critical")
	})
	gtest.C(t, func(t *gtest.T) {
		buffer := bytes.NewBuffer(nil)
		l := New()
		l.SetWriter(buffer)
		l.Debug(ctx, "test1")
		t.Assert(gstr.Contains(buffer.String(), defaultLevelPrefixes[LEVEL_DEBU]), true)

		buffer.Reset()

		l.SetLevelPrefix(LEVEL_DEBU, "debug")
		l.Debug(ctx, "test2")
		t.Assert(gstr.Contains(buffer.String(), defaultLevelPrefixes[LEVEL_DEBU]), false)
		t.Assert(gstr.Contains(buffer.String(), "debug"), true)

		buffer.Reset()
		l.SetLevelPrefixes(map[int]string{
			LEVEL_ERRO: "error",
		})
		l.Error(ctx, "test3")
		t.Assert(gstr.Contains(buffer.String(), defaultLevelPrefixes[LEVEL_ERRO]), false)
		t.Assert(gstr.Contains(buffer.String(), "error"), true)
	})
}
