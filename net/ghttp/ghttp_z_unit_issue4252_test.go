// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Regression tests for upload slice types passed to request validation.

package ghttp

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gvalid"
)

// newIssue4252Request creates an in-memory multipart request without saving files.
func newIssue4252Request(t testing.TB, name string, count int) *Request {
	t.Helper()
	var (
		body   bytes.Buffer
		writer = multipart.NewWriter(&body)
	)
	for i := 0; i < count; i++ {
		part, err := writer.CreateFormFile(name, fmt.Sprintf("file%d.txt", i))
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write([]byte("upload")); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	var req = httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set(HeaderContentType, writer.FormDataContentType())
	t.Cleanup(func() {
		if req.MultipartForm != nil {
			if err := req.MultipartForm.RemoveAll(); err != nil {
				t.Error(err)
			}
		}
	})
	return &Request{Request: req, Server: &Server{config: NewConfig()}}
}

// issue4252NamedFiles checks user-defined slice types as well as UploadFiles.
type issue4252NamedFiles []*UploadFile

// Issue4252Embedded supplies an exported embedded upload field.
type Issue4252Embedded struct {
	Files UploadFiles `p:"files" v:"array|issue4252-declared-type"`
}

// Test_Issue4252_FieldTypes covers binding names, slice types and embedded fields.
func Test_Issue4252_FieldTypes(t *testing.T) {
	var expectedType reflect.Type
	gvalid.RegisterRule("issue4252-declared-type", func(_ context.Context, in gvalid.RuleFuncInput) error {
		if reflect.TypeOf(in.Value.Interface()) != expectedType {
			return gerror.Newf("expected %v, got %T", expectedType, in.Value.Interface())
		}
		return nil
	})
	for _, tc := range []struct {
		name  string
		field string
		req   any
	}{
		{"slice", "files", &struct {
			Files []*UploadFile `p:"files" v:"array|issue4252-declared-type"`
		}{}},
		{"named", "files", &struct {
			Files issue4252NamedFiles `p:"files" v:"array|issue4252-declared-type"`
		}{}},
		{"json tag", "uploads", &struct {
			Files UploadFiles `json:"uploads,omitempty" v:"array|issue4252-declared-type"`
		}{}},
		{"field name", "Files", &struct {
			Files UploadFiles `p:"uploads" v:"array|issue4252-declared-type"`
		}{}},
		{"fuzzy field name", "f_i_l_e_s", &struct {
			Files UploadFiles `p:"uploads" v:"array|issue4252-declared-type"`
		}{}},
		{"embedded", "files", &struct{ Issue4252Embedded }{}},
		{"pointer to request", "files", new(*struct {
			Files UploadFiles `p:"files" v:"array|issue4252-declared-type"`
		})},
		{"single file", "files", &struct {
			Files *UploadFile `p:"files" v:"issue4252-declared-type"`
		}{}},
	} {
		reqType := reflect.TypeOf(tc.req)
		for reqType.Kind() == reflect.Pointer {
			reqType = reqType.Elem()
		}
		field, ok := reqType.FieldByName("Files")
		if !ok {
			t.Fatal("missing Files field")
		}
		expectedType = field.Type
		for _, count := range []int{1, 2} {
			if tc.name == "single file" && count > 1 {
				continue
			}
			t.Run(fmt.Sprintf("%s/%d", tc.name, count), func(t *testing.T) {
				gtest.C(t, func(t *gtest.T) {
					var r = newIssue4252Request(t, tc.field, count)
					t.AssertNil(r.Parse(tc.req))
					value := reflect.ValueOf(tc.req).Elem()
					if value.Kind() == reflect.Pointer {
						value = value.Elem()
					}
					field := value.FieldByName("Files")
					if field.Kind() == reflect.Slice {
						t.Assert(field.Len(), count)
					}
				})
			})
		}
	}
}

// Test_Issue4252_LengthRules checks missing files and upload count constraints.
func Test_Issue4252_LengthRules(t *testing.T) {
	for _, count := range []int{0, 1, 2, 3} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					r   = newIssue4252Request(t, "files", count)
					req struct {
						Files UploadFiles `p:"files" v:"required|array|min-length:1|max-length:2"`
						Note  string      `d:"default note"`
					}
				)
				err := r.Parse(&req)
				if count == 0 || count > 2 {
					t.AssertNE(err, nil)
				} else {
					t.AssertNil(err)
				}
				t.Assert(len(req.Files), count)
				t.Assert(req.Note, "default note")
			})
		})
	}
}

// Test_Issue4252_NonUploadValues keeps ordinary parameters outside upload normalization.
func Test_Issue4252_NonUploadValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			r    = newIssue4252Request(t, "other", 1)
			req  struct{ Files UploadFiles }
			data = r.GetRequestMap()
		)
		data["Files"] = "unchanged"
		t.AssertNil(r.mergeUploadFiles(data, &req))
		t.AssertEQ(data["Files"], "unchanged")
		t.AssertNE(r.mergeUploadFiles(data, 1), nil)
	})
}

// Test_Issue4252_ExplicitMapping keeps parameter binding independent of validation normalization.
func Test_Issue4252_ExplicitMapping(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			r   = newIssue4252Request(t, "files", 1)
			req struct {
				Files  UploadFiles
				Single *UploadFile
			}
		)
		t.AssertNil(r.GetRequestStruct(&req, map[string]string{"files": "Single"}))
		t.AssertNE(req.Single, nil)
		t.Assert(req.Single.Filename, "file0.txt")
		t.Assert(len(req.Files), 0)
	})
}

// Test_Issue4252 checks that custom validation sees the declared upload slice type.
func Test_Issue4252(t *testing.T) {
	gvalid.RegisterRule("issue4252-upload-files", func(_ context.Context, in gvalid.RuleFuncInput) error {
		if _, ok := in.Value.Interface().(UploadFiles); !ok {
			return gerror.Newf("expected UploadFiles, got %T", in.Value.Interface())
		}
		return nil
	})
	for _, count := range []int{1, 2} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					r   = newIssue4252Request(t, "files", count)
					req struct {
						Files UploadFiles `p:"files" v:"issue4252-upload-files"`
					}
				)
				t.AssertNil(r.Parse(&req))
				t.Assert(len(req.Files), count)
				// Public parameter retrieval retains its single-file/multiple-file contract.
				if count == 1 {
					_, ok := r.GetRequestMap()["files"].(*UploadFile)
					t.Assert(ok, true)
				} else {
					_, ok := r.GetRequestMap()["files"].(UploadFiles)
					t.Assert(ok, true)
				}
			})
		})
	}
}
