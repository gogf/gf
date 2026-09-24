// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT License was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// These tests cover multipart lifetime and error recovery during parameter reload.

package ghttp_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// reloadMultipartBody builds an upload that exceeds the test server's memory limit.
func reloadMultipartBody(t *testing.T, field string) (*bytes.Buffer, string) {
	t.Helper()
	var (
		body   = new(bytes.Buffer)
		writer = multipart.NewWriter(body)
	)
	if field != "" {
		if err := writer.WriteField("field", field); err != nil {
			t.Fatal(err)
		}
	}
	part, err := writer.CreateFormFile("upload", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = io.WriteString(part, "payload"); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}
	return body, writer.FormDataContentType()
}

// reloadMultipartRequest creates a request with a parsed, disk-backed upload.
func reloadMultipartRequest(t *testing.T, field string) *ghttp.Request {
	t.Helper()
	var (
		body, contentType = reloadMultipartBody(t, field)
		server            = g.Server(guid.S())
		r                 = &ghttp.Request{
			Request: httptest.NewRequest(http.MethodPost, "/?old=1", body),
			Server:  server,
		}
	)
	server.SetFormParsingMemory(1)
	r.Header.Set("Content-Type", contentType)
	r.GetFormMap()
	form := r.MultipartForm
	if form == nil || len(form.File["upload"]) != 1 {
		t.Fatal("expected a parsed upload")
	}
	t.Cleanup(func() {
		if err := form.RemoveAll(); err != nil {
			t.Error(err)
		}
		if r.MultipartForm != nil && r.MultipartForm != form {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				t.Error(err)
			}
		}
	})
	return r
}

// Test_Request_ReloadParam_UnchangedMultipart preserves uploads across repeated reloads.
func Test_Request_ReloadParam_UnchangedMultipart(t *testing.T) {
	for _, field := range []string{"kept", ""} {
		t.Run("field="+field, func(t *testing.T) {
			r := reloadMultipartRequest(t, field)
			gtest.C(t, func(t *gtest.T) {
				form := r.MultipartForm
				for i := 0; i < 3; i++ {
					r.GetRequestMap()
					r.MakeBodyRepeatableRead(true)
					r.URL.RawQuery = "new=2"
					r.ReloadParam()
					r.GetRequestMap()
					t.Assert(r.GetQueryMap(), g.Map{"new": "2"})
					t.Assert(r.MultipartForm == form, true)
					t.Assert(r.FormValue("field"), field)
					if field != "" {
						t.Assert(r.GetFormMap()["field"], field)
					}
					file, err := form.File["upload"][0].Open()
					t.AssertNil(err)
					content, err := io.ReadAll(file)
					closeErr := file.Close()
					t.AssertNil(err)
					t.AssertNil(closeErr)
					t.Assert(string(content), "payload")
				}
			})
		})
	}
}

// Test_Request_ReloadParam_StandardMultipart covers parsing through embedded HTTP methods.
func Test_Request_ReloadParam_StandardMultipart(t *testing.T) {
	middleBody, middleContentType := reloadMultipartBody(t, "mid")
	s := g.Server(guid.S())
	s.BindHandler("/reload/standard", func(r *ghttp.Request) {
		r.ReloadParam()
		for _, replacement := range []struct{ body, contentType string }{
			{middleBody.String(), middleContentType},
			{`{"field":"new"}`, "application/json"},
		} {
			// Force disk-backed uploads without invoking framework form parsing.
			if err := r.ParseMultipartForm(1); err != nil {
				r.Response.WriteExit(err)
			}
			file, upload, err := r.FormFile("upload")
			if err != nil {
				r.Response.WriteExit(err)
			}
			if err = file.Close(); err != nil {
				r.Response.WriteExit(err)
			}
			r.Body = io.NopCloser(bytes.NewBufferString(replacement.body))
			r.ContentLength = int64(len(replacement.body))
			r.Header.Set("Content-Type", replacement.contentType)
			r.ReloadParam()
			file, err = upload.Open()
			if err == nil {
				if closeErr := file.Close(); closeErr != nil {
					r.Response.WriteExit(closeErr)
				}
			}
			if !os.IsNotExist(err) {
				r.Response.WriteExit("old upload was not removed")
			}
		}
		r.Response.WriteJson(r.GetRequestMap())
	})
	s.SetDumpRouterMap(false)
	startArrayServer(t, s)
	body, contentType := reloadMultipartBody(t, "old")
	gtest.C(t, func(t *gtest.T) {
		response := g.Client().ContentType(contentType).PostContent(
			ctx, fmt.Sprintf("http://127.0.0.1:%d/reload/standard", s.GetListenedPort()), body.String(),
		)
		t.Assert(response, `{"field":"new"}`)
	})
}

// reloadBodyValue intentionally has an incomparable dynamic type to check Body tracking.
type reloadBodyValue struct {
	io.ReadCloser
	padding []byte
}

// Test_Request_ReloadParam_ReplacedMultipart removes the old upload only on replacement.
func Test_Request_ReloadParam_ReplacedMultipart(t *testing.T) {
	for _, multipartReplacement := range []bool{false, true} {
		name := "json"
		if multipartReplacement {
			name = "multipart"
		}
		t.Run(name, func(t *testing.T) {
			r := reloadMultipartRequest(t, "old")
			upload := r.MultipartForm.File["upload"][0]
			var (
				body        = bytes.NewBufferString(`{"field":"new"}`)
				contentType = "application/json"
			)
			if multipartReplacement {
				body, contentType = reloadMultipartBody(t, "new")
			}
			r.Body = reloadBodyValue{ReadCloser: io.NopCloser(body)}
			r.ContentLength = int64(body.Len())
			r.Header.Set("Content-Type", contentType)
			r.ReloadParam()
			gtest.C(t, func(t *gtest.T) {
				file, err := upload.Open()
				if err == nil {
					t.AssertNil(file.Close())
				}
				t.Assert(os.IsNotExist(err), true)
				t.Assert(r.GetRequestMap()["field"], "new")
				t.AssertNil(r.GetError())
				if multipartReplacement {
					r.ReloadParam()
					t.Assert(r.GetRequestMap()["field"], "new")
				}
			})
		})
	}
}

// Test_Request_ReloadParam_ErrorState replaces stale errors only on full reload.
func Test_Request_ReloadParam_ErrorState(t *testing.T) {
	t.Run("explicit request error", func(t *testing.T) {
		gtest.C(t, func(t *gtest.T) {
			r := &ghttp.Request{Request: httptest.NewRequest(http.MethodGet, "/", nil)}
			err := gerror.New("middleware error")
			r.SetError(err)
			r.GetRequestMap()
			t.Assert(r.GetError() == err, true)
			r.ReloadQuery()
			t.Assert(r.GetError() == err, true)
			r.ReloadParam()
			t.AssertNil(r.GetError())
		})
	})
	for _, body := range []string{`{"items":[{"id":2}]}`, `[{"id":2}]`, `[]`, `{"new":}`} {
		t.Run(body, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				r := &ghttp.Request{
					Request: httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"old":}`)),
					Server:  g.Server(guid.S()),
				}
				r.Header.Set("Content-Type", "application/json")
				r.GetRequestMap()
				oldError := r.GetError()
				t.Assert(gerror.Code(oldError), gcode.CodeInvalidParameter)
				r.ReloadQuery()
				t.Assert(r.GetError() == oldError, true)
				r.Body = io.NopCloser(bytes.NewBufferString(body))
				r.ContentLength = int64(len(body))
				r.ReloadParam()
				t.AssertNil(r.GetError())
				r.GetRequestMap()
				if body == `{"new":}` {
					t.Assert(gerror.Code(r.GetError()), gcode.CodeInvalidParameter)
					t.Assert(r.GetError() == oldError, false)
				} else {
					t.AssertNil(r.GetError())
				}
			})
		})
	}
}
