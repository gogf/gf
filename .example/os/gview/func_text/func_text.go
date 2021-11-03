package main

import (
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/util/gutil"
)

func main() {
	gutil.Dump(gview.ParseContent(`{{"<div>测试</div>去掉HTML标签<script>var v=1;</script>"|text}}`, nil))
}
