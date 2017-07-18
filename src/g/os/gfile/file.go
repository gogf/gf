package gfile

import (
    "os"
    "path/filepath"
    "log"
    "io"
    "io/ioutil"
    "sort"
)

// 封装了常用的文件操作方法，如需更详细的文件控制，请查看官方os包

// 文件分隔符
var Separator = string(filepath.Separator)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        log.Println(err)
        return false
    }
    return os.IsExist(err)
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
    s, err := os.Stat(path)
    if (err != nil) {
        log.Println(err)
        return false
    }
    return s.IsDir()
}

// 获取文件或目录信息
func Info(path string) os.FileInfo {
    info, err := os.Stat(path)
    if err != nil {
        log.Println(err)
        return nil
    }
    return info
}

// 文件移动/重命名
func Move(src string, dst string) {
    err := os.Rename(src, dst)
    if err != nil {
        log.Println(err)
    }
}

// 文件复制
func Copy(src string, dst string) bool {
    result       := true
    srcFile, err := os.Open(src)
    if err != nil {
        result = false
        log.Println(err)
    }
    dstFile, err := os.Create(dst)
    if err != nil {
        result = false
        log.Println(err)
    }
    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        result = false
        log.Println(err)
    }
    err = dstFile.Sync()
    if err != nil {
        result = false
        log.Println(err)
    }
    srcFile.Close()
    dstFile.Close()
    return result
}

// 文件删除
func Remove(path string) {
    err := os.Remove(path)
    if err != nil {
        log.Println(err)
    }
}

// 文件是否可读
func Readable(path string) bool {
    result    := true
    file, err := os.OpenFile(path, os.O_RDONLY, 0666)
    if err != nil {
        log.Println(err)
        result = false
    }
    file.Close()
    return result
}

// 文件是否可写
func Writable(path string) bool {
    result    := true
    file, err := os.OpenFile(path, os.O_WRONLY, 0666)
    if err != nil {
        log.Println(err)
        result = false
    }
    file.Close()
    return result
}

// 修改文件/目录权限
func Chmod(path string, mode os.FileMode) bool {
    result := true
    err    := os.Chmod(path, mode)
    if err != nil {
        log.Println(err)
        result = false
    }
    return result
}

// 打开目录，并返回其下一级子目录名称列表，按照文件名称大小写进行排序
func ScanDir(path string) []string {
    f, err := os.Open(path)
    if err != nil {
        log.Println(err)
        return nil
    }
    list, err := f.Readdirnames(-1)
    f.Close()
    if err != nil {
        log.Println(err)
        return nil
    }
    sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
    return list
}

// 将所给定的路径转换为绝对路径
// 并判断文件路径是否存在，如果文件不存在，那么返回空字符串
func RealPath(path string) string {
    p, err := filepath.Abs(path)
    if err != nil {
        log.Println(err)
        return ""
    }
    if !Exists(p) {
        return ""
    }
    return p
}

// 读取文件内容
func GetContents(path string) []byte {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        log.Println(err)
        return nil
    }
    return data
}

// 写入文件内容
func putContents(path string, data []byte, flag int, perm os.FileMode) bool {
    result := true
    f, err := os.OpenFile(path, flag, perm)
    if err == nil {
        n, err := f.Write(data)
        if err == nil && n < len(data) {
            err = io.ErrShortWrite
        }
        if err1 := f.Close(); err == nil {
            err = err1
        }
    }
    if err != nil {
        log.Println(err)
        result = false
    }
    return result
}

// 写入文件内容
func PutContents(path string, content string) bool {
    return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
}

// 追加内容到文件末尾
func PutContentsAppend(path string, content string) bool {
    return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
}