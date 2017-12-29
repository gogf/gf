// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 全局配置管理对象
package gconfig

import (
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/encoding/gjson"
)

const (
    gDEFAULT_CONFIG_FILE = "config" // 默认的配置管理文件名称
)

// 配置管理对象
type Config struct {
    path  string                   // 配置文件存放目录，绝对路径
    jsons *gmap.StringInterfaceMap // 配置文件对象
}

// 生成一个配置管理对象
func New(path string) *Config {
    return &Config{
        path  : path,
        jsons : gmap.NewStringInterfaceMap(),
    }
}

// 判断从哪个配置文件中获取内容
func (c *Config) file(files []string) string {
    file := gDEFAULT_CONFIG_FILE
    if len(files) > 0 {
        file = files[0]
    }
    return file + ".json"
}

// 添加配置文件到配置管理器中，第二个参数为非必须，如果不输入表示添加进入默认的配置名称中
func (c *Config) getJson(files []string) *gjson.Json {
    file := c.file(files)
    if r := c.jsons.Get(file); r != nil {
        return r.(*gjson.Json)
    }
    path := c.path + gfile.Separator + file
    if j, err := gjson.Load(path); err == nil {
        c.jsons.Set(file, j)
        return j
    }
    return nil
}

// 获取配置项，当不存在时返回nil
func (c *Config) Get(pattern string, files...string) interface{} {
    if j := c.getJson(files); j != nil {
        return j.Get(pattern)
    }
    return nil
}

// 获得一个键值对关联数组/哈希表，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (c *Config) GetMap(pattern string, files...string)  map[string]interface{} {
    if j := c.getJson(files); j != nil {
        return j.GetMap(pattern)
    }
    return nil
}

// 获得一个数组[]interface{}，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (c *Config) GetArray(pattern string, files...string)  []interface{} {
    if j := c.getJson(files); j != nil {
        return j.GetArray(pattern)
    }
    return nil
}

// 返回指定json中的string
func (c *Config) GetString(pattern string, files...string) string {
    if j := c.getJson(files); j != nil {
        return j.GetString(pattern)
    }
    return ""
}

// 返回指定json中的bool
func (c *Config) GetBool(pattern string, files...string) bool {
    if j := c.getJson(files); j != nil {
        return j.GetBool(pattern)
    }
    return false
}

// 返回指定json中的float64
func (c *Config) GetFloat64(pattern string, files...string) float64 {
    if j := c.getJson(files); j != nil {
        return j.GetFloat64(pattern)
    }
    return 0
}

// 返回指定json中的float64->int
func (c *Config) GetInt(pattern string, files...string)  int {
    return int(c.GetFloat64(pattern))
}

// 返回指定json中的float64->int64
func (c *Config) GetInt64(pattern string, files...string)  int64 {
    return int64(c.GetFloat64(pattern))
}
