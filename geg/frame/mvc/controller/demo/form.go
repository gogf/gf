package demo

import (
<<<<<<< HEAD
    "gitee.com/johng/gf/g/net/ghttp"
    "fmt"
)

func Form(r *ghttp.Request) {
    fmt.Println(r.GetPostMap())
    fmt.Println(r.GetPostString("name"))
    fmt.Println(r.GetPostString("age"))
=======
	"fmt"
	"github.com/gogf/gf/g/net/ghttp"
)

func Form(r *ghttp.Request) {
	fmt.Println(r.GetPostMap())
	fmt.Println(r.GetPostString("name"))
	fmt.Println(r.GetPostString("age"))
>>>>>>> upstream/master

}

func FormShow(r *ghttp.Request) {
<<<<<<< HEAD
    r.Response.Write(`
=======
	r.Response.Write(`
>>>>>>> upstream/master
<html>
<head>
    <title>表单提交</title>
</head>
    <body>
        <form enctype="application/x-www-form-urlencoded" action="/form" method="post">
            <input type="input" name="name" />
            <input type="input" name="age" />
            <input type="submit" value="submit" />
        </form>
    </body>
</html>
`)
}

func init() {
<<<<<<< HEAD
    ghttp.GetServer().BindHandler("/form",      Form)
    ghttp.GetServer().BindHandler("/form/show", FormShow)
}
=======
	ghttp.GetServer().BindHandler("/form", Form)
	ghttp.GetServer().BindHandler("/form/show", FormShow)
}
>>>>>>> upstream/master
