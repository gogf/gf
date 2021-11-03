package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gview"
)

func main() {
	if c, err := gview.ParseContent(`{{"<div>测试</div>模板引擎默认处理HTML标签<script>alert(\"test\");</script>\n"}}`, nil); err == nil {
		g.Dump(c)
	} else {
		g.Dump(c)
	}
	if c, err := gview.ParseContent(`{{"<div>测试</div>去掉HTML标签<script>alert(\"test\");</script>\n"|text}}`, nil); err == nil {
		g.Dump(c)
	} else {
		g.Dump(c)
	}
	if c, err := gview.ParseContent(`{{"<div>测试</div>保留HTML标签<script>alert(\"test\");</script>\n"|html}}`, nil); err == nil {
		g.Dump(c)
	} else {
		g.Dump(c)
	}
}
