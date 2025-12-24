// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"strings"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_cRunApp_getWatchPaths_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			WatchPaths: []string{"."},
		}
		watchPaths := app.getWatchPaths()

		t.AssertGT(len(watchPaths), 0)
		for _, v := range watchPaths {
			t.Log(v)
		}
	})
}

func Test_cRunApp_getWatchPaths_EmptyWatchPaths(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			WatchPaths: []string{},
		}
		watchPaths := app.getWatchPaths()

		// Should default to current directory "."
		t.AssertGT(len(watchPaths), 0)
	})
}

func Test_cRunApp_getWatchPaths_CustomIgnorePattern(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := &cRunApp{
			WatchPaths:     []string{"testdata"},
			IgnorePatterns: []string{"2572"},
		}
		watchPaths := app.getWatchPaths()

		// Ensure the "2572" directory is not watched directly.
		for _, path := range watchPaths {
			t.Log("watch path:", path)
			t.AssertNE(true, strings.HasSuffix(path, "2572"))
		}
		t.AssertGT(len(watchPaths), 0)
	})
}
