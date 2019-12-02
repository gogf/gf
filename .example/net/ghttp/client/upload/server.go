package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"io"
)

// Upload uploads files to /tmp .
func Upload(r *ghttp.Request) {
	saveDir := "/tmp/"
	for _, item := range r.GetMultipartFiles("upload-file") {
		file, err := item.Open()
		if err != nil {
			r.Response.Write(err)
			return
		}
		defer file.Close()

		f, err := gfile.Create(saveDir + gfile.Basename(item.Filename))
		if err != nil {
			r.Response.Write(err)
			return
		}
		defer f.Close()

		if _, err := io.Copy(f, file); err != nil {
			r.Response.Write(err)
			return
		}
	}
	r.Response.Write("upload successfully")
}

// UploadShow shows uploading simgle file page.
func UploadShow(r *ghttp.Request) {
	r.Response.Write(`
    <html>
    <head>
        <title>GF Upload File Demo</title>
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

// UploadShowBatch shows uploading multiple files page.
func UploadShowBatch(r *ghttp.Request) {
	r.Response.Write(`
    <html>
    <head>
        <title>GF Upload Files Demo</title>
    </head>
        <body>
            <form enctype="multipart/form-data" action="/upload" method="post">
                <input type="file" name="upload-file" />
                <input type="file" name="upload-file" />
                <input type="submit" value="upload" />
            </form>
        </body>
    </html>
    `)
}

func main() {
	s := g.Server()
	s.Group("/upload", func(group *ghttp.RouterGroup) {
		group.ALL("/", Upload)
		group.ALL("/show", UploadShow)
		group.ALL("/batch", UploadShowBatch)
	})
	s.SetPort(8199)
	s.Run()
}
