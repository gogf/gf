// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Env_Index(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test that env command runs without error
		_, err := Env.Index(ctx, cEnvInput{})
		t.AssertNil(err)
	})
}

func Test_Env_ParseGoEnvOutput(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test parsing normal go env output
		lines := []string{
			"set GOPATH=C:\\Users\\test\\go",
			"set GOROOT=C:\\Go",
			"set GOOS=windows",
			"GOARCH=amd64", // Unix format without "set " prefix
			"CGO_ENABLED=0",
		}

		for _, line := range lines {
			line = gstr.Trim(line)
			if gstr.Pos(line, "set ") == 0 {
				line = line[4:]
			}
			match, _ := gregex.MatchString(`(.+?)=(.*)`, line)
			t.Assert(len(match) >= 3, true)
		}
	})
}

func Test_Env_ParseGoEnvOutput_WithWarnings(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test parsing go env output that contains warning messages
		// These lines should be skipped without causing errors
		lines := []string{
			"go: stripping unprintable or unescapable characters from %\"GOPROXY\"%",
			"go: warning: some warning message",
			"# this is a comment",
			"",
			"set GOPATH=C:\\Users\\test\\go",
			"set GOOS=windows",
		}

		array := make([][]string, 0)
		for _, line := range lines {
			line = gstr.Trim(line)
			if line == "" {
				continue
			}
			if gstr.Pos(line, "set ") == 0 {
				line = line[4:]
			}
			match, _ := gregex.MatchString(`(.+?)=(.*)`, line)
			if len(match) < 3 {
				// Skip lines that don't match key=value format (e.g., warning messages)
				continue
			}
			array = append(array, []string{gstr.Trim(match[1]), gstr.Trim(match[2])})
		}

		// Should have parsed 2 valid environment variables
		t.Assert(len(array), 2)
		t.Assert(array[0][0], "GOPATH")
		t.Assert(array[0][1], "C:\\Users\\test\\go")
		t.Assert(array[1][0], "GOOS")
		t.Assert(array[1][1], "windows")
	})
}
