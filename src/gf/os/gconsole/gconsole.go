package gconsole

import (
	"os"
    "regexp"
    "errors"
)

// 包变量结构体，存储console option解析数据
// console option命令行参数
var co struct {
    inited  bool
    values  []string
    options map[string]string
}

// 检查并初始化console参数
func init() {
    if !co.inited {
        co.options = make(map[string]string)
        reg       := regexp.MustCompile(`\-\-{0,1}(\w+?)=(.+)`)
        for i := 1; i < len(os.Args); i++ {
            result := reg.FindStringSubmatch(os.Args[i])
            if len(result) > 1 {
                co.options[result[1]] = result[2]
            } else {
                co.values = append(co.values, os.Args[i])
            }
        }
        co.inited = true
    }
}

// 返回所有的命令行参数values
func GetValues() []string {
    return co.values
}

// 返回所有的命令行参数options
func GetOptions() map[string]string {
    return co.options
}

// 获得一条指定索引位置的value参数
func GetValue(index uint8) (string, error) {
    if index < uint8(len(co.values)) {
        return co.values[index], nil
    }
    return "", errors.New("index out of range")
}

// 获得一条指定索引位置的option参数
func GetOption(key string) (string, error) {
    option, ok := co.options[key]
    if ok {
        return option, nil
    }
    return "", errors.New("no mapping value")
}
