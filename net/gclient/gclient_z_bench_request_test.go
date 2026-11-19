// Benchmark for the multipart request parameter handling path.
package gclient

import "testing"

func BenchmarkNewMultipartRequest(b *testing.B) {
	c := &Client{}
	params := "name=john&age=18&email=john%40example.com&tags=a%2Cb&remark=hello"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := c.newMultipartRequest("POST", "http://example.com/upload", params); err != nil {
			b.Fatal(err)
		}
	}
}
