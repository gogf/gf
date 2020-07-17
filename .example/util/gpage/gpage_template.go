package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/os/gview"
)

func main() {
	s := g.Server()
	s.BindHandler("/page/template/{page}.html", func(r *ghttp.Request) {
		page := r.GetPage(100, 10)
		page.UrlTemplate = "/order/list/{.page}.html"
		buffer, _ := gview.ParseContent(`
        <html>
            <head>
                <style>
                    a,span {padding:8px; font-size:16px;}
                    div{margin:5px 5px 20px 5px}
                </style>
            </head>
            <body>
                <div>{{.page1}}</div>
                <div>{{.page2}}</div>
                <div>{{.page3}}</div>
                <div>{{.page4}}</div>
            </body>
        </html>
        `, g.Map{
			"page1": page.GetContent(1),
			"page2": page.GetContent(2),
			"page3": page.GetContent(3),
			"page4": page.GetContent(4),
		})
		r.Response.Write(buffer)
	})
	s.SetPort(8199)
	s.Run()
}
