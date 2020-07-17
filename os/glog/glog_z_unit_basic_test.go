// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package glog

import (
	"bytes"
	"github.com/jin502437344/gf/test/gtest"
	"github.com/jin502437344/gf/text/gstr"
	"testing"
)

func Test_Print(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Print(1, 2, 3)
		l.Println(1, 2, 3)
		l.Printf("%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), "["), 0)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 3)
	})
}

func Test_Debug(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Debug(1, 2, 3)
		l.Debugf("%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_DEBU]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Info(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Info(1, 2, 3)
		l.Infof("%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_INFO]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Notice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Notice(1, 2, 3)
		l.Noticef("%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_NOTI]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Warning(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Warning(1, 2, 3)
		l.Warningf("%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_WARN]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Error(1, 2, 3)
		l.Errorf("%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}
