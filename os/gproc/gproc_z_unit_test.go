// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gproc_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	envKeyPPid = "GPROC_PPID"
)

func Test_ShellExec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s, err := gproc.ShellExec(gctx.New(), `echo 123`)
		t.AssertNil(err)
		t.Assert(s, "123\n")
	})
	// error
	gtest.C(t, func(t *gtest.T) {
		_, err := gproc.ShellExec(gctx.New(), `NoneExistCommandCall`)
		t.AssertNE(err, nil)
	})
}

func Test_Pid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(os.Getpid(), gproc.Pid())
	})
}

func Test_IsChild(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		originalPPid := os.Getenv(envKeyPPid)
		defer os.Setenv(envKeyPPid, originalPPid)

		os.Setenv(envKeyPPid, "1234")
		t.Assert(true, gproc.IsChild())

		os.Unsetenv(envKeyPPid)
		t.Assert(false, gproc.IsChild())
	})
}

func Test_SetPPid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := gproc.SetPPid(1234)
		t.AssertNil(err)
		t.Assert("1234", os.Getenv(envKeyPPid))

		err = gproc.SetPPid(0)
		t.AssertNil(err)
		t.Assert("", os.Getenv(envKeyPPid))
	})
}

func Test_StartTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		result := gproc.StartTime()
		t.Assert(result, gproc.StartTime())
	})
}

func Test_Uptime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		result := gproc.Uptime()
		t.AssertGT(result, 0)
	})
}

func Test_SearchBinary_FoundInPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "testbinary")
		gfile.Create(tempFile)
		os.Chmod(tempFile, 0755)

		originalPath := os.Getenv("PATH")
		os.Setenv("PATH", tempDir+string(os.PathListSeparator)+originalPath)
		defer os.Setenv("PATH", originalPath)

		result := gproc.SearchBinary("testbinary")
		t.Assert(result, tempFile)
	})
}

func Test_SearchBinary_NotFound(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		result := gproc.SearchBinary("nonexistentbinary")
		t.Assert(result, "")
	})
}

func Test_SearchBinaryPath_FoundInPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "testbinary")
		gfile.Create(tempFile)
		os.Chmod(tempFile, 0755)

		originalPath := os.Getenv("PATH")
		os.Setenv("PATH", tempDir+string(os.PathListSeparator)+originalPath)
		defer os.Setenv("PATH", originalPath)

		result := gproc.SearchBinaryPath("testbinary")
		t.Assert(result, tempFile)
	})
}

func Test_SearchBinaryPath_NotFound(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		result := gproc.SearchBinaryPath("nonexistentbinary")
		t.Assert(result, "")
	})
}

func Test_PPidOS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ppid := gproc.PPidOS()
		expectedPpid := os.Getppid()
		t.Assert(ppid, expectedPpid)
	})
}

func Test_PPid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		customPPid := 12345
		os.Setenv("GPROC_PPID", gconv.String(customPPid))
		defer os.Unsetenv("GPROC_PPID")

		t.Assert(gproc.PPid(), customPPid)
	})
}
