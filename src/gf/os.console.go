package gf

import (
	"os"
    "regexp"
)

type gtConsoleValue  struct {
    values []string
}
type gtConsoleOption struct {
    options map[string]string
}

// 终端操作结构体封装
type console struct {
    inited bool            // console终端参数是否已经初始化
    Value  gtConsoleValue  // console终端参数-命令参数列表
    Option gtConsoleOption // console终端参数-选项参数列表
}

// 终端管理对象(全局)
var Console console

// 检查并初始化console参数，在包加载的时候触发
func init() {
    if !Console.inited {
        Console.Option.options = make(map[string]string)
        reg       := regexp.MustCompile(`\-\-{0,1}(\w+?)=(.+)`)
        for i := 1; i < len(os.Args); i++ {
            result := reg.FindStringSubmatch(os.Args[i])
            if len(result) > 1 {
                Console.Option.options[result[1]] = result[2]
            } else {
                Console.Value.values = append(Console.Value.values, os.Args[i])
            }
        }
        Console.inited = true
    }
}

// 返回所有的命令行参数values
func (c gtConsoleValue) GetAll() []string {
    return c.values
}

// 返回所有的命令行参数options
func (c gtConsoleOption) GetAll() map[string]string {
    return c.options
}

// 获得一条指定索引位置的value参数
func (c gtConsoleValue) GetIndex(index uint8) (string, bool) {
    if index < uint8(len(Console.Value.values)) {
        return c.values[index], true
    }
    return "", false
}

// 获得一条指定索引位置的option参数
func (c gtConsoleOption) GetIndex(key string) (string, bool) {
    option, ok := c.options[key]
    if ok {
        return option, true
    }
    return "", false
}
