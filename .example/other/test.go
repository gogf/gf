package main

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
)

func main() {
	arrayIM := garray.NewStrArraySize(0, 2000000)
	arrayNoIM := garray.NewStrArraySize(0, 2000000)
	content := gfile.GetContents("/Users/john/Downloads/keys.txt")
	for _, line := range gstr.Split(content, "\n") {
		if gstr.HasPrefix(line, "up:") ||
			gstr.HasPrefix(line, "helper") ||
			gstr.HasPrefix(line, "u5_0_fl") ||
			gstr.HasPrefix(line, "ubul#") ||
			gstr.HasPrefix(line, "uf:") ||
			gstr.HasPrefix(line, "ufr:") ||
			gstr.HasPrefix(line, "friend:") ||
			gstr.HasPrefix(line, "ubl:") {
			arrayIM.Append(line)
		} else {
			arrayNoIM.Append(line)
		}
	}
	gfile.PutContents("/Users/john/Downloads/keysIM.txt", arrayIM.Join("\n"))
	gfile.PutContents("/Users/john/Downloads/keysNoIM.txt", arrayNoIM.Join("\n"))
}
