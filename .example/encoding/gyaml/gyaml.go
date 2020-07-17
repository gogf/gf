package main

import (
	"fmt"

	"github.com/jin502437344/gf/encoding/gyaml"
)

func main() {
	var yamlStr string = `
#即表示url属性值；
url: http://www.wolfcode.cn 
#即表示server.host属性的值；
server:
    host: http://www.wolfcode.cn 
#数组，即表示server为[a,b,c]
server:
    - 120.168.117.21
    - 120.168.117.22
    - 120.168.117.23
#常量
pi: 3.14   #定义一个数值3.14
hasChild: true  #定义一个boolean值
name: '你好YAML'   #定义一个字符串
`

	i, err := gyaml.Decode([]byte(yamlStr))
	fmt.Println(err)
	fmt.Println(i)
}
