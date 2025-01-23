package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// HelloReq hello request
type HelloReq struct {
	g.Meta `path:"/hello" method:"get" sort:"1"`
	Name   string `v:"required" dc:"Your name"`
}

// HelloRes hello response
type HelloRes struct {
	Reply string `dc:"Reply content"`
}

// Hello Controller
type Hello struct{}

// Say function
func (Hello) Say(ctx context.Context, req *HelloReq) (res *HelloRes, err error) {
	g.Log().Debugf(ctx, `receive say: %+v`, req)
	res = &HelloRes{
		Reply: fmt.Sprintf(`Hi %s`, req.Name),
	}
	return
}

// upload file request
type UploadReq struct {
	g.Meta `path:"/upload" method:"POST" tags:"Upload" mime:"multipart/form-data" summary:"上传文件"`
	Files  []*ghttp.UploadFile `json:"files" type:"file" dc:"选择上传多文件"`
	File   *ghttp.UploadFile   `p:"file" type:"file" dc:"选择上传文件"`
	Msg    string              `dc:"消息"`
}

// upload file response
type UploadRes struct {
	FilesName []string `json:"files_name"`
	FileName  string   `json:"file_name"`
	Msg       string   `json:"msg"`
}

// upload file
func (Hello) Upload(ctx context.Context, req *UploadReq) (res *UploadRes, err error) {
	g.Log().Debugf(ctx, `receive say: %+v`, req)
	res = &UploadRes{
		Msg: req.Msg,
	}
	if req.File != nil {
		res.FileName = req.File.Filename
	}
	if len(req.Files) > 0 {
		var filesName []string
		for _, file := range req.Files {
			filesName = append(filesName, file.Filename)
		}
		res.FilesName = filesName
	}
	return
}

const (
	// MySwaggerUITemplate is the custom Swagger UI template.
	MySwaggerUITemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<meta name="description" content="SwaggerUI"/>
	<title>SwaggerUI</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui.min.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui-bundle.js" crossorigin></script>
<script>
	window.onload = () => {
		window.ui = SwaggerUIBundle({
			url:    '{SwaggerUIDocUrl}',
			dom_id: '#swagger-ui',
		});
	};
</script>
</body>
</html>
`
	// OpenapiUITemplate is the OpenAPI UI template.
	OpenapiUITemplate = `
	<!doctype html>
	<html lang="en">
	  <head>
		<meta charset="UTF-8" />
		<title>test</title>
	  </head>
	  <body>
		<div id="openapi-ui-container" spec-url="{SwaggerUIDocUrl}" theme="light"></div>
		<script src="https://cdn.jsdelivr.net/npm/openapi-ui-dist@latest/lib/openapi-ui.umd.js"></script>
	  </body>
	</html>
`
)

func main() {
	s := g.Server()
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			new(Hello),
		)
	})
	s.SetSwaggerUITemplate(MySwaggerUITemplate)
	// s.SetSwaggerUITemplate(OpenapiUITemplate) // files support
	s.Run()
}
