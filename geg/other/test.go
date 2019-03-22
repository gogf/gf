package main

import (
    "fmt"
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/util/gvalid"
)

type User struct {
    Id   int
    name string
}

func NewMessage(mobile string, msgType int, content string, templateId string, param string) *Message {
    return &Message{
        Mobile:     mobile,
        Type:       msgType,
        Content:    content,
        TemplateId: templateId,
        Param:      param,
        config:     make(map[string]interface{}),
    }
}

func main() {
    message := NewMessage("123",1,"456","333","1,3,9")
    fmt.Println(gconv.Map(message))
    if e := gvalid.CheckStruct(message,nil); e != nil {
        g.Dump(e.Maps())
    }

}