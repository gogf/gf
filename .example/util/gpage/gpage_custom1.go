package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/os/gview"
	"github.com/jin502437344/gf/text/gstr"
	"github.com/jin502437344/gf/util/gpage"
)

// wrapContent wraps each of the page tag with html li and ul.
func wrapContent(page *gpage.Page) string {
	content := page.GetContent(4)
	content = gstr.ReplaceByMap(content, map[string]string{
		"<span":  "<li><span",
		"/span>": "/span></li>",
		"<a":     "<li><a",
		"/a>":    "/a></li>",
	})
	return "<ul>" + content + "</ul>"
}

func main() {
	s := g.Server()
	s.BindHandler("/page/custom1/*page", func(r *ghttp.Request) {
		page := r.GetPage(100, 10)
		content := wrapContent(page)
		buffer, _ := gview.ParseContent(`
        <html>
            <head>
                <style>
                    a,span {padding:8px; font-size:16px;}
                    div{margin:5px 5px 20px 5px}
                </style>
            </head>
            <body>
                <div>{{.page}}</div>
            </body>
        </html>
        `, g.Map{
			"page": content,
		})
		r.Response.Write(buffer)
	})
	s.SetPort(10000)
	s.Run()
}
