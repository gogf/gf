// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 配置管理.
// 配置文件格式支持：json, xml, toml, yaml/yml
package gcfg

import (
    "errors"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/container/gvar"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/os/gfsnotify"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gspath"
)

const (
    DEFAULT_CONFIG_FILE = "config.toml" // 默认的配置管理文件名称
)

// 配置管理对象
type Config struct {
    name   *gtype.String            // 默认配置文件名称
    paths  *gspath.SPath            // 搜索目录路径
    jsons  *gmap.StringInterfaceMap // 配置文件对象
    vc     *gtype.Bool              // 层级检索是否执行分隔符冲突检测(默认为false，检测会比较影响检索效率)
}

// 生成一个配置管理对象
func New(path string, file...string) *Config {
    name := DEFAULT_CONFIG_FILE
    if len(file) > 0 {
        name = file[0]
    }
    c := &Config {
        name   : gtype.NewString(name),
        paths  : gspath.New(),
        jsons  : gmap.NewStringInterfaceMap(),
        vc     : gtype.NewBool(),
    }
    c.SetPath(path)
    return c
}

// 判断从哪个配置文件中获取内容，返回配置文件的绝对路径
func (c *Config) filePath(file...string) string {
    name := c.name.Val()
    if len(file) > 0 {
        name = file[0]
    }
    return c.paths.Search(name)
}

// 设置配置管理器的配置文件存放目录绝对路径
func (c *Config) SetPath(path string) error {
    if rp, err := c.paths.Set(path); err != nil {
        glog.Error("gcfg.SetPath failed:", err.Error())
        return err
    } else {
        c.jsons.Clear()
        glog.Debug("gcfg.SetPath:", rp)
    }
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
    if rp, err := c.paths.Add(path); err != nil {
        glog.Debug("gcfg.AddPath failed:", err.Error())
        return err
    } else {
        glog.Debug("gcfg.AddPath:", rp)
    }
    return nil
}

// 获取指定文件的绝对路径，默认获取默认的配置文件路径
func (c *Config) GetFilePath(file...string) string {
    name := c.name.Val()
    if len(file) > 0 {
        name = file[0]
    }
    return c.paths.Search(name)
}

// 设置配置管理对象的默认文件名称
func (c *Config) SetFileName(name string) {
    glog.Debug("gcfg.SetFileName:", name)
    c.name.Set(name)
}

// 添加配置文件到配置管理器中，第二个参数为非必须，如果不输入表示添加进入默认的配置名称中
func (c *Config) getJson(file...string) *gjson.Json {
    fpath := c.filePath(file...)
    if r := c.jsons.Get(fpath); r != nil {
        return r.(*gjson.Json)
    }
    if j, err := gjson.Load(fpath); err == nil {
        j.SetViolenceCheck(c.vc.Val())
        c.addMonitor(fpath)
        c.jsons.Set(fpath, j)
        return j
    } else {
        glog.Errorfln(`gcfg.Load config file "%s" failed: %s`, fpath, err.Error())
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
func (c *Config) GetVar(pattern string, file...string) *gvar.Var {
    if j := c.getJson(file...); j != nil {
        return gvar.New(j.Get(pattern))
    }
    return nil
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

// 添加文件监控
func (c *Config) addMonitor(path string) {
    // 防止多goroutine同时调用
    if c.jsons.Get(path) != nil {
        return
    }
    gfsnotify.Add(path, func(event *gfsnotify.Event) {
        // 删除文件内容缓存，下一次查询会自动更新
        c.jsons.Remove(event.Path)
    })
}
