package gconsole

import (
	"os"
    "regexp"
    "errors"
)

// 命令行参数列表
type gConsoleValue  struct {
    values []string
}

// 命令行选项列表
type gConsoleOption struct {
    options map[string]string
}

// 终端管理对象(全局)
var Value      gConsoleValue  // console终端参数-命令参数列表
var Option     gConsoleOption // console终端参数-选项参数列表
var cmdFuncMap = make(map[string]func()) // 终端命令及函数地址对应表

// 检查并初始化console参数，在包加载的时候触发
// 初始化时执行，不影响运行时性能
func init() {
    Option.options = make(map[string]string)
    reg       := regexp.MustCompile(`\-\-{0,1}(\w+?)=(.+)`)
    for i := 0; i < len(os.Args); i++ {
        result := reg.FindStringSubmatch(os.Args[i])
        if len(result) > 1 {
            Option.options[result[1]] = result[2]
        } else {
            Value.values = append(Value.values, os.Args[i])
        }
    }
}

// 返回所有的命令行参数values
func (c gConsoleValue) GetAll() []string {
    return c.values
}

// 返回所有的命令行参数options
func (c gConsoleOption) GetAll() map[string]string {
    return c.options
}

// 获得一条指定索引位置的value参数
func (c gConsoleValue) GetIndex(index uint8) (string, bool) {
    if index < uint8(len(c.values)) {
        return c.values[index], true
    }
    return "", false
}

// 获得一条指定索引位置的option参数
func (c gConsoleOption) GetIndex(key string) (string, bool) {
    option, ok := c.options[key]
    if ok {
        return option, true
    }
    return "", false
}

// 绑定命令行参数及对应的命令函数，注意参数是函数的内存地址
// 如果操作失败返回错误信息
func BindHandle (cmd string, f func()) error {
    _, ok := cmdFuncMap[cmd]
    if ok {
        return errors.New("duplicated handle for command:" + cmd)
    } else {
        cmdFuncMap[cmd] = f
        return nil
    }
}

// 执行命令对应的函数
func RunHandle (cmd string) error {
    handle, ok := cmdFuncMap[cmd]
    if ok {
        handle()
        return nil
    } else {
        return errors.New("no handle found for command:" + cmd)
    }
}

// 自动识别命令参数并执行命令参数对应的函数
func AutoRun () error {
    cmd, ok := Value.GetIndex(1);
    if ok {
        if handle, ok := cmdFuncMap[cmd]; ok {
            handle()
            return nil
        } else {
            return errors.New("no handle found for command:" + cmd)
        }
    } else {
        return errors.New("no command found")
    }

}
