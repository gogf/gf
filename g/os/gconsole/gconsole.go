package gconsole

import (
	"os"
    "regexp"
    "errors"
    "strconv"
    "strings"
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
func (c gConsoleValue) Get(index uint8) string {
    if index < uint8(len(c.values)) {
        return c.values[index]
    }
    return ""
}

// 类型转换
func (c gConsoleValue) GetInt(key uint8) int {
    v := c.Get(key)
    if v != "" {
        i, _ := strconv.Atoi(v)
        return i
    }
    return 0
}

// 类型转换bool
func (c gConsoleValue) GetBool(key uint8) bool {
    v := c.Get(key)
    v  = strings.ToLower(v)
    if v != "" && v != "0" && v != "false" {
        return true
    }
    return false
}

// 获得一条指定索引位置的option参数
func (c gConsoleOption) Get(key string) string {
    option, ok := c.options[key]
    if ok {
        return option
    }
    return ""
}

// 类型转换int
func (c gConsoleOption) GetInt(key string) int {
    v := c.Get(key)
    if v != "" {
        i, _ := strconv.Atoi(v)
        return i
    }
    return 0
}

// 类型转换bool
func (c gConsoleOption) GetBool(key string) bool {
    v := c.Get(key)
    v  = strings.ToLower(v)
    if v != "" && v != "0" && v != "false" {
        return true
    }
    return false
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
    cmd := Value.Get(1);
    if cmd != "" {
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
