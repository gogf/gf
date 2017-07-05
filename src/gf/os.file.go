package gf

import (
    "os"
    "path/filepath"
)

var File gstFile

type gstFile struct {
    Separator string
}

func init() {
    File.Separator = string(filepath.Separator)
}

// 判断所给路径文件/文件夹是否存在
func (f gstFile) Exists(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}

// 判断所给路径是否为文件夹
func (f gstFile) IsDir(path string) bool {
    s, err := os.Stat(path)
    if (err == nil) {
        return s.IsDir()
    }
    return false
}
