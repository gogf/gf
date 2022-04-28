package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
)

func main() {
	fileContent := gfile.GetContents(`/Users/john/Workspace/Go/GOPATH/src/git.code.oa.com/Khaos/eros/app/khaos-oss/internal/logic/resource/resource_horizontal_downgrade.go`)
	matches, err := gregex.MatchAllString(`func \(\w+ (.+?)\) ([\s\S]+?) {`, fileContent)
	fmt.Println(err)
	g.Dump(matches)
}
