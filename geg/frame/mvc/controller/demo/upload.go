package demo

import (
<<<<<<< HEAD
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/net/ghttp"
)

func Upload(r *ghttp.Request) {
    if f, h, e := r.FormFile("upload-file"); e == nil {
        defer f.Close()
        fname  := gfile.Basename(h.Filename)
        buffer := make([]byte, h.Size)
        f.Read(buffer)
        gfile.PutBinContents("/tmp/" + fname, buffer)
        r.Response.Write(fname + " uploaded successly")
    } else {
        r.Response.Write(e.Error())
    }
}

func UploadShow(r *ghttp.Request) {
    r.Response.Write(`
=======
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/gfile"
)

func Upload(r *ghttp.Request) {
	if f, h, e := r.FormFile("upload-file"); e == nil {
		defer f.Close()
		fname := gfile.Basename(h.Filename)
		buffer := make([]byte, h.Size)
		f.Read(buffer)
		gfile.PutBinContents("/tmp/"+fname, buffer)
		r.Response.Write(fname + " uploaded successly")
	} else {
		r.Response.Write(e.Error())
	}
}

func UploadShow(r *ghttp.Request) {
	r.Response.Write(`
>>>>>>> upstream/master
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

func init() {
<<<<<<< HEAD
    ghttp.GetServer().BindHandler("/upload",      Upload)
    ghttp.GetServer().BindHandler("/upload/show", UploadShow)
}
=======
	ghttp.GetServer().BindHandler("/upload", Upload)
	ghttp.GetServer().BindHandler("/upload/show", UploadShow)
}
>>>>>>> upstream/master
