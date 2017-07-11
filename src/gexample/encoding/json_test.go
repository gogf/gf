package main

import (
    "testing"
    "g/encoding/gjson"
)

// go test json_test.go -bench=".*"

func BenchmarkJsonDecode(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        json := `{"name":"中国","age":31,"list":[["a","b","c"],["d","e","f"]],"item":{"title":"make\"he moon","name":"make'he moon","content":"'[}]{[}he moon"}}`
        gjson.Decode(&json)
    }
}

func BenchmarkJsonDecodeWithGo(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        json := `{"name":"中国","age":31,"list":[["a","b","c"],["d","e","f"]],"item":{"title":"make\"he moon","name":"make'he moon","content":"'[}]{[}he moon"}}`
        go gjson.Decode(&json)
    }
}
