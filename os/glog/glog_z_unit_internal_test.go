// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	ctx = context.TODO()
)

const (
	// fatalAsyncHelperEnv marks the child process that calls Fatal/Fatalf on an async logger.
	fatalAsyncHelperEnv = "GF_GLOG_FATAL_ASYNC_HELPER"
	// fatalAsyncHelperNotExitedCode is the child process exit code used when Fatal/Fatalf
	// return without exiting the process, which fails the parent assertion.
	fatalAsyncHelperNotExitedCode = 3
)

func Test_Print(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Print(ctx, 1, 2, 3)
		l.Printf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), "["), 0)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Debug(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Debug(ctx, 1, 2, 3)
		l.Debugf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_DEBU]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Info(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Info(ctx, 1, 2, 3)
		l.Infof(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_INFO]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Notice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Notice(ctx, 1, 2, 3)
		l.Noticef(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_NOTI]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Warning(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Warning(ctx, 1, 2, 3)
		l.Warningf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_WARN]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Error(ctx, 1, 2, 3)
		l.Errorf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gstr.Count(w.String(), defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 2)
	})
}

// Test_Fatal_Async checks that Fatal/Fatalf still print their content when the logger
// is async. They exit the process, so they are called in a child process.
func Test_Fatal_Async(t *testing.T) {
	if method := os.Getenv(fatalAsyncHelperEnv); method != "" {
		l := New()
		l.SetAsync(true)
		l.SetStack(false)
		switch method {
		case "Fatalf":
			l.Fatalf(ctx, "%d %d %d", 1, 2, 3)
		default:
			l.Fatal(ctx, 1, 2, 3)
		}
		// Fatal/Fatalf always exit the process, so reaching here means they did not,
		// which must be reported to the parent assertion instead of inlining the test.
		os.Exit(fatalAsyncHelperNotExitedCode)
	}
	gtest.C(t, func(t *gtest.T) {
		for _, method := range []string{"Fatal", "Fatalf"} {
			cmd := exec.Command(os.Args[0], "-test.run=^Test_Fatal_Async$")
			cmd.Env = append(os.Environ(), fatalAsyncHelperEnv+"="+method)
			output, err := cmd.CombinedOutput()
			t.AssertNE(err, nil)
			t.Assert(cmd.ProcessState.ExitCode(), 1)
			t.Assert(gstr.Count(string(output), defaultLevelPrefixes[LEVEL_FATA]), 1)
			t.Assert(gstr.Count(string(output), "1 2 3"), 1)
		}
	})
}

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
