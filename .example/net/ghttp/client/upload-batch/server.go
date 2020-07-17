package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
)

// Upload uploads files to /tmp .
func Upload(r *ghttp.Request) {
	saveDirPath := "/tmp/"
	files := r.GetUploadFiles("upload-file")
	if _, err := files.Save(saveDirPath); err != nil {
		r.Response.WriteExit(err)
	}
	r.Response.WriteExit("upload successfully")
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
		group.POST("/", Upload)
		group.ALL("/show", UploadShow)
		group.ALL("/batch", UploadShowBatch)
	})
	s.SetPort(8199)
	s.Run()
}
