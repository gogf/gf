package main

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type UploadReq struct {
	g.Meta `path:"/upload" method:"POST" tags:"Upload" mime:"multipart/form-data" summary:"上传文件"`
	File   *ghttp.UploadFile `p:"file" type:"file" dc:"选择上传文件"`
	Msg    string            `dc:"消息"`
}
type UploadRes struct {
	FileName string `json:"fileName"`
}

type cUpload struct{}

func (u cUpload) Upload(ctx context.Context, req *UploadReq) (*UploadRes, error) {
	if req.File != nil {
		return &UploadRes{
			FileName: req.File.Filename,
		}, nil
	}
	return nil, nil
}

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(cUpload{})
	})
	s.SetClientMaxBodySize(600 * 1024 * 1024) // 600M
	s.SetPort(8199)
	s.SetAccessLogEnabled(true)
	s.Run()
}

// curl --location 'http://127.0.0.1:8199/upload' \
// --form 'file=@"/D:/下载/goframe-v2.5.pdf"' \
// --form 'msg="666"'
