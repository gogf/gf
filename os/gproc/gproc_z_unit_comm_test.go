// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_getCommPidSystemFolderPaths(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(getCommPidSystemFolderPaths("linux"), []string{"/var/tmp", "/var/run"})
		t.Assert(getCommPidSystemFolderPaths("darwin"), []string{"/var/tmp", "/var/run"})
		t.Assert(len(getCommPidSystemFolderPaths("windows")), 0)
	})
}

func Test_getCommPidFolderPath_NotInDriveRootOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows only")
	}
	gtest.C(t, func(t *gtest.T) {
		folderPath, err := getCommPidFolderPath()
		t.AssertNil(err)
		absPath, err := filepath.Abs(folderPath)
		t.AssertNil(err)
		// "/var/tmp" would resolve to "<drive>:\var\tmp" in the root of the current drive.
		t.Assert(strings.HasPrefix(absPath, filepath.VolumeName(absPath)+`\var\`), false)
	})
}
