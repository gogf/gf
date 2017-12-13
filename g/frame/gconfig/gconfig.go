// 全局配置管理对象
package gconfig

import (
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/os/gfile"
    "errors"
)

// 配置对象
var config = gmap.NewStringInterfaceMap()

// 获取配置
func Get(k string) interface{} {
    return config.Get(k)
}

func GetInt(k string) int {
    if v := config.Get(k); v != nil {
        if r, ok := v.(int); ok {
            return r
        }
    }
    return 0
}

func GetString(k string) string {
    if v := config.Get(k); v != nil {
        if r, ok := v.(string); ok {
            return r
        }
    }
    return ""
}

// 适用于json文件配置，在设置的时候通过gjson进行解析后再保存
func GetJson(k string) *gjson.Json {
    if v := config.Get(k); v != nil {
        if r, ok := v.(*gjson.Json); ok {
            return r
        }
    }
    return nil
}

// 设置配置
func Set(k string, v interface{}) {
    config.Set(k, v)
}

// 加载json文件配置
func Load(key string, path string) error {
    content := gfile.GetContents(path)
    if len(content) == 0 {
        return errors.New("load json file failed, path: " + path)
    }
    if json, err := gjson.DecodeToJson(content); err == nil {
        config.Set(key, json)
    } else {
        return err
    }
    return nil
}