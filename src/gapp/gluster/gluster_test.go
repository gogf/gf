package main

import (
    "testing"
    "g/net/ghttp"
)

func BenchmarkKVSet(b *testing.B) {
    r := ghttp.Post("http://127.0.0.1:4168/kv", "{\"this_is_key\":\"this_is_value\"}")
    r.Close()
}

func BenchmarkKVGet(b *testing.B) {
    r := ghttp.Get("http://127.0.0.1:4168/kv?k=this_is_key")
    r.Close()
}

