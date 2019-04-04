// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/container/garray"
)


// 如果给定绝对路径将会去掉其中的相对路径符号后返回；
// 如果是给定的相对路径，那么将会按照以下路径优先级搜索文件(重复路径会去重)：
// prioritySearchPaths、当前工作目录、二进制文件目录、源码main包目录(开发环境下)
func Search(name string, prioritySearchPaths...string) (realPath string, err error) {
    // 是否绝对路径
    realPath = RealPath(name)
    if realPath != "" {
        return
    }
    // 相对路径搜索
    array := garray.NewStringArray(true)
    // 自定义优先路径
    array.Append(prioritySearchPaths...)
    // 用户工作目录
    array.Append(Pwd())
    // 二进制所在目录
    array.Append(SelfDir())
    // 源码main包目录
    if path := MainPkgPath(); path != "" {
        array.Append(path)
    }
    // 路径去重
    array.Unique()
    // 执行相对路径搜索
    array.RLockFunc(func(array []string) {
        path := ""
        for _, v := range array {
            path = RealPath(v + Separator + name)
            if path != "" {
                realPath = path
                break
            }
        }
    })
    // 目录不存在错误处理
    if realPath == "" {
        buffer := bytes.NewBuffer(nil)
        buffer.WriteString(fmt.Sprintf("cannot find file/folder \"%s\" in following paths:", name))
        array.RLockFunc(func(array []string) {
            for k, v := range array {
                buffer.WriteString(fmt.Sprintf("\n%d. %s", k + 1,  v))
            }
        })
        err = errors.New(buffer.String())
    }
    return
}
