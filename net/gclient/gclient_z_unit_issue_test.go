// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Issue3748(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write(
			r.GetBody(),
		)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	clientHost := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetHeader("Content-Type", "application/json")
		data := map[string]any{
			"name":  "@file:",
			"value": "json",
		}
		client.SetPrefix(clientHost)
		content := client.PostContent(ctx, "/", data)
		t.Assert(content, `{"name":"@file:","value":"json"}`)
	})

	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetHeader("Content-Type", "application/xml")
		data := map[string]any{
			"name":  "@file:",
			"value": "xml",
		}
		client.SetPrefix(clientHost)
		content := client.PostContent(ctx, "/", data)
		t.Assert(content, `<doc><name>@file:</name><value>xml</value></doc>`)
	})

	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		data := map[string]any{
			"name":  "@file:",
			"value": "x-www-form-urlencoded",
		}
		client.SetPrefix(clientHost)
		content := client.PostContent(ctx, "/", data)
		t.Assert(strings.Contains(content, `Content-Disposition: form-data; name="value"`), true)
		t.Assert(strings.Contains(content, `Content-Disposition: form-data; name="name"`), true)
		t.Assert(strings.Contains(content, "\r\n@file:"), true)
		t.Assert(strings.Contains(content, "\r\nx-www-form-urlencoded"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		data := "@file:"
		client.SetPrefix(clientHost)
		_, err := client.Post(ctx, "/", data)
		t.AssertNil(err)
	})
}

// https://github.com/gogf/gf/issues/4156
func Test_Issue4156(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/upload", func(r *ghttp.Request) {
		// Return the fieldName value received
		r.Response.Write(r.Get("fieldName"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	clientHost := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		// When posting form with file upload, if value contains '=', it should not be truncated.
		data := g.Map{
			"file":      "@file:" + gtest.DataPath("upload", "file1.txt"),
			"fieldName": "aaa=1&b=2",
		}
		content := client.PostContent(ctx, "/upload", data)
		// The complete value should be received, not truncated at '='
		t.Assert(content, "aaa=1&b=2")
	})
}

// Test_Issue4156_MultipleSpecialChars tests file upload with various special characters in field values.
func Test_Issue4156_MultipleSpecialChars(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/upload", func(r *ghttp.Request) {
		r.Response.Write(r.Get("field"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	clientHost := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	time.Sleep(100 * time.Millisecond)

	// Test with multiple equals signs
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":  "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field": "a=1=2=3",
		}
		content := client.PostContent(ctx, "/upload", data)
		t.Assert(content, "a=1=2=3")
	})

	// Test with multiple ampersands
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":  "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field": "a&b&c&d",
		}
		content := client.PostContent(ctx, "/upload", data)
		t.Assert(content, "a&b&c&d")
	})

	// Test with percent sign
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":  "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field": "100%complete",
		}
		content := client.PostContent(ctx, "/upload", data)
		t.Assert(content, "100%complete")
	})

	// Test with plus sign
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":  "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field": "1+2+3",
		}
		content := client.PostContent(ctx, "/upload", data)
		t.Assert(content, "1+2+3")
	})

	// Test with spaces
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":  "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field": "hello world test",
		}
		content := client.PostContent(ctx, "/upload", data)
		t.Assert(content, "hello world test")
	})

	// Test with mixed special characters
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":  "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field": "key=value&foo=bar%20test+plus",
		}
		content := client.PostContent(ctx, "/upload", data)
		t.Assert(content, "key=value&foo=bar%20test+plus")
	})
}

// Test_Issue4156_MultipleFields tests file upload with multiple fields containing special characters.
func Test_Issue4156_MultipleFields(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/upload", func(r *ghttp.Request) {
		// Return all field values as JSON-like format
		r.Response.Writef("field1=%s,field2=%s,field3=%s",
			r.Get("field1"), r.Get("field2"), r.Get("field3"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	clientHost := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":   "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field1": "a=1",
			"field2": "b&2",
			"field3": "c%3",
		}
		content := client.PostContent(ctx, "/upload", data)
		t.Assert(strings.Contains(content, "field1=a=1"), true)
		t.Assert(strings.Contains(content, "field2=b&2"), true)
		t.Assert(strings.Contains(content, "field3=c%3"), true)
	})
}

// Test_Issue4156_NoFileUpload tests that normal POST without file upload still works correctly.
func Test_Issue4156_NoFileUpload(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/post", func(r *ghttp.Request) {
		r.Response.Write(r.Get("field"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	clientHost := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	time.Sleep(100 * time.Millisecond)

	// Test normal POST with special characters (no file upload)
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"field": "a=1&b=2",
		}
		content := client.PostContent(ctx, "/post", data)
		t.Assert(content, "a=1&b=2")
	})

	// Test POST with Content-Type: application/x-www-form-urlencoded
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		client.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		data := g.Map{
			"field": "value=with=equals&and&ampersand",
		}
		content := client.PostContent(ctx, "/post", data)
		t.Assert(content, "value=with=equals&and&ampersand")
	})
}

// Test_Issue4156_PreEncodedValue tests that pre-encoded values are handled correctly.
func Test_Issue4156_PreEncodedValue(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/upload", func(r *ghttp.Request) {
		r.Response.Write(r.Get("field"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	clientHost := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	time.Sleep(100 * time.Millisecond)

	// Test with already URL-encoded value - should preserve the encoding
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":  "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field": "value%3Dwith%26encoding", // User wants to send literal %3D
		}
		content := client.PostContent(ctx, "/upload", data)
		// The literal %3D and %26 should be preserved
		t.Assert(content, "value%3Dwith%26encoding")
	})
}

// Test_Issue4156_EmptyAndSpecialValues tests edge cases with empty and special values.
func Test_Issue4156_EmptyAndSpecialValues(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/upload", func(r *ghttp.Request) {
		r.Response.Write(r.Get("field"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	clientHost := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	time.Sleep(100 * time.Millisecond)

	// Test with value starting with =
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":  "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field": "=startWithEquals",
		}
		content := client.PostContent(ctx, "/upload", data)
		t.Assert(content, "=startWithEquals")
	})

	// Test with value ending with =
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":  "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field": "endWithEquals=",
		}
		content := client.PostContent(ctx, "/upload", data)
		t.Assert(content, "endWithEquals=")
	})

	// Test with only special characters
	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetPrefix(clientHost)
		data := g.Map{
			"file":  "@file:" + gtest.DataPath("upload", "file1.txt"),
			"field": "=&=&=",
		}
		content := client.PostContent(ctx, "/upload", data)
		t.Assert(content, "=&=&=")
	})
}
