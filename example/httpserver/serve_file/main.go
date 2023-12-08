package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var pathMap = map[string]string{
	"/aaa/": "./",
}

// ServeFile serves the file to the response.
func ServeFile(r *ghttp.Request) {
	truePath := r.URL.Path
	// Replace the path prefix.
	for k, v := range pathMap {
		if strings.HasPrefix(truePath, k) {
			truePath = strings.Replace(truePath, k, v, 1)
			break
		}
	}

	// Use file from dist.
	file, err := os.Open(truePath)
	if err != nil {
		r.Response.WriteStatus(http.StatusForbidden)
		return
	}
	defer file.Close()

	// Clear the response buffer before file serving.
	// It ignores all custom buffer content and uses the file content.
	r.Response.ClearBuffer()

	info, _ := file.Stat()
	if info.IsDir() {
		r.Response.WriteStatus(http.StatusForbidden)
	} else {
		r.Response.ServeContent(info.Name(), info.ModTime(), file)
	}
}

func main() {
	s := g.Server()
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.BindHandler("/*", ServeFile)
	s.SetPort(8080)
	s.Run()
}

// http://127.0.0.1:8080/aaa/main.go
