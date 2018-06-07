package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    r.Use(gin.Recovery())
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "哈喽世界！")
    })
    r.Run(":8199")
}