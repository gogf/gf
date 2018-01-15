package groutine_test

import (
    "testing"
    "gitee.com/johng/gf/g/os/groutine"
)

func test() {
    num := 0
    for i := 0; i < 1000000; i++ {
        num += i
    }
}

var pool = groutine.New()

func BenchmarkGroutine(b *testing.B) {
    for i := 0; i < b.N; i++ {
        pool.Add(test)
    }
    //pool.Close()
}

//func BenchmarkGoRoutine(b *testing.B) {
//    t := gtime.Microsecond()
//    b.N = 100000
//    for i := 0; i < b.N; i++ {
//        go test()
//    }
//    fmt.Println("BenchmarkGoRoutine costs:", gtime.Microsecond() - t)
//}