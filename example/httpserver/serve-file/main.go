package main

import (
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// pathMap is used for URL mapping
var pathMap = map[string]string{
	"/aaa/": "/tmp/",
}

// ServeFile serves the file to the response.
func ServeFile(r *ghttp.Request) {
	truePath := r.URL.Path
	hasPrefix := false
	// Replace the path prefix.
	for k, v := range pathMap {
		if strings.HasPrefix(truePath, k) {
			truePath = strings.Replace(truePath, k, v, 1) // Replace only once.
			hasPrefix = true
			break
		}
	}

	if !hasPrefix {
		r.Response.WriteStatus(http.StatusForbidden)
		return
	}

	r.Response.ServeFile(truePath)
}

func main() {
	s := g.Server()
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.BindHandler("/*", ServeFile)
	s.SetPort(8080)
	s.Run()
}

// http://127.0.0.1:8080/aaa/main.go
