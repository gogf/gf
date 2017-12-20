package main
import (
    "fmt"
    "gitee.com/johng/gf/g/frame/ginstance"
)

func main() {
    db := ginstance.Database()
    list, _ := db.Table("test").Select()
    fmt.Println(list[0])
}