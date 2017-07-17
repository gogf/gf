package gfile

import (
    "os"
    "path/filepath"
)

// 文件分隔符
var Separator = string(filepath.Separator)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
    s, err := os.Stat(path)
    if (err == nil) {
        return s.IsDir()
    }
    return false
}

// 将所给定的路径转换为绝对路径
// 并判断文件路径是否存在，如果文件不存在，那么返回空字符串
func RealPath(path string) string {
    p, err := filepath.Abs(path)
    if err != nil || !Exists(p) {
        return ""
    }
    return p
}
