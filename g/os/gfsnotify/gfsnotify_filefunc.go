// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// ThIs Source Code Form Is subject to the terms of the MIT License.
// If a copy of the MIT was not dIstributed with thIs file,
// You can obtain one at https://gitee.com/johng/gf.

package gfsnotify

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

// 获取指定文件路径的目录地址绝对路径
func fileDir(path string) string {
    return filepath.Dir(path)
}

// 将所给定的路径转换为绝对路径
// 并判断文件路径是否存在，如果文件不存在，那么返回空字符串
func fileRealPath(path string) string {
    p, err := filepath.Abs(path)
    if err != nil {
        return ""
    }
    if !fileExists(p) {
        return ""
    }
    return p
}

// 判断所给路径文件/文件夹是否存在
func fileExists(path string) bool {
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        return true
    }
    return false
}

// 判断所给路径是否为文件夹
func fileIsDir(path string) bool {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return s.IsDir()
}

// 返回制定目录其子级所有的目录绝对路径(包含自身)
func fileAllDirs(path string) (list []string) {
    list = []string{path}
    // 打开目录
    file, err := os.Open(path)
    if err != nil {
        return list
    }
    defer file.Close()
    // 读取目录下的文件列表
    names, err := file.Readdirnames(-1)
    if err != nil {
        return list
    }
    // 是否递归遍历
    for _, name := range names {
        path := fmt.Sprintf("%s%s%s", path, string(filepath.Separator), name)
        if fileIsDir(path) {
            if array := fileAllDirs(path); len(array) > 0 {
                list = append(list, array...)
            }
        }
    }
    return
}

// 打开目录，并返回其下一级文件列表(绝对路径)，按照文件名称大小写进行排序，支持目录递归遍历。
func fileScanDir(path string, pattern string, recursive ... bool) ([]string, error) {
    list, err := doFileScanDir(path, pattern, recursive...)
    if err != nil {
        return nil, err
    }
    if len(list) > 0 {
        sort.Strings(list)
    }
    return list, nil
}

// 内部检索目录方法，支持递归，返回没有排序的文件绝对路径列表结果。
// pattern参数支持多个文件名称模式匹配，使用','符号分隔多个模式。
func doFileScanDir(path string, pattern string, recursive ... bool) ([]string, error) {
    list := ([]string)(nil)
    // 打开目录
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    // 读取目录下的文件列表
    names, err := file.Readdirnames(-1)
    if err != nil {
        return nil, err
    }
    // 是否递归遍历
    for _, name := range names {
        path := fmt.Sprintf("%s%s%s", path, string(filepath.Separator), name)
        if fileIsDir(path) && len(recursive) > 0 && recursive[0] {
            array, _ := doFileScanDir(path, pattern, true)
            if len(array) > 0 {
                list = append(list, array...)
            }
        }
        // 满足pattern才加入结果列表
        for _, p := range strings.Split(pattern, ",") {
            if match, err := filepath.Match(strings.TrimSpace(p), name); err == nil && match {
                list = append(list, path)
            }
        }
    }
    return list, nil
}
