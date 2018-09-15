package main

import (
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/util/gutil"
)

func main() {
    gutil.Dump(gview.ParseContent(`{{"<div>测试</div>去掉HTML标签<script>var v=1;</script>"|text}}`, nil))
}
