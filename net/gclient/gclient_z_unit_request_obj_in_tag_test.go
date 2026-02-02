// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_DoRequestObj_InTag_Mixed tests mixed parameter types with in tag
func Test_DoRequestObj_InTag_Mixed(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/user/{id}", func(r *ghttp.Request) {
		// Verify path parameter
		pathId := r.Get("id").String()

		// Verify query parameter
		queryPage := r.URL.Query().Get("page")

		// Verify header parameter
		headerToken := r.Header.Get("Authorization")

		// Verify cookie parameter
		cookieSession := r.Cookie.Get("session")

		// Verify body parameters
		bodyMap := r.GetBodyMap()
		bodyName := gconv.String(bodyMap["name"])
		bodyAge := gconv.Int(bodyMap["age"])

		// Return verification result
		r.Response.Writef("path_id=%s,query_page=%s,header_token=%s,cookie_session=%s,body_name=%s,body_age=%d",
			pathId, queryPage, headerToken, cookieSession, bodyName, bodyAge)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			g.Meta  `path:"/user/{id}" method:"post"`
			Id      int    `in:"path"`
			Page    int    `in:"query" json:"page"`
			Token   string `in:"header" json:"Authorization"`
			Session string `in:"cookie" json:"session"`
			Name    string `json:"name"`
			Age     int    `json:"age"`
		}

		var res string
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			DoRequestObj(context.Background(), &Req{
				Id:      123,
				Page:    1,
				Token:   "Bearer xxx",
				Session: "session-id",
				Name:    "john",
				Age:     25,
			}, &res)

		t.AssertNil(err)
		// Verify each parameter is in the correct location
		t.Assert(res, "path_id=123,query_page=1,header_token=Bearer xxx,cookie_session=session-id,body_name=john,body_age=25")
	})
}

// Test_DoRequestObj_InTag_QuerySlice tests slice query parameters
func Test_DoRequestObj_InTag_QuerySlice(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/filter", func(r *ghttp.Request) {
		// Get slice parameters from URL query (not from r.Get which only returns first value)
		ids := r.URL.Query()["ids"]
		tags := r.URL.Query()["tags"]

		// Verify we got all values
		r.Response.Writef("ids_count=%d,ids_values=%v,tags_count=%d,tags_values=%v",
			len(ids), ids, len(tags), tags)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			g.Meta `path:"/filter" method:"get"`
			Ids    []int    `in:"query" json:"ids"`
			Tags   []string `in:"query" json:"tags"`
		}

		var res string
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			DoRequestObj(context.Background(), &Req{
				Ids:  []int{1, 2, 3},
				Tags: []string{"go", "web"},
			}, &res)

		t.AssertNil(err)
		// Verify all slice values are sent
		t.Assert(res, "ids_count=3,ids_values=[1 2 3],tags_count=2,tags_values=[go web]")
	})
}

// Test_DoRequestObj_InTag_QueryMap tests map query parameters
func Test_DoRequestObj_InTag_QueryMap(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/filter", func(r *ghttp.Request) {
		// Verify map parameters are flattened to filter[key] format in query
		name := r.URL.Query().Get("filter[name]")
		age := r.URL.Query().Get("filter[age]")

		// Verify they are NOT in the body
		bodyContent := string(r.GetBody())

		r.Response.Writef("query_name=%s,query_age=%s,body_empty=%v",
			name, age, bodyContent == "" || bodyContent == "{}")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			g.Meta `path:"/filter" method:"get"`
			Filter map[string]string `in:"query" json:"filter"`
		}

		var res string
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			DoRequestObj(context.Background(), &Req{
				Filter: map[string]string{"name": "john", "age": "25"},
			}, &res)

		t.AssertNil(err)
		// Verify map is flattened to query parameters, not in body
		t.Assert(res, "query_name=john,query_age=25,body_empty=true")
	})
}

// Test_DoRequestObj_InTag_EmbeddedStruct tests embedded struct
func Test_DoRequestObj_InTag_EmbeddedStruct(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/list", func(r *ghttp.Request) {
		// Verify embedded struct fields are flattened to query
		queryPage := r.URL.Query().Get("page")
		querySize := r.URL.Query().Get("size")
		queryKeyword := r.URL.Query().Get("keyword")

		// Verify they are NOT in body
		bodyContent := string(r.GetBody())

		r.Response.Writef("query_page=%s,query_size=%s,query_keyword=%s,body_empty=%v",
			queryPage, querySize, queryKeyword, bodyContent == "" || bodyContent == "{}")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type Pagination struct {
			Page int `in:"query" json:"page"`
			Size int `in:"query" json:"size"`
		}

		type Req struct {
			g.Meta `path:"/list" method:"get"`
			Pagination
			Keyword string `in:"query" json:"keyword"`
		}

		var res string
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			DoRequestObj(context.Background(), &Req{
				Pagination: Pagination{Page: 1, Size: 20},
				Keyword:    "golang",
			}, &res)

		t.AssertNil(err)
		// Verify embedded struct fields are flattened to query
		t.Assert(res, "query_page=1,query_size=20,query_keyword=golang,body_empty=true")
	})
}

// Test_DoRequestObj_InTag_NamedStructQuery tests named struct with in:"query"
func Test_DoRequestObj_InTag_NamedStructQuery(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/list", func(r *ghttp.Request) {
		// Verify named struct fields with in:"query" are flattened to query
		queryPage := r.URL.Query().Get("page")
		querySize := r.URL.Query().Get("size")
		queryKeyword := r.URL.Query().Get("keyword")

		// Verify they are NOT in body
		bodyContent := string(r.GetBody())

		r.Response.Writef("query_page=%s,query_size=%s,query_keyword=%s,body_empty=%v",
			queryPage, querySize, queryKeyword, bodyContent == "" || bodyContent == "{}")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type Pagination struct {
			Page int `json:"page"`
			Size int `json:"size"`
		}

		type Req struct {
			g.Meta     `path:"/list" method:"get"`
			Pagination Pagination `in:"query"`
			Keyword    string     `in:"query" json:"keyword"`
		}

		var res string
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			DoRequestObj(context.Background(), &Req{
				Pagination: Pagination{Page: 1, Size: 20},
				Keyword:    "golang",
			}, &res)

		t.AssertNil(err)
		// Verify named struct with in:"query" is flattened to query
		t.Assert(res, "query_page=1,query_size=20,query_keyword=golang,body_empty=true")
	})
}

// Test_DoRequestObj_InTag_RequestInspection demonstrates how to inspect the full request
func Test_DoRequestObj_InTag_RequestInspection(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/inspect/{id}", func(r *ghttp.Request) {
		// Comprehensive request inspection
		type RequestInfo struct {
			// Path parameters
			PathId string `json:"path_id"`

			// Query parameters
			QueryPage string   `json:"query_page"`
			QueryTags []string `json:"query_tags"`
			QueryRaw  string   `json:"query_raw"`

			// Headers
			HeaderToken   string `json:"header_token"`
			HeaderVersion string `json:"header_version"`

			// Cookies
			CookieSession string `json:"cookie_session"`

			// Body
			BodyContent string `json:"body_content"`
			BodyName    string `json:"body_name"`
			BodyAge     int    `json:"body_age"`

			// Request metadata
			Method      string `json:"method"`
			URL         string `json:"url"`
			ContentType string `json:"content_type"`
		}

		info := RequestInfo{
			// Extract path parameter
			PathId: r.Get("id").String(),

			// Extract query parameters
			QueryPage: r.URL.Query().Get("page"),
			QueryTags: r.URL.Query()["tags"],
			QueryRaw:  r.URL.RawQuery,

			// Extract headers
			HeaderToken:   r.Header.Get("Authorization"),
			HeaderVersion: r.Header.Get("X-Version"),

			// Extract cookies
			CookieSession: r.Cookie.Get("session").String(),

			// Extract body
			BodyContent: string(r.GetBody()),

			// Request metadata
			Method:      r.Method,
			URL:         r.URL.String(),
			ContentType: r.Header.Get("Content-Type"),
		}

		// Parse body JSON
		bodyMap := r.GetBodyMap()
		info.BodyName = gconv.String(bodyMap["name"])
		info.BodyAge = gconv.Int(bodyMap["age"])

		r.Response.WriteJson(info)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			g.Meta  `path:"/inspect/{id}" method:"post"`
			Id      int      `in:"path"`
			Page    int      `in:"query" json:"page"`
			Tags    []string `in:"query" json:"tags"`
			Token   string   `in:"header" json:"Authorization"`
			Version string   `in:"header" json:"X-Version"`
			Session string   `in:"cookie" json:"session"`
			Name    string   `json:"name"`
			Age     int      `json:"age"`
		}

		type RequestInfo struct {
			PathId        string   `json:"path_id"`
			QueryPage     string   `json:"query_page"`
			QueryTags     []string `json:"query_tags"`
			QueryRaw      string   `json:"query_raw"`
			HeaderToken   string   `json:"header_token"`
			HeaderVersion string   `json:"header_version"`
			CookieSession string   `json:"cookie_session"`
			BodyContent   string   `json:"body_content"`
			BodyName      string   `json:"body_name"`
			BodyAge       int      `json:"body_age"`
			Method        string   `json:"method"`
			URL           string   `json:"url"`
			ContentType   string   `json:"content_type"`
		}

		var res RequestInfo
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).ContentJson().
			DoRequestObj(context.Background(), &Req{
				Id:      123,
				Page:    1,
				Tags:    []string{"go", "web", "api"},
				Token:   "Bearer secret-token",
				Version: "v2.0",
				Session: "session-abc123",
				Name:    "Alice",
				Age:     30,
			}, &res)

		t.AssertNil(err)

		// Verify path parameter
		t.Assert(res.PathId, "123")

		// Verify query parameters
		t.Assert(res.QueryPage, "1")
		t.Assert(len(res.QueryTags), 3)
		t.Assert(res.QueryTags[0], "go")
		t.Assert(res.QueryTags[1], "web")
		t.Assert(res.QueryTags[2], "api")
		t.Assert(strings.Contains(res.QueryRaw, "page=1"), true)
		t.Assert(strings.Contains(res.QueryRaw, "tags=go"), true)
		t.Assert(strings.Contains(res.QueryRaw, "tags=web"), true)
		t.Assert(strings.Contains(res.QueryRaw, "tags=api"), true)

		// Verify headers
		t.Assert(res.HeaderToken, "Bearer secret-token")
		t.Assert(res.HeaderVersion, "v2.0")

		// Verify cookies
		t.Assert(res.CookieSession, "session-abc123")

		// Verify body
		t.Assert(res.BodyName, "Alice")
		t.Assert(res.BodyAge, 30)
		t.Assert(strings.Contains(res.BodyContent, `"name":"Alice"`), true)
		t.Assert(strings.Contains(res.BodyContent, `"age":30`), true)

		// Verify request metadata
		t.Assert(res.Method, "POST")
		t.Assert(strings.Contains(res.URL, "/inspect/123"), true)
		t.Assert(strings.Contains(res.ContentType, "application/json"), true)
	})
}

// Test_DoRequestObj_InTag_FileUpload tests file upload with in tag
func Test_DoRequestObj_InTag_FileUpload(t *testing.T) {
	// Create test file
	testFile := gfile.Temp(guid.S())
	defer gfile.Remove(testFile)
	gfile.PutContents(testFile, "test file content for upload")

	s := g.Server(guid.S())
	s.BindHandler("/upload/{id}", func(r *ghttp.Request) {
		// Verify path parameter
		pathId := r.Get("id").String()

		// Verify query parameter
		queryCategory := r.URL.Query().Get("category")

		// Verify header
		headerToken := r.Header.Get("Authorization")

		// Verify file upload
		file := r.GetUploadFile("file")
		var fileContent string
		if file != nil {
			content, _ := file.Open()
			defer content.Close()
			data := make([]byte, file.Size)
			content.Read(data)
			fileContent = string(data)
		}

		// Verify other form field
		description := r.Get("description").String()

		r.Response.Writef("path_id=%s,query_category=%s,header_token=%s,file_content=%s,description=%s",
			pathId, queryCategory, headerToken, fileContent, description)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			g.Meta      `path:"/upload/{id}" method:"post"`
			Id          int    `in:"path"`
			Category    string `in:"query" json:"category"`
			Token       string `in:"header" json:"Authorization"`
			File        string `json:"file"`        // File upload with @file: prefix
			Description string `json:"description"` // Regular form field
		}

		var res string
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			DoRequestObj(context.Background(), &Req{
				Id:          123,
				Category:    "documents",
				Token:       "Bearer upload-token",
				File:        "@file:" + testFile,
				Description: "Test document",
			}, &res)

		t.AssertNil(err)
		t.Assert(res, "path_id=123,query_category=documents,header_token=Bearer upload-token,file_content=test file content for upload,description=Test document")
	})
}

// Test_DoRequestObj_InTag_MultipleFileUpload tests multiple file upload
func Test_DoRequestObj_InTag_MultipleFileUpload(t *testing.T) {
	// Create test files
	testFile1 := gfile.Temp(guid.S())
	testFile2 := gfile.Temp(guid.S())
	defer func() {
		gfile.Remove(testFile1)
		gfile.Remove(testFile2)
	}()
	gfile.PutContents(testFile1, "content1")
	gfile.PutContents(testFile2, "content2")

	s := g.Server(guid.S())
	s.BindHandler("/upload", func(r *ghttp.Request) {
		file1 := r.GetUploadFile("file1")
		file2 := r.GetUploadFile("file2")

		var content1, content2 string
		if file1 != nil {
			f, _ := file1.Open()
			defer f.Close()
			data := make([]byte, file1.Size)
			f.Read(data)
			content1 = string(data)
		}
		if file2 != nil {
			f, _ := file2.Open()
			defer f.Close()
			data := make([]byte, file2.Size)
			f.Read(data)
			content2 = string(data)
		}

		r.Response.Writef("file1=%s,file2=%s", content1, content2)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			g.Meta `path:"/upload" method:"post"`
			File1  string `json:"file1"`
			File2  string `json:"file2"`
		}

		var res string
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			DoRequestObj(context.Background(), &Req{
				File1: "@file:" + testFile1,
				File2: "@file:" + testFile2,
			}, &res)

		t.AssertNil(err)
		t.Assert(res, "file1=content1,file2=content2")
	})
}
