package main

import (
    "fmt"
    "reflect"
    "gitee.com/johng/gf/g/database/gdb"
)

func main() {
    var value interface{}
    value = gdb.Map{"a":1}

    refValue := reflect.ValueOf(value)

    if refValue.Kind() == reflect.Map {
            keys := refValue.MapKeys()
            for _, k := range keys {
                fmt.Println(k, refValue.MapIndex(k).Interface())
            }
    }

}
