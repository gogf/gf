package gfile

import (
    "os"
    "path/filepath"
    "log"
    "io"
    "io/ioutil"
    "sort"
    "fmt"
)

// 封装了常用的文件操作方法，如需更详细的文件控制，请查看官方os包

// 文件分隔符
var Separator = string(filepath.Separator)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
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

// 判断所给路径是否为文件
func IsFile(path string) bool {
    return !IsDir(path)
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

// 修改时间
func MTime(path string) int64 {
    f, e := os.Stat(path)
    if e != nil {
        log.Println(e)
        return 0
    }
    return f.ModTime().Unix()
}

// 文件大小(bytes)
func Size(path string) int64 {
    f, e := os.Stat(path)
    if e != nil {
        log.Println(e)
        return 0
    }
    return f.Size()
}

// 格式化文件大小
func ReadableSize(path string) string {
    return FormatSize(float64(Size(path)))
}

// 格式化文件大小
func FormatSize(raw float64) string {
    var t float64 = 1024
    var d float64 = 1

    if raw < t {
        return fmt.Sprintf("%.2fB", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fK", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fM", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fG", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fT", raw/d)
    }

    d *= 1024
    t *= 1024

    if raw < t {
        return fmt.Sprintf("%.2fP", raw/d)
    }

    return "TooLarge"
}

// 文件移动/重命名
func Move(src string, dst string) error {
    err := os.Rename(src, dst)
    if err != nil {
        log.Println(err)
    }
    return err
}


// 文件移动/重命名
func Rename(src string, dst string) error {
    return Move(src, dst)
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

// 获取当前执行文件的绝对路径
func SelfPath() string {
    p, _ := filepath.Abs(os.Args[0])
    return p
}

// 获取当前执行文件的目录绝对路径
func SelfDir() string {
    return filepath.Dir(SelfPath())
}

// 获取指定文件路径的文件名称
func Basename(path string) string {
    return filepath.Base(path)
}

// 获取指定文件路径的目录地址
func Dir(path string) string {
    return filepath.Dir(path)
}

// 获取指定文件路径的文件扩展名
func Ext(path string) string {
    return filepath.Ext(path)
}