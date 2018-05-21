package main

import "fmt"


func test(a []byte) {
    fmt.Println(a)
}

func main() {
    b := []byte{0,1,2,3,4,5,6,7,8,9}
    a := []byte{}
    a = append(a, b[0:2]...)
    a = append(a, b[7:10]...)
    test(a)
    fmt.Println(b)
}
