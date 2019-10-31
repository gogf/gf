package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"io"
)

// Upload uploads file to /tmp .
func Upload(r *ghttp.Request) {
	f, h, e := r.FormFile("upload-file")
	if e != nil {
		r.Response.Write(e)
	}
	defer f.Close()
	savePath := "/tmp/" + gfile.Basename(h.Filename)
	file, err := gfile.Create(savePath)
	if err != nil {
		r.Response.Write(err)
		return
	}
	defer file.Close()
	if _, err := io.Copy(file, f); err != nil {
		r.Response.Write(err)
		return
	}
	r.Response.Write("upload successfully")
}

// UploadShow shows uploading page.
func UploadShow(r *ghttp.Request) {
	r.Response.Write(`
    <html>
    <head>
        <title>GF UploadFile Demo</title>
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
	s.Group("/upload", func(g *ghttp.RouterGroup) {
		g.ALL("/", Upload)
		g.ALL("/show", UploadShow)
	})
	s.SetPort(8199)
	s.Run()
}
