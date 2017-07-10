package main

import (
    "g"
    "testing"
)

func BenchmarkJsonDecodeUsingGlobalVar(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        json := `{"name":"中国","age":31,"list":[["a","b","c"],["d","e","f"]],"item":{"title":"make\"he moon","name":"make'he moon","content":"'[}]{[}he moon"}}`
        g.Json.Decode(&json)
    }
}

func BenchmarkJsonDecodeUsingPackage(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        json := `{"name":"中国","age":31,"list":[["a","b","c"],["d","e","f"]],"item":{"title":"make\"he moon","name":"make'he moon","content":"'[}]{[}he moon"}}`
        g.JsonDecode(&json)
    }
}

func BenchmarkJsonDecodeUsingGlobalVarWithGo(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        json := `{"name":"中国","age":31,"list":[["a","b","c"],["d","e","f"]],"item":{"title":"make\"he moon","name":"make'he moon","content":"'[}]{[}he moon"}}`
        go g.Json.Decode(&json)
    }
}

func BenchmarkJsonDecodeUsingPackageWithGo(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        json := `{"name":"中国","age":31,"list":[["a","b","c"],["d","e","f"]],"item":{"title":"make\"he moon","name":"make'he moon","content":"'[}]{[}he moon"}}`
        go g.JsonDecode(&json)
    }
}
