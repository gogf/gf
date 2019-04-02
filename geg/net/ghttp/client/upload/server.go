package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/gfile"
)

// 执行文件上传处理，上传到系统临时目录 /tmp
func Upload(r *ghttp.Request) {
	if f, h, e := r.FormFile("upload-file"); e == nil {
		defer f.Close()
		name := gfile.Basename(h.Filename)
		buffer := make([]byte, h.Size)
		f.Read(buffer)
		gfile.PutBinContents("/tmp/"+name, buffer)
		r.Response.Write(name + " uploaded successly")
	} else {
		r.Response.Write(e.Error())
	}
}

// 展示文件上传页面
func UploadShow(r *ghttp.Request) {
	r.Response.Write(`
    <html>
    <head>
        <title>上传文件</title>
    </head>
        <body>
            <form enctype="multipart/form-data" action="/upload" method="post">
                <input type="file" name="upload-file" />
                <input type="submit" value="upload" />
            </form>
        </body>
    </html>
    `)
}

func main() {
	s := g.Server()
	s.BindHandler("/upload", Upload)
	s.BindHandler("/upload/show", UploadShow)
	s.SetPort(8199)
	s.Run()
}
