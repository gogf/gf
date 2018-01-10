package demo

import (
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/net/ghttp"
    "fmt"
)

func Upload(r *ghttp.Request) {
    fmt.Println(r.GetPostMap(nil))
    fmt.Println(r.GetPostString("name"))
    //fmt.Println(string(r.GetRaw()))
    return
    if f, h, e := r.FormFile("upload-file"); e == nil {
        defer f.Close()
        buffer := make([]byte, h.Size)
        f.Read(buffer)
        gfile.PutBinContents("/tmp/" + h.Filename, buffer)
        r.Response.WriteString(fmt.Sprintf("%s upload success, input value:%s", h.Filename, r.GetPostString("name")))
    } else {
        glog.Error(e)
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