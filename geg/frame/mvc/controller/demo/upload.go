package demo

import (
    "fmt"
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
        r.Response.WriteString(fmt.Sprintf("%s upload success, input value:%s", fname, r.GetPostString("name")))
    } else {
        r.Response.WriteString(e.Error())
    }
}

func UploadShow(r *ghttp.Request) {
    r.Response.WriteString(`
<html>
<head>
    <title>上传文件</title>
</head>
    <body>
        <form enctype="multipart/form-data" action="/upload" method="post">
            <input type="input" name="name" />
            <input type="file" name="upload-file" />
            <input type="submit" value="upload" />
        </form>
    </body>
</html>
`)
}

func init() {
    ghttp.GetServer().BindHandler("/upload",      Upload)
    ghttp.GetServer().BindHandler("/upload/show", UploadShow)
}