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
	"github.com/gogf/gf/v2/util/guid"
)

// Test_Client_Upload_QueryParams_Basic tests basic file upload with query parameters
func Test_Client_Upload_QueryParams_Basic(t *testing.T) {
	s := createUploadServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Create temporary file with random name for upload
		tempDir := gfile.Temp(guid.S())
		gfile.Mkdir(tempDir)
		uploadFilePath := gfile.Join(tempDir, "random_upload_"+guid.S()+".txt")
		gfile.PutContents(uploadFilePath, "test content for upload "+guid.S())
		defer gfile.Remove(tempDir)

		// Test 1: Basic file upload with query parameters
		data := g.Map{
			"file": "@file:" + uploadFilePath,
		}

		// Add query parameters using Query method
		resp := c.Query(g.Map{
			"category": "documents",
			"type":     "text",
		}).PostContent(context.Background(), "/upload", data)

		// Verify that both query parameters and file content are present
		t.Assert(strings.Contains(resp, "category=[documents]"), true)
		t.Assert(strings.Contains(resp, "type=[text]"), true)
		t.Assert(strings.Contains(resp, "test content for upload"), true) // actual file content
	})
}

// Test_Client_Upload_QueryParams_MixedScenarios tests various combinations of file upload with query parameters
func Test_Client_Upload_QueryParams_MixedScenarios(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/upload", func(r *ghttp.Request) {
		tmpPath := gfile.Temp(guid.S())
		err := gfile.Mkdir(tmpPath)
		gtest.AssertNil(err)
		defer gfile.Remove(tmpPath)

		file := r.GetUploadFile("file")
		_, err = file.Save(tmpPath)
		gtest.AssertNil(err)

		// Get query parameters from URL
		queryParams := r.URL.Query()
		var queryResult string
		for k, v := range queryParams {
			queryResult += fmt.Sprintf("%s=%v,", k, v)
		}

		r.Response.Write(
			"query_params=" + queryResult +
				"file_content=" + gfile.GetContents(gfile.Join(tmpPath, gfile.Basename(file.Filename))) +
				"title=" + r.Get("title").String(),
		)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Create temporary file with random name for upload
		tempDir := gfile.Temp(guid.S())
		gfile.Mkdir(tempDir)
		uploadFilePath := gfile.Join(tempDir, "random_upload_"+guid.S()+".txt")
		gfile.PutContents(uploadFilePath, "upload test content "+guid.S())
		defer gfile.Remove(tempDir)

		// Test 1: File upload with URL query parameters and client query parameters
		resp := c.Query(g.Map{
			"page": 1,
			"size": 10,
		}).PostContent(context.Background(), "/upload?filter=all&sort=date", g.Map{
			"file":  "@file:" + uploadFilePath,
			"title": "test file",
		})

		// Check that URL parameters, query parameters, and form fields are all present
		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
		t.Assert(strings.Contains(resp, "filter=[all]"), true)
		t.Assert(strings.Contains(resp, "sort=[date]"), true)
		t.Assert(strings.Contains(resp, "upload test content"), true) // actual file content
		t.Assert(strings.Contains(resp, "title=test file"), true)     // form field should be present
	})
}

// Test_Client_Upload_QueryParams_Conflicts tests conflict resolution between different parameter sources
func Test_Client_Upload_QueryParams_Conflicts(t *testing.T) {
	s := createUploadServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Create temporary file with random name for upload
		tempDir := gfile.Temp(guid.S())
		gfile.Mkdir(tempDir)
		uploadFilePath := gfile.Join(tempDir, "random_upload_"+guid.S()+".txt")
		gfile.PutContents(uploadFilePath, "conflict test "+guid.S())
		defer gfile.Remove(tempDir)

		// Test conflict resolution: queryParams should override URL params
		resp := c.Query(g.Map{
			"conflict": "from_query", // Higher priority
		}).PostContent(context.Background(), "/upload?conflict=from_url", g.Map{
			"file": "@file:" + uploadFilePath,
		})

		// Query parameters should override URL parameters
		t.Assert(strings.Contains(resp, "conflict=[from_query]"), true)
		t.Assert(strings.Contains(resp, "conflict=[from_url]"), false)
		t.Assert(strings.Contains(resp, "conflict"), true) // Should appear once with correct value
	})
}

// Test_Client_Upload_QueryParams_SliceArrays tests file upload with slice/array query parameters
func Test_Client_Upload_QueryParams_SliceArrays(t *testing.T) {
	s := createUploadServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Create temporary file with random name for upload
		tempDir := gfile.Temp(guid.S())
		gfile.Mkdir(tempDir)
		uploadFilePath := gfile.Join(tempDir, "random_upload_"+guid.S()+".txt")
		gfile.PutContents(uploadFilePath, "slice test content "+guid.S())
		defer gfile.Remove(tempDir)

		// Test slice/array query parameters with file upload
		resp := c.Query(g.Map{
			"tags": []string{"tag1", "tag2", "tag3"},
			"ids":  []int{1, 2, 3},
		}).PostContent(context.Background(), "/upload", g.Map{
			"file": "@file:" + uploadFilePath,
		})

		// Check that slice parameters are properly expanded
		t.Assert(strings.Contains(resp, "tag1"), true)
		t.Assert(strings.Contains(resp, "tag2"), true)
		t.Assert(strings.Contains(resp, "tag3"), true)
		t.Assert(strings.Contains(resp, "1"), true)
		t.Assert(strings.Contains(resp, "2"), true)
		t.Assert(strings.Contains(resp, "3"), true)
		t.Assert(strings.Contains(resp, "slice test content"), true) // actual file content
	})
}

// Test_Client_Upload_QueryParams_Chaining tests file upload with chained query parameter methods
func Test_Client_Upload_QueryParams_Chaining(t *testing.T) {
	s := createUploadServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Create temporary file with random name for upload
		tempDir := gfile.Temp(guid.S())
		gfile.Mkdir(tempDir)
		uploadFilePath := gfile.Join(tempDir, "random_upload_"+guid.S()+".txt")
		gfile.PutContents(uploadFilePath, "chaining test content "+guid.S())
		defer gfile.Remove(tempDir)

		// Test chained query parameter methods with file upload
		chainedClient := c.QueryPair("status", "active").
			QueryPair("priority", "high").
			SetQuery("category", "important")

		resp := chainedClient.PostContent(context.Background(), "/upload", g.Map{
			"file": "@file:" + uploadFilePath,
		})

		// Check that all chained query parameters are present
		t.Assert(strings.Contains(resp, "status=[active]"), true)
		t.Assert(strings.Contains(resp, "priority=[high]"), true)
		t.Assert(strings.Contains(resp, "category=[important]"), true)
		t.Assert(strings.Contains(resp, "chaining test content"), true) // actual file content
	})
}

// Test_Client_Upload_QueryParams_SpecialCharacters tests file upload with special characters in query parameters
func Test_Client_Upload_QueryParams_SpecialCharacters(t *testing.T) {
	s := createUploadServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Create temporary file with random name for upload
		tempDir := gfile.Temp(guid.S())
		gfile.Mkdir(tempDir)
		uploadFilePath := gfile.Join(tempDir, "random_upload_"+guid.S()+".txt")
		gfile.PutContents(uploadFilePath, "special chars test "+guid.S())
		defer gfile.Remove(tempDir)

		// Test special characters in query parameters with file upload
		resp := c.Query(g.Map{
			"query":   "hello world",
			"email":   "test@example.com",
			"path":    "/data/file.txt",
			"chinese": "中文测试",
			"symbols": "!@#$%^&*()",
		}).PostContent(context.Background(), "/upload", g.Map{
			"file": "@file:" + uploadFilePath,
		})

		// Check that special characters are properly handled
		t.Assert(strings.Contains(resp, "hello world"), true)
		t.Assert(strings.Contains(resp, "test@example.com"), true)
		t.Assert(strings.Contains(resp, "/data/file.txt"), true)
		t.Assert(strings.Contains(resp, "中文测试"), true)
		t.Assert(strings.Contains(resp, "!@#$"), true)               // At least some symbols
		t.Assert(strings.Contains(resp, "special chars test"), true) // actual file content
	})
}

// Test_Client_Upload_QueryParams_GetMergedURL tests GetMergedURL functionality with file upload scenarios
func Test_Client_Upload_QueryParams_GetMergedURL(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/upload", func(r *ghttp.Request) {
		r.Response.Write("ok")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test GetMergedURL with query parameters but no actual upload (since upload changes the request method to multipart)
		// We'll test the URL building functionality separately

		// For GET requests with query parameters
		mergedURL, err := c.Query(g.Map{
			"page": 1,
			"size": 10,
		}).GetMergedURL(context.Background(), "GET", "/api/data", g.Map{
			"filter": "active",
		})

		t.AssertNil(err)
		t.Assert(strings.Contains(mergedURL, "page=1"), true)
		t.Assert(strings.Contains(mergedURL, "size=10"), true)
		t.Assert(strings.Contains(mergedURL, "filter=active"), true)
	})
}

// Test_Client_Upload_QueryParams_NestedStruct tests file upload with nested struct query parameters
func Test_Client_Upload_QueryParams_NestedStruct(t *testing.T) {
	s := createUploadServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Create temporary file with random name for upload
		tempDir := gfile.Temp(guid.S())
		gfile.Mkdir(tempDir)
		uploadFilePath := gfile.Join(tempDir, "random_upload_"+guid.S()+".txt")
		gfile.PutContents(uploadFilePath, "nested test content "+guid.S())
		defer gfile.Remove(tempDir)

		// Define a struct for query parameters
		type FilterParams struct {
			Page     int      `json:"page"`
			Size     int      `json:"size"`
			Category string   `json:"category"`
			Tags     []string `json:"tags"`
		}

		params := FilterParams{
			Page:     1,
			Size:     20,
			Category: "documents",
			Tags:     []string{"important", "review"},
		}

		// Test struct query parameters with file upload
		resp := c.QueryParams(params).PostContent(context.Background(), "/upload", g.Map{
			"file": "@file:" + uploadFilePath,
		})

		// Check that struct fields are properly converted to query parameters
		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[20]"), true)
		t.Assert(strings.Contains(resp, "category=[documents]"), true)
		t.Assert(strings.Contains(resp, "important"), true)           // Tag value
		t.Assert(strings.Contains(resp, "review"), true)              // Tag value
		t.Assert(strings.Contains(resp, "nested test content"), true) // actual file content
	})
}

// Test_Client_Upload_QueryParams_EmptyValues tests file upload with empty/nil query parameter values
func Test_Client_Upload_QueryParams_EmptyValues(t *testing.T) {
	s := createUploadServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Create temporary file with random name for upload
		tempDir := gfile.Temp(guid.S())
		gfile.Mkdir(tempDir)
		uploadFilePath := gfile.Join(tempDir, "random_upload_"+guid.S()+".txt")
		gfile.PutContents(uploadFilePath, "empty values test "+guid.S())
		defer gfile.Remove(tempDir)

		// Test with empty and nil values in query parameters
		resp := c.Query(g.Map{
			"empty_string": "",
			"zero_int":     0,
			"nil_value":    nil, // This should be skipped
			"normal":       "value",
		}).PostContent(context.Background(), "/upload", g.Map{
			"file": "@file:" + uploadFilePath,
		})

		// Empty values should be included, nil values should be skipped
		t.Assert(strings.Contains(resp, "empty_string=[]"), true)   // Empty string becomes []
		t.Assert(strings.Contains(resp, "zero_int=[0]"), true)      // Zero value should be included
		t.Assert(strings.Contains(resp, "nil_value="), false)       // Nil value should be skipped
		t.Assert(strings.Contains(resp, "normal=[value]"), true)    // Normal value should be present
		t.Assert(strings.Contains(resp, "empty values test"), true) // actual file content
	})
}

// Test_Client_Upload_QueryParams_NoUrlEncode tests file upload with no URL encoding enabled
func Test_Client_Upload_QueryParams_NoUrlEncode(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/upload", func(r *ghttp.Request) {
		tmpPath := gfile.Temp(guid.S())
		err := gfile.Mkdir(tmpPath)
		gtest.AssertNil(err)
		defer gfile.Remove(tmpPath)

		file := r.GetUploadFile("file")
		_, err = file.Save(tmpPath)
		gtest.AssertNil(err)

		r.Response.Write(
			"raw_query=" + r.URL.RawQuery +
				"file_content=" + gfile.GetContents(gfile.Join(tmpPath, gfile.Basename(file.Filename))),
		)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		// Create temporary file with random name for upload
		tempDir := gfile.Temp(guid.S())
		gfile.Mkdir(tempDir)
		uploadFilePath := gfile.Join(tempDir, "random_upload_"+guid.S()+".txt")
		gfile.PutContents(uploadFilePath, "no encode test "+guid.S())
		defer gfile.Remove(tempDir)

		// Test with no URL encoding
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		c.SetNoUrlEncode(true)

		resp := c.Query(g.Map{
			"path":  "/data/binlog",
			"query": "user=admin&role=super",
		}).PostContent(context.Background(), "/upload", g.Map{
			"file": "@file:" + uploadFilePath,
		})

		// Check that special characters are not encoded
		t.Assert(strings.Contains(resp, "path=/data/binlog"), true)           // Not %2Fdata%2Fbinlog
		t.Assert(strings.Contains(resp, "query=user=admin&role=super"), true) // Not encoded
		t.Assert(strings.Contains(resp, "no encode test"), true)              // actual file content
	})
}

// createUploadServer creates a server for file upload testing
func createUploadServer() *ghttp.Server {
	s := g.Server(guid.S())
	s.BindHandler("/upload", func(r *ghttp.Request) {
		tmpPath := gfile.Temp(guid.S())
		err := gfile.Mkdir(tmpPath)
		gtest.AssertNil(err)
		defer gfile.Remove(tmpPath)

		file := r.GetUploadFile("file")
		_, err = file.Save(tmpPath)
		gtest.AssertNil(err)

		// Get query parameters from URL
		queryParams := r.URL.Query()
		var queryResult string
		for k, v := range queryParams {
			queryResult += fmt.Sprintf("%s=%v,", k, v)
		}

		r.Response.Write(
			"query_params=" + queryResult +
				"file_content=" + gfile.GetContents(gfile.Join(tmpPath, gfile.Basename(file.Filename))),
		)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	return s
}
