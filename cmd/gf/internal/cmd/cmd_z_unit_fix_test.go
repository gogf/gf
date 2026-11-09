// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Fix_doFixV25Content(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			content = gtest.DataContent(`fix`, `fix25_content.go`)
			f       = cFix{}
		)
		_, err := f.doFixV25Content(content)
		t.AssertNil(err)
	})
}

func Test_Fix_doFixV25Content_WithReplacement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			f       = cFix{}
			content = `s.BindHookHandlerByMap("/path", map[string]ghttp.HandlerFunc{
				ghttp.HookBeforeServe: func(r *ghttp.Request) {},
			})`
		)
		newContent, err := f.doFixV25Content(content)
		t.AssertNil(err)
		// Verify the replacement was made
		t.Assert(gstr.Contains(newContent, "map[ghttp.HookName]ghttp.HandlerFunc"), true)
		t.Assert(gstr.Contains(newContent, "map[string]ghttp.HandlerFunc"), false)
	})
}

func Test_Fix_doFixV25Content_NoMatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			f       = cFix{}
			content = `package main

func main() {
	fmt.Println("Hello World")
}
`
		)
		newContent, err := f.doFixV25Content(content)
		t.AssertNil(err)
		// Content should remain unchanged
		t.Assert(newContent, content)
	})
}

func Test_Fix_doFixV25Content_MultipleMatches(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			f       = cFix{}
			content = `
s.BindHookHandlerByMap("/path1", map[string]ghttp.HandlerFunc{})
s.BindHookHandlerByMap("/path2", map[string]ghttp.HandlerFunc{})
`
		)
		newContent, err := f.doFixV25Content(content)
		t.AssertNil(err)
		// Both should be replaced
		count := gstr.Count(newContent, "map[ghttp.HookName]ghttp.HandlerFunc")
		t.Assert(count, 2)
	})
}

func Test_Fix_doFixV25Content_EmptyContent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			f       = cFix{}
			content = ""
		)
		newContent, err := f.doFixV25Content(content)
		t.AssertNil(err)
		t.Assert(newContent, "")
	})
}

func Test_Fix_doFixV25Content_ComplexPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			f       = cFix{}
			content = `s.BindHookHandlerByMap("/api/v1/user/{id}/profile", map[string]ghttp.HandlerFunc{
				ghttp.HookBeforeServe: func(r *ghttp.Request) {
					r.Response.Write("before")
				},
			})`
		)
		newContent, err := f.doFixV25Content(content)
		t.AssertNil(err)
		t.Assert(gstr.Contains(newContent, "map[ghttp.HookName]ghttp.HandlerFunc"), true)
	})
}
