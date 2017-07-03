package gf

import (
	"os"
    "regexp"
    "errors"
)

// 命令行参数列表
type gtConsoleValue  struct {
    values []string
}

// 命令行选项列表
type gtConsoleOption struct {
    options map[string]string
}

// 终端操作结构体封装
type gtConsole struct {
    Value      gtConsoleValue         // console终端参数-命令参数列表
    Option     gtConsoleOption        // console终端参数-选项参数列表
    cmdFuncMap map[string]gtEmptyFunc // 终端命令及函数地址对应表
}

// 终端管理对象(全局)
var Console gtConsole

// 检查并初始化console参数，在包加载的时候触发
func init() {
    Console.cmdFuncMap     = make(map[string]gtEmptyFunc)
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
    if index < uint8(len(c.values)) {
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

// 绑定命令行参数及对应的命令函数，注意参数是函数的内存地址
// 如果操作失败返回错误信息
func (c gtConsole) BindHandle (cmd string, f gtEmptyFunc) error {
    _, ok := c.cmdFuncMap[cmd]
    if ok {
        return errors.New("duplicated handle for command:" + cmd)
    } else {
        c.cmdFuncMap[cmd] = f
        return nil
    }
}

// 执行命令对应的函数
func (c gtConsole) RunHandle (cmd string) error {
    handle, ok := c.cmdFuncMap[cmd]
    if ok {
        handle()
        return nil
    } else {
        return errors.New("no handle found for command:" + cmd)
    }
}
