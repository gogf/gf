// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gfile provides easy-to-use operations for file system.
// 
// 文件管理.
package gfile

import (
    "bytes"
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gregex"
    "gitee.com/johng/gf/g/util/gstr"
    "io"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
    "runtime"
    "sort"
    "strings"
    "time"
)

// 文件分隔符
const (
    Separator     = string(filepath.Separator)
    // 默认的文件打开权限
    gDEFAULT_PERM = 0666
)

var (
    // 源码的main包所在目录，仅仅会设置一次
    mainPkgPath   = gtype.NewString()

    // 编译时的 GOROOT 数值
    goRootOfBuild = gtype.NewString()
)

// 给定文件的绝对路径创建文件
func Mkdir(path string) error {
    err  := os.MkdirAll(path, os.ModePerm)
    if err != nil {
        return err
    }
    return nil
}

// 给定文件的绝对路径创建文件
func Create(path string) error {
    dir := Dir(path)
    if !Exists(dir) {
        Mkdir(dir)
    }
    f, err  := os.Create(path)
    if err != nil {
        return err
    }
    f.Close()
    return nil
}

// 打开文件(os.O_RDWR|os.O_CREATE, 0666)
func Open(path string) (*os.File, error) {
    f, err  := os.OpenFile(path, os.O_RDWR|os.O_CREATE, gDEFAULT_PERM)
    if err != nil {
        return nil, err
    }
    return f, nil
}

// 打开文件(带flag)
func OpenWithFlag(path string, flag int) (*os.File, error) {
    f, err  := os.OpenFile(path, flag, gDEFAULT_PERM)
    if err != nil {
        return nil, err
    }
    return f, nil
}

// 打开文件(带flag&perm)
func OpenWithFlagPerm(path string, flag int, perm int) (*os.File, error) {
    f, err  := os.OpenFile(path, flag, os.FileMode(perm))
    if err != nil {
        return nil, err
    }
    return f, nil
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        return true
    }
    return false
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return s.IsDir()
}

// 获取当前工作目录(SelfDir()方法的别名)
func Pwd() string {
    return SelfDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return !s.IsDir()
}

// 获取文件或目录信息
func Info(path string) *os.FileInfo {
    info, err := os.Stat(path)
    if err != nil {
        return nil
    }
    return &info
}

// 文件移动/重命名
func Move(src string, dst string) error {
    return os.Rename(src, dst)
}


// 文件移动/重命名
func Rename(src string, dst string) error {
    return Move(src, dst)
}

// 文件复制
func Copy(src string, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        return err
    }
    err = dstFile.Sync()
    if err != nil {
        return err
    }
    srcFile.Close()
    dstFile.Close()
    return nil
}

// 返回目录下的文件名列表
func DirNames(path string) ([]string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    list, err := f.Readdirnames(-1)
    f.Close()
    if err != nil {
        return nil, err
    }
    return list, nil
}

// 文件名正则匹配查找，第二个可选参数指定返回的列表是否仅为文件名(非绝对路径)，默认返回绝对路径
func Glob(pattern string, onlyNames...bool) ([]string, error) {
    if list, err := filepath.Glob(pattern); err == nil {
        if len(onlyNames) > 0 && onlyNames[0] && len(list) > 0 {
            array := make([]string, len(list))
            for k, v := range list {
                array[k] = Basename(v)
            }
            return array, nil
        }
        return list, nil
    } else {
        return nil, err
    }
}

// 文件/目录删除
func Remove(path string) error {
    return os.RemoveAll(path)
}

// 文件是否可读
func IsReadable(path string) bool {
    result    := true
    file, err := os.OpenFile(path, os.O_RDONLY, gDEFAULT_PERM)
    if err != nil {
        result = false
    }
    file.Close()
    return result
}

// 文件是否可写
func IsWritable(path string) bool {
    result := true
    if IsDir(path) {
        // 如果是目录，那么创建一个临时文件进行写入测试
        tfile := strings.TrimRight(path, Separator) + Separator + gconv.String(time.Now().UnixNano())
        err   := Create(tfile)
        if err != nil || !Exists(tfile){
            result = false
        } else {
            Remove(tfile)
        }
    } else {
        // 如果是文件，那么判断文件是否可打开
        file, err := os.OpenFile(path, os.O_WRONLY, gDEFAULT_PERM)
        if err != nil {
            result = false
        }
        file.Close()
    }
    return result
}

// 修改文件/目录权限
func Chmod(path string, mode os.FileMode) error {
    return os.Chmod(path, mode)
}

// 打开目录，并返回其下一级文件列表(绝对路径)，按照文件名称大小写进行排序，支持目录递归遍历。
func ScanDir(path string, pattern string, recursive ... bool) ([]string, error) {
    list, err := doScanDir(path, pattern, recursive...)
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
func doScanDir(path string, pattern string, recursive ... bool) ([]string, error) {
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
        path := fmt.Sprintf("%s%s%s", path, Separator, name)
        if IsDir(path) && len(recursive) > 0 && recursive[0] {
            array, _ := doScanDir(path, pattern, true)
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

// 将所给定的路径转换为绝对路径
// 并判断文件路径是否存在，如果文件不存在，那么返回空字符串
func RealPath(path string) string {
    p, err := filepath.Abs(path)
    if err != nil {
        return ""
    }
    if !Exists(p) {
        return ""
    }
    return p
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

// 获取指定文件路径的目录地址绝对路径
func Dir(path string) string {
    return filepath.Dir(path)
}

// 获取指定文件路径的文件扩展名(包含"."号)
func Ext(path string) string {
    return filepath.Ext(path)
}

// 获取用户主目录
func Home() (string, error) {
    u, err := user.Current()
    if nil == err {
        return u.HomeDir, nil
    }
    if "windows" == runtime.GOOS {
        return homeWindows()
    }
    return homeUnix()
}

func homeUnix() (string, error) {
    if home := os.Getenv("HOME"); home != "" {
        return home, nil
    }
    var stdout bytes.Buffer
    cmd := exec.Command("sh", "-c", "eval echo ~$USER")
    cmd.Stdout = &stdout
    if err := cmd.Run(); err != nil {
        return "", err
    }

    result := strings.TrimSpace(stdout.String())
    if result == "" {
        return "", errors.New("blank output when reading home directory")
    }

    return result, nil
}

func homeWindows() (string, error) {
    drive := os.Getenv("HOMEDRIVE")
    path  := os.Getenv("HOMEPATH")
    home  := drive + path
    if drive == "" || path == "" {
        home = os.Getenv("USERPROFILE")
    }
    if home == "" {
        return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
    }

    return home, nil
}

// 获取入口函数文件所在目录(main包文件目录),
// **仅对源码开发环境有效(即仅对生成该可执行文件的系统下有效)**
func MainPkgPath() string {
    path := mainPkgPath.Val()
    if path != "" {
        return path
    }
    f      := ""
    goroot := runtime.GOROOT()
    // runtime.GOROOT() 在windows下有可能是以'\'符号分隔，
    // 而 runtime.Caller(i) 获取到的文件路径却是以'/'符号分隔，
    // 因此这里统一转换为'/'符号再进行比较
    goroot  = gstr.Replace(goroot, "\\", "/")
    for i := 1; i < 10000; i++ {
        if _, file, _, ok := runtime.Caller(i); ok {
            // 不包含go源码路径
            if file != "" && goroot != "" &&
                !gregex.IsMatchString("^" + goroot, file) &&
                !strings.EqualFold("<autogenerated>", file) {
                f = file
            }
        } else {
            break
        }
    }
    if f != "" {
        for {
            p := Dir(f)
            if p == f {
                break
            }
            // 会自动扫描源码，寻找main包
            if paths, err := ScanDir(p, "*.go"); err == nil && len(paths) > 0 {
                for _, path := range paths {
                    if gregex.IsMatchString(`package\s+main`, GetContents(path)) {
                        mainPkgPath.Set(p)
                        return p
                    }
                }
            }
            f = p
        }
    }
    return ""
}

// 系统临时目录
func TempDir() string {
    return os.TempDir()
}