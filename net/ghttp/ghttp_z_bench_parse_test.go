// Benchmark for the standard handler signature parse chain:
// func(context.Context, *Req)(*Res, error).
//
// Measurement tooling only, no assertions on the results.
package ghttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/util/guid"
)

type benchParseReq struct {
	Name    string `json:"name"    v:"required|length:1,30#name is required"`
	Age     int    `json:"age"     v:"between:1,100#age should be between 1 and 100"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Tags    string `json:"tags"`
	Remark  string `json:"remark"`
}

type benchParseReqNoTag struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Tags    string `json:"tags"`
	Remark  string `json:"remark"`
}

type benchParseRes struct {
	Name string `json:"name"`
}

var (
	benchParseJSONBody = `{"name":"john","age":18,"email":"john@example.com","address":"street 1","tags":"a,b","remark":"hello"}`
	benchParseQuery    = "/bench?name=john&age=18&email=john@example.com&address=street+1&tags=a,b&remark=hello"

	benchParseHandler = func(ctx context.Context, req *benchParseReq) (*benchParseRes, error) {
		return &benchParseRes{Name: req.Name}, nil
	}
	benchParseHandlerNoTag = func(ctx context.Context, req *benchParseReqNoTag) (*benchParseRes, error) {
		return &benchParseRes{Name: req.Name}, nil
	}
)

// newBenchParseServer creates a server with given handler registered on /bench.
func newBenchParseServer(name string, handler any) (*Server, handlerFuncInfo) {
	s := GetServer(name + "-" + guid.S())
	s.sessionManager = gsession.New(time.Hour, gsession.NewStorageMemory())
	funcInfo, err := s.checkAndCreateFuncInfo(handler, "", "", "")
	if err != nil {
		panic(err)
	}
	s.doBindHandler(context.Background(), doBindHandlerInput{
		Pattern:  "ALL:/bench",
		FuncInfo: funcInfo,
	})
	return s, funcInfo
}

// newBenchParseRequest creates a request and populates the routing information,
// as it is done in the real serving procedure.
func newBenchParseRequest(s *Server, method, path, contentType, body string) *Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	r := newRequest(s, req, httptest.NewRecorder())
	r.handlers, r.serveHandler, r.hasHookHandler, r.hasServeHandler = s.getHandlersWithCache(r)
	return r
}

func newBenchParseFormBody() string {
	return url.Values{
		"name":    {"john"},
		"age":     {"18"},
		"email":   {"john@example.com"},
		"address": {"street 1"},
		"tags":    {"a,b"},
		"remark":  {"hello"},
	}.Encode()
}

// BenchmarkParse_RequestOnly measures the request object creation and routing
// only, which is the common overhead of the following benchmarks.
func BenchmarkParse_RequestOnly(b *testing.B) {
	s, _ := newBenchParseServer("bench-parse-req", benchParseHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := newBenchParseRequest(s, http.MethodPost, "/bench", "application/json", benchParseJSONBody)
		_ = r
	}
}

// BenchmarkParse_JSONPost measures the parse chain for a typical JSON post request.
func BenchmarkParse_JSONPost(b *testing.B) {
	s, _ := newBenchParseServer("bench-parse-json", benchParseHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := newBenchParseRequest(s, http.MethodPost, "/bench", "application/json", benchParseJSONBody)
		var req benchParseReq
		if err := r.Parse(&req); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParse_JSONPost_NoTag measures the parse chain with a request struct
// that has no validation tags at all.
func BenchmarkParse_JSONPost_NoTag(b *testing.B) {
	s, _ := newBenchParseServer("bench-parse-json-notag", benchParseHandlerNoTag)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := newBenchParseRequest(s, http.MethodPost, "/bench", "application/json", benchParseJSONBody)
		var req benchParseReqNoTag
		if err := r.Parse(&req); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParse_QueryGet measures the parse chain for a typical query get request.
func BenchmarkParse_QueryGet(b *testing.B) {
	s, _ := newBenchParseServer("bench-parse-query", benchParseHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := newBenchParseRequest(s, http.MethodGet, benchParseQuery, "", "")
		var req benchParseReq
		if err := r.Parse(&req); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParse_FormPost measures the parse chain for a form urlencoded post request.
func BenchmarkParse_FormPost(b *testing.B) {
	s, _ := newBenchParseServer("bench-parse-form", benchParseHandler)
	body := newBenchParseFormBody()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := newBenchParseRequest(s, http.MethodPost, "/bench", "application/x-www-form-urlencoded", body)
		var req benchParseReq
		if err := r.Parse(&req); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseHandlerFunc_JSONPost measures the whole router function of a
// standard handler: request object creation, parse and the reflection call.
func BenchmarkParseHandlerFunc_JSONPost(b *testing.B) {
	s, funcInfo := newBenchParseServer("bench-parse-func", benchParseHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := newBenchParseRequest(s, http.MethodPost, "/bench", "application/json", benchParseJSONBody)
		funcInfo.Func(r)
		if r.error != nil {
			b.Fatal(r.error)
		}
	}
}

// BenchmarkParseHandlerFunc_JSONPost_NoTag measures the whole router function
// of a standard handler whose request struct has no validation tag.
func BenchmarkParseHandlerFunc_JSONPost_NoTag(b *testing.B) {
	s, funcInfo := newBenchParseServer("bench-parse-func-notag", benchParseHandlerNoTag)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := newBenchParseRequest(s, http.MethodPost, "/bench", "application/json", benchParseJSONBody)
		funcInfo.Func(r)
		if r.error != nil {
			b.Fatal(r.error)
		}
	}
}

// BenchmarkServeHTTP_RequestBuildOnly measures the httptest request/recorder
// building cost itself, the common overhead of BenchmarkServeHTTP_JSONPost.
func BenchmarkServeHTTP_RequestBuildOnly(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/bench", strings.NewReader(benchParseJSONBody))
		req.Header.Set("Content-Type", "application/json")
		_ = httptest.NewRecorder()
	}
}

// BenchmarkServeHTTP_JSONPost measures the full server serving procedure for a
// standard handler JSON post request, as the end-to-end reference number.
func BenchmarkServeHTTP_JSONPost(b *testing.B) {
	s, _ := newBenchParseServer("bench-serve-http", benchParseHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/bench", strings.NewReader(benchParseJSONBody))
		req.Header.Set("Content-Type", "application/json")
		s.ServeHTTP(httptest.NewRecorder(), req)
	}
}
