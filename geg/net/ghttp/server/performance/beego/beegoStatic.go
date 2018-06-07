package main

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func main() {
    beego.Get("/",func(ctx *context.Context){
        ctx.Output.Body([]byte("哈喽世界！"))
    })
    beego.BeeLogger.SetLevel(0)
    beego.Run(":8199")
}