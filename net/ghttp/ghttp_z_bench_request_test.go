// Benchmarks for the request string handling paths.
package ghttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const benchRawQuery = "name=john&tags=a&tags=b&tags=c&age=18&email=john%40example.com&remark=hello%20world"

func BenchmarkParseRepeatedQueryValues(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := parseRepeatedQueryValues(benchRawQuery)
		if len(m) == 0 {
			b.Fatal("no repeated values parsed")
		}
	}
}

func BenchmarkGetClientIp(b *testing.B) {
	r := &Request{Request: httptest.NewRequest(http.MethodGet, "/", nil)}
	r.Header.Set(HeaderXForwardedFor, "203.0.113.7, 198.51.100.1, 192.0.2.9")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.clientIp = ""
		if ip := r.GetClientIp(); ip != "203.0.113.7" {
			b.Fatalf("unexpected ip: %s", ip)
		}
	}
}
