package main

import (
    "g/encoding/gjson"
    "log"
    "g/core/types/gmap"
)

type T struct {
    I int
    J string
}

func main() {
    m := gmap.NewStringStringMap()
    m.Set("name", "john")
    s := gjson.Encode(*m.Clone())
    var t = T {1, *s}
    s2 := gjson.Encode(t)
    log.Println(*s2)

    var t2 = T {}
    err := gjson.DecodeTo(s2, &t2)
    log.Println(err)
    log.Println(t2)

    m2 := make(map[string]string)
    err = gjson.DecodeTo(&t2.J, &m2)
    log.Println(err)
    log.Println(m2)
}