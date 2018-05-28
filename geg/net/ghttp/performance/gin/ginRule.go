package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    r.Use(gin.Recovery())
    r.GET("/:name", func(c *gin.Context) {
        c.String(http.StatusOK, c.Param("name"))
    })
    r.Run(":8199")
}