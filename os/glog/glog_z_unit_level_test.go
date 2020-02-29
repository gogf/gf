// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/text/gstr"
	"testing"
)

func Test_LevelPrefix(t *testing.T) {
	gtest.Case(t, func() {
		l := New()
		gtest.Assert(l.GetLevelPrefix(LEVEL_DEBU), defaultLevelPrefixes[LEVEL_DEBU])
		gtest.Assert(l.GetLevelPrefix(LEVEL_INFO), defaultLevelPrefixes[LEVEL_INFO])
		gtest.Assert(l.GetLevelPrefix(LEVEL_NOTI), defaultLevelPrefixes[LEVEL_NOTI])
		gtest.Assert(l.GetLevelPrefix(LEVEL_WARN), defaultLevelPrefixes[LEVEL_WARN])
		gtest.Assert(l.GetLevelPrefix(LEVEL_ERRO), defaultLevelPrefixes[LEVEL_ERRO])
		gtest.Assert(l.GetLevelPrefix(LEVEL_CRIT), defaultLevelPrefixes[LEVEL_CRIT])
		l.SetLevelPrefix(LEVEL_DEBU, "debug")
		gtest.Assert(l.GetLevelPrefix(LEVEL_DEBU), "debug")
		l.SetLevelPrefixes(map[int]string{
			LEVEL_CRIT: "critical",
		})
		gtest.Assert(l.GetLevelPrefix(LEVEL_DEBU), "debug")
		gtest.Assert(l.GetLevelPrefix(LEVEL_INFO), defaultLevelPrefixes[LEVEL_INFO])
		gtest.Assert(l.GetLevelPrefix(LEVEL_NOTI), defaultLevelPrefixes[LEVEL_NOTI])
		gtest.Assert(l.GetLevelPrefix(LEVEL_WARN), defaultLevelPrefixes[LEVEL_WARN])
		gtest.Assert(l.GetLevelPrefix(LEVEL_ERRO), defaultLevelPrefixes[LEVEL_ERRO])
		gtest.Assert(l.GetLevelPrefix(LEVEL_CRIT), "critical")
	})
	gtest.Case(t, func() {
		buffer := bytes.NewBuffer(nil)
		l := New()
		l.SetWriter(buffer)
		l.Debug("test1")
		gtest.Assert(gstr.Contains(buffer.String(), defaultLevelPrefixes[LEVEL_DEBU]), true)

		buffer.Reset()

		l.SetLevelPrefix(LEVEL_DEBU, "debug")
		l.Debug("test2")
		gtest.Assert(gstr.Contains(buffer.String(), defaultLevelPrefixes[LEVEL_DEBU]), false)
		gtest.Assert(gstr.Contains(buffer.String(), "debug"), true)

		buffer.Reset()
		l.SetLevelPrefixes(map[int]string{
			LEVEL_ERRO: "error",
		})
		l.Error("test3")
		gtest.Assert(gstr.Contains(buffer.String(), defaultLevelPrefixes[LEVEL_ERRO]), false)
		gtest.Assert(gstr.Contains(buffer.String(), "error"), true)
	})
}
