// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcfg provides reading, caching and managing for configuration files.
// 
// 配置管理,
// 配置文件格式支持：json, xml, toml, yaml/yml
package gcfg

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/container/garray"
    "github.com/gogf/gf/g/container/gmap"
    "github.com/gogf/gf/g/container/gtype"
    "github.com/gogf/gf/g/container/gvar"
    "github.com/gogf/gf/g/encoding/gjson"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/os/gfsnotify"
    "github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/os/gspath"
)

const (
    // 默认的配置管理文件名称
    DEFAULT_CONFIG_FILE = "config.toml"
)

// 配置管理对象
type Config struct {
    name   *gtype.String            // 默认配置文件名称
    paths  *garray.StringArray      // 搜索目录路径
    jsons  *gmap.StringInterfaceMap // 配置文件对象
    vc     *gtype.Bool              // 层级检索是否执行分隔符冲突检测(默认为false，检测会比较影响检索效率)
}

// New returns a new configuration management object.
//
// 生成一个配置管理对象
func New(path string, file...string) *Config {
    name := DEFAULT_CONFIG_FILE
    if len(file) > 0 {
        name = file[0]
    }
    c := &Config {
        name   : gtype.NewString(name),
        paths  : garray.NewStringArray(),
        jsons  : gmap.NewStringInterfaceMap(),
        vc     : gtype.NewBool(),
    }
    if len(path) > 0 {
        c.SetPath(path)
    }
    return c
}

// filePath returns the absolute configuration file path for the given filename by <file>.
//
// 判断从哪个配置文件中获取内容，返回配置文件的绝对路径
func (c *Config) filePath(file...string) (path string) {
    name := c.name.Val()
    if len(file) > 0 {
        name = file[0]
    }
    path = c.GetFilePath(name)
    if path == "" {
        buffer := bytes.NewBuffer(nil)
        if c.paths.Len() > 0 {
            buffer.WriteString(fmt.Sprintf("[gcfg] cannot find config file \"%s\" in following paths:", name))
            c.paths.RLockFunc(func(array []string) {
                for k, v := range array {
                    buffer.WriteString(fmt.Sprintf("\n%d. %s",k + 1,  v))
                }
            })
        } else {
            buffer.WriteString(fmt.Sprintf("[gcfg] cannot find config file \"%s\" with no path set/add", name))
        }
        glog.Error(buffer.String())
    }
    return path
}

// 设置配置管理器的配置文件存放目录绝对路径
func (c *Config) SetPath(path string) error {
    realPath := gfile.RealPath(path)
    if realPath == "" {
        err := errors.New(fmt.Sprintf(`path "%s" does not exist`, path))
        glog.Error(fmt.Sprintf(`[gcfg] SetPath failed: %s`, err.Error()))
        return err
    }
    // 重复判断
    if c.paths.Search(realPath) != -1 {
        return nil
    }
    c.jsons.Clear()
    c.paths.Clear()
    c.paths.Append(realPath)
    //glog.Debug("[gcfg] SetPath:", realPath)
    return nil
}

// 设置是否执行层级冲突检查，当键名中存在层级符号时需要开启该特性，默认为关闭。
// 开启比较耗性能，也不建议允许键名中存在分隔符，最好在应用端避免这种情况。
func (c *Config) SetViolenceCheck(check bool) {
    c.vc.Set(check)
    c.Reload()
}

// 添加配置管理器的配置文件搜索路径
func (c *Config) AddPath(path string) error {
    realPath := gfile.RealPath(path)
    if realPath == "" {
        err := errors.New(fmt.Sprintf(`path "%s" does not exist`, path))
        glog.Error(fmt.Sprintf(`[gcfg] AddPath failed: %s`, err.Error()))
        return err
    }
    // 重复判断
    if c.paths.Search(realPath) != -1 {
        return nil
    }
    c.paths.Append(realPath)
    //glog.Debug("[gcfg] AddPath:", realPath)
    return nil
}

// 查找配置文件，获取指定配置文件的绝对路径，默认获取默认的配置文件路径；
// 当指定的配置文件不存在时，返回空字符串，并且不会报错。
func (c *Config) GetFilePath(file...string) (path string) {
    name := c.name.Val()
    if len(file) > 0 {
        name = file[0]
    }
    c.paths.RLockFunc(func(array []string) {
        for _, v := range array {
            // 查找当前目录
            if path, _ = gspath.Search(v, name); path != "" {
                break
            }
            // 查找当前目录下的config子目录
            if path, _ = gspath.Search(v, "config" + gfile.Separator + name); path != "" {
                break
            }
        }
    })
    return
}

// 设置配置管理对象的默认文件名称
func (c *Config) SetFileName(name string) {
    //glog.Debug("[gcfg] SetFileName:", name)
    c.name.Set(name)
}

// 获取配置管理对象的默认文件名称
func (c *Config) GetFileName() string {
    return c.name.Val()
}

// 添加配置文件到配置管理器中，第二个参数为非必须，如果不输入表示添加进入默认的配置名称中
// 内部带缓存控制功能。
func (c *Config) getJson(file...string) *gjson.Json {
    name := c.name.Val()
    if len(file) > 0 {
        name = file[0]
    }
    r := c.jsons.GetOrSetFuncLock(name, func() interface{} {
        filePath := c.filePath(file...)
        if filePath == "" {
            return nil
        }
        if j, err := gjson.Load(filePath); err == nil {
            j.SetViolenceCheck(c.vc.Val())
            // 添加配置文件监听，如果有任何变化，删除文件内容缓存，下一次查询会自动更新
            gfsnotify.Add(filePath, func(event *gfsnotify.Event) {
                c.jsons.Remove(name)
            })
            return j
        } else {
            glog.Criticalfln(`[gcfg] Load config file "%s" failed: %s`, filePath, err.Error())
        }
        return nil
    })
    if r != nil {
        return r.(*gjson.Json)
    }
    return nil
}

// 获取配置项，当不存在时返回nil
func (c *Config) Get(pattern string, file...string) interface{} {
    if j := c.getJson(file...); j != nil {
        return j.Get(pattern)
    }
    return nil
}

// 获得配置项，返回动态变量
func (c *Config) GetVar(pattern string, file...string) gvar.VarRead {
    if j := c.getJson(file...); j != nil {
        return gvar.New(j.Get(pattern), true)
    }
    return gvar.New(nil, true)
}

// 判断指定的配置项是否存在
func (c *Config) Contains(pattern string, file...string) bool {
    if j := c.getJson(file...); j != nil {
        return j.Contains(pattern)
    }
    return false
}

// 获得一个键值对关联数组/哈希表，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (c *Config) GetMap(pattern string, file...string)  map[string]interface{} {
    if j := c.getJson(file...); j != nil {
        return j.GetMap(pattern)
    }
    return nil
}

// 获得一个数组[]interface{}，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (c *Config) GetArray(pattern string, file...string)  []interface{} {
    if j := c.getJson(file...); j != nil {
        return j.GetArray(pattern)
    }
    return nil
}

// 返回指定json中的string
func (c *Config) GetString(pattern string, file...string) string {
    if j := c.getJson(file...); j != nil {
        return j.GetString(pattern)
    }
    return ""
}

func (c *Config) GetStrings(pattern string, file...string) []string {
    if j := c.getJson(file...); j != nil {
        return j.GetStrings(pattern)
    }
    return nil
}

func (c *Config) GetInterfaces(pattern string, file...string) []interface{} {
    if j := c.getJson(file...); j != nil {
        return j.GetInterfaces(pattern)
    }
    return nil
}

// 返回指定json中的bool
func (c *Config) GetBool(pattern string, file...string) bool {
    if j := c.getJson(file...); j != nil {
        return j.GetBool(pattern)
    }
    return false
}

// 返回指定json中的float32
func (c *Config) GetFloat32(pattern string, file...string) float32 {
    if j := c.getJson(file...); j != nil {
        return j.GetFloat32(pattern)
    }
    return 0
}

// 返回指定json中的float64
func (c *Config) GetFloat64(pattern string, file...string) float64 {
    if j := c.getJson(file...); j != nil {
        return j.GetFloat64(pattern)
    }
    return 0
}

func (c *Config) GetFloats(pattern string, file...string) []float64 {
    if j := c.getJson(file...); j != nil {
        return j.GetFloats(pattern)
    }
    return nil
}

// 返回指定json中的float64->int
func (c *Config) GetInt(pattern string, file...string)  int {
    if j := c.getJson(file...); j != nil {
        return j.GetInt(pattern)
    }
    return 0
}


func (c *Config) GetInt8(pattern string, file...string)  int8 {
    if j := c.getJson(file...); j != nil {
        return j.GetInt8(pattern)
    }
    return 0
}

func (c *Config) GetInt16(pattern string, file...string)  int16 {
    if j := c.getJson(file...); j != nil {
        return j.GetInt16(pattern)
    }
    return 0
}

func (c *Config) GetInt32(pattern string, file...string)  int32 {
    if j := c.getJson(file...); j != nil {
        return j.GetInt32(pattern)
    }
    return 0
}

func (c *Config) GetInt64(pattern string, file...string)  int64 {
    if j := c.getJson(file...); j != nil {
        return j.GetInt64(pattern)
    }
    return 0
}

func (c *Config) GetInts(pattern string, file...string) []int {
    if j := c.getJson(file...); j != nil {
        return j.GetInts(pattern)
    }
    return nil
}

// 返回指定json中的float64->uint
func (c *Config) GetUint(pattern string, file...string)  uint {
    if j := c.getJson(file...); j != nil {
        return j.GetUint(pattern)
    }
    return 0
}

func (c *Config) GetUint8(pattern string, file...string)  uint8 {
    if j := c.getJson(file...); j != nil {
        return j.GetUint8(pattern)
    }
    return 0
}

func (c *Config) GetUint16(pattern string, file...string)  uint16 {
    if j := c.getJson(file...); j != nil {
        return j.GetUint16(pattern)
    }
    return 0
}

func (c *Config) GetUint32(pattern string, file...string)  uint32 {
    if j := c.getJson(file...); j != nil {
        return j.GetUint32(pattern)
    }
    return 0
}

func (c *Config) GetUint64(pattern string, file...string)  uint64 {
    if j := c.getJson(file...); j != nil {
        return j.GetUint64(pattern)
    }
    return 0
}

func (c *Config) GetToStruct(pattern string, objPointer interface{}, file...string) error {
    if j := c.getJson(file...); j != nil {
        return j.GetToStruct(pattern, objPointer)
    }
    return errors.New("config file not found")
}

// 清空当前配置文件缓存，强制重新从磁盘文件读取配置文件内容
func (c *Config) Reload() {
    c.jsons.Clear()
}

