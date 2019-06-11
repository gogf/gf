package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gfile"
)

func main() {
	path := "/tmp/temp"
	content := `123
456
789
`
	gfile.PutContents(path, content)
	fmt.Println(gfile.Size(path))
	fmt.Println(gfile.GetBinContentsTilCharByPath(path, '\n', 0))
	fmt.Println(gfile.GetBinContentsTilCharByPath(path, '\n', 3))
	fmt.Println(gfile.GetBinContentsTilCharByPath(path, '\n', 8))
	fmt.Println(gfile.GetBinContentsTilCharByPath(path, '\n', 12))
}
