package main
import (
    "sort"
    "fmt"
)

func main() {
    fnums := []uint64{0,12,2,4,5,5,2,uint64(10)}
    sort.Slice(fnums, func(i, j int) bool { return fnums[i] < fnums[j] })
    fmt.Println(fnums)
}