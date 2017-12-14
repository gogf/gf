package main
import (
    "fmt"
    "regexp"
    "strings"
    "strconv"
    "gitee.com/johng/gf/g/database/gdb"
)

func Add(path string, name ... string) {
    fmt.Println(name)
}

func main() {
    nodestr := "账号@地址:端口 ,密码 , 数据库名称, 数据库类型, 集群角色 , 字符编码, 负载均衡优先级 , 12345"
    reg, _  := regexp.Compile(`(.+)@(.+):([^,]+),([^,]+),([^,]+),([^,]+)`)
    match   := reg.FindStringSubmatch(nodestr)
    if match != nil {
        node := gdb.ConfigNode{
            User : strings.TrimSpace(match[1]),
            Host : strings.TrimSpace(match[2]),
            Port : strings.TrimSpace(match[3]),
            Pass : strings.TrimSpace(match[4]),
            Name : strings.TrimSpace(match[5]),
            Type : strings.TrimSpace(match[6]),
        }
        extra := strings.Split(nodestr[len(match[0]) + 1:], ",")
        if len(extra) > 0 {
            node.Role = strings.TrimSpace(extra[0])
        }
        if len(extra) > 1 {
            node.Charset = strings.TrimSpace(extra[1])
        }
        if len(extra) > 2 {
            node.Priority, _ = strconv.Atoi(strings.TrimSpace(extra[2]))
        }
        if len(extra) > 3 {
            index        := len(extra[0]) + len(extra[1]) + len(extra[2]) + 3
            node.Linkinfo = strings.TrimSpace(nodestr[len(match[0]) + 1 + index:])
        }
        fmt.Println(node)
    }

}