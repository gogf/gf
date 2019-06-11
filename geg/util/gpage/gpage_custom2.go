package main

import (
<<<<<<< HEAD
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/util/gpage"
=======
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/gview"
	"github.com/gogf/gf/g/util/gpage"
>>>>>>> upstream/master
)

// 自定义分页名称
func pageContent(page *gpage.Page) string {
<<<<<<< HEAD
    page.NextPageTag  = "NextPage"
    page.PrevPageTag  = "PrevPage"
    page.FirstPageTag = "HomePage"
    page.LastPageTag  = "LastPage"
    pageStr := page.FirstPage()
    pageStr += page.PrevPage()
    pageStr += page.PageBar("current-page")
    pageStr += page.NextPage()
    pageStr += page.LastPage()
    return pageStr
}

func main() {
    s := ghttp.GetServer()
    s.BindHandler("/page/custom2/*page", func(r *ghttp.Request){
        page      := gpage.New(100, 10, r.Get("page"), r.URL.String(), r.Router.Uri)
        buffer, _ := gview.ParseContent(`
=======
	page.NextPageTag = "NextPage"
	page.PrevPageTag = "PrevPage"
	page.FirstPageTag = "HomePage"
	page.LastPageTag = "LastPage"
	pageStr := page.FirstPage()
	pageStr += page.PrevPage()
	pageStr += page.PageBar("current-page")
	pageStr += page.NextPage()
	pageStr += page.LastPage()
	return pageStr
}

func main() {
	s := ghttp.GetServer()
	s.BindHandler("/page/custom2/*page", func(r *ghttp.Request) {
		page := gpage.New(100, 10, r.Get("page"), r.URL.String(), r.Router)
		buffer, _ := gview.ParseContent(`
>>>>>>> upstream/master
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
<<<<<<< HEAD
            "page" : gview.HTML(pageContent(page)),
        })
        r.Response.Write(buffer)
    })
    s.SetPort(10000)
    s.Run()
}
=======
			"page": pageContent(page),
		})
		r.Response.Write(buffer)
	})
	s.SetPort(10000)
	s.Run()
}
>>>>>>> upstream/master
