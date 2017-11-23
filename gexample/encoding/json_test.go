package main

import (
    "testing"
    "gf/g/encoding/gjson"
    "encoding/json"
    "log"
)

// go test json_test.go -bench=".*"

var data = `[{"CityId":1, "CityName":"北京", "ProvinceId":1, "CityOrder":1}, {"CityId":5, "CityName":"成都", "ProvinceId":27, "CityOrder":1}]`

func BenchmarkJsonDecode(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        gjson.Decode(data)
    }
}

func BenchmarkJsonDecodeByUnmarshal(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        var citys interface{}
        if err := json.Unmarshal([]byte(data), &citys); err != nil {
            log.Fatalf("JSON unmarshaling failed: %s", err)
        }
        //fmt.Println(citys)
    }
}

