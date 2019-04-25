// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gfile provides easy-to-use operations for file system.
// 
// 文件管理.
package gfile

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/text/gregex"
	"github.com/gogf/gf/g/util/gconv"
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

const (
    // 文件分隔符
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

// Create directories recursively.
//
// 给定目录的绝对路径创建目录(递归创建)。
func Mkdir(path string) error {
    err  := os.MkdirAll(path, os.ModePerm)
    if err != nil {
        return err
    }
    return nil
}

// Create file with given path recursively.
//
// 给定文件的绝对路径创建文件。
func Create(path string) (*os.File, error) {
    dir := Dir(path)
    if !Exists(dir) {
        Mkdir(dir)
    }
    return os.Create(path)
}

// Open file/directory with readonly.
//
// 只读打开文件
func Open(path string) (*os.File, error) {
    return os.Open(path)
}

// Open file/directory with given <flag> and <perm>.
//
// 打开文件(带flag&perm)
func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
    return os.OpenFile(path, flag, perm)
}

// Open file/directory with default perm and given <flag>.
//
// 打开文件(带flag)
func OpenWithFlag(path string, flag int) (*os.File, error) {
    f, err  := os.OpenFile(path, flag, gDEFAULT_PERM)
    if err != nil {
        return nil, err
    }
    return f, nil
}

// Open file/directory with given <flag> and <perm>.
//
// 打开文件(带flag&perm)
func OpenWithFlagPerm(path string, flag int, perm int) (*os.File, error) {
    f, err  := os.OpenFile(path, flag, os.FileMode(perm))
    if err != nil {
        return nil, err
    }
    return f, nil
}

// Check whether given path exist.
//
// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        return true
    }
    return false
}

// Check whether given path a directory.
//
// 判断所给路径是否为文件夹
func IsDir(path string) bool {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return s.IsDir()
}

// Get current working directory absolute path.
//
// 获取当前工作目录(注意与SelfDir的区别).
func Pwd() string {
    path, _ := os.Getwd()
    return path
}

// Check whether given path a file(not a directory).
//
// 判断所给路径是否为文件
func IsFile(path string) bool {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return !s.IsDir()
}

// See Stat.
//
// Stat 方法的别名。
func Info(path string) (os.FileInfo, error) {
    return Stat(path)
}

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
//
// 获取文件或目录信息.
func Stat(path string) (os.FileInfo, error) {
    return os.Stat(path)
}

// Move renames (moves) src to dst path.
//
// 文件移动/重命名
func Move(src string, dst string) error {
    return os.Rename(src, dst)
}

// Rename renames (moves) src to dst path.
//
// 文件移动/重命名.
func Rename(src string, dst string) error {
    return Move(src, dst)
}

// Copy file from src to dst.
//
// 文件复制.
// @TODO 支持目录复制.
func Copy(src string, dst string) error {
    srcFile, err := Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    dstFile, err := Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()
    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        return err
    }
    err = dstFile.Sync()
    if err != nil {
        return err
    }
    return nil
}

// Get sub-file names of path.
//
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

// Glob returns the names of all files matching pattern or nil
// if there is no matching file. The syntax of patterns is the same
// as in Match. The pattern may describe hierarchical names such as
// /usr/*/bin/ed (assuming the Separator is '/').
//
// Glob ignores file system errors such as I/O errors reading directories.
// The only possible returned error is ErrBadPattern, when pattern
// is malformed.
//
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

// Remove file/directory with <path> parameter.
//
// 文件/目录删除
func Remove(path string) error {
    return os.RemoveAll(path)
}

// Check whether given <path> is readable.
//
// 文件是否可读(支持文件/目录)
func IsReadable(path string) bool {
    result    := true
    file, err := os.OpenFile(path, os.O_RDONLY, gDEFAULT_PERM)
    if err != nil {
        result = false
    }
    file.Close()
    return result
}

// Check whether given <path> is writable.
//
// 文件是否可写(支持文件/目录)
// @TODO 改进性能，利用 golang.org/x/sys 来实现跨平台的权限判断。
func IsWritable(path string) bool {
    result := true
    if IsDir(path) {
        // 如果是目录，那么创建一个临时文件进行写入测试
        tmpFile := strings.TrimRight(path, Separator) + Separator + gconv.String(time.Now().UnixNano())
        if f, err := Create(tmpFile); err != nil || !Exists(tmpFile){
            result = false
        } else {
            f.Close()
            Remove(tmpFile)
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

// See os.Chmod.
//
// 修改文件/目录权限
func Chmod(path string, mode os.FileMode) error {
    return os.Chmod(path, mode)
}

// Get all sub-files(absolute) of given <path>,
// can be recursively with given parameter <recursive> true.
//
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

// See filepath.Abs.
//
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

// Get absolute file path of current running process(binary).
//
// 获取当前执行文件的绝对路径
func SelfPath() string {
    p, _ := filepath.Abs(os.Args[0])
    return p
}

// Get absolute directory path of current running process(binary).
//
// 获取当前执行文件的目录绝对路径
func SelfDir() string {
    return filepath.Dir(SelfPath())
}

// See filepath.Base.
//
// 获取指定文件路径的文件名称
func Basename(path string) string {
    return filepath.Base(path)
}

// See filepath.Dir.
//
// 获取指定文件路径的目录地址绝对路径.
func Dir(path string) string {
    return filepath.Dir(path)
}

// See filepath.Ext.
//
// 获取指定文件路径的文件扩展名(包含"."号)
func Ext(path string) string {
    return filepath.Ext(path)
}

// Get absolute home directory path of current user.
//
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

// Get absolute file path of main file, which contains the entrance function main.
// Available in develop environment.
//
// 获取入口函数文件所在目录(main包文件目录),
// **仅对源码开发环境有效(即仅对生成该可执行文件的系统下有效)**。
// 注意：该方法被第一次调用时，如果是在异步的goroutine中，该方法可能无法获取到main包路径。
func MainPkgPath() string {
    path := mainPkgPath.Val()
    if path != "" {
    	if path == "-" {
    		return ""
	    }
        return path
    }
    for i := 1; i < 10000; i++ {
        if _, file, _, ok := runtime.Caller(i); ok {
	        if gregex.IsMatchString(`package\s+main`, GetContents(file)) {
	        	path = Dir(file)
		        mainPkgPath.Set(path)
		        return path
	        }
        } else {
            break
        }
    }
    // 找不到，下次不用再检索了
	mainPkgPath.Set("-")
    return ""
}

// See os.TempDir().
//
// 系统临时目录
func TempDir() string {
    return os.TempDir()
}
