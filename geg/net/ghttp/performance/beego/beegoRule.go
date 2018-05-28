package main

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
)

func main() {
    beego.Get("/:name",func(ctx *context.Context){
        ctx.Output.Body([]byte(ctx.Input.Param(":name")))
    })
    beego.BeeLogger.SetLevel(0)
    beego.Run(":8199")
}