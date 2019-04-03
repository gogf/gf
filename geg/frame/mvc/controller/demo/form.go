package demo

import (
	"fmt"
	"github.com/gogf/gf/g/net/ghttp"
)

func Form(r *ghttp.Request) {
	fmt.Println(r.GetPostMap())
	fmt.Println(r.GetPostString("name"))
	fmt.Println(r.GetPostString("age"))

}

func FormShow(r *ghttp.Request) {
	r.Response.Write(`
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
	ghttp.GetServer().BindHandler("/form", Form)
	ghttp.GetServer().BindHandler("/form/show", FormShow)
}
