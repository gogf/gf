// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package grouter

import (
    "sync"
    "sort"
    "bytes"
    "errors"
    "strings"
    "gitee.com/johng/gf/g/util/gregx"
    "gitee.com/johng/gf/g/container/gmap"
)

// 路由管理对象
type Router struct {
    dmu    sync.RWMutex // 解析规则互斥锁
    pmu    sync.RWMutex // 打包规则互斥锁
    dkeys  []string     // 解析规则排序键名
    pkeys  []string     // 打包规则排序键名
    drules *gmap.StringStringMap // 解析规则
    prules *gmap.StringStringMap // 打包规则
}


func New() *Router {
    return &Router{
        drules : gmap.NewStringStringMap(),
        prules : gmap.NewStringStringMap(),
    }
}

// 设置解析规则，例如：静态分页
// `\/([\w\.\-]+)\/([\w\.\-]+)\/page\/([\d\.\-]+)[\/\?]*`, "/user/list/page/2"
func (r *Router) SetRule(rule, replace string) {
    r.drules.Set(rule, replace)
    r.updateDispatchKeys()
}

// 批量设置解析规则
func (r *Router) SetRules(rules map[string]string) {
    r.drules.BatchSet(rules)
    r.updateDispatchKeys()
}

// 删除解析规则
func (r *Router) RemoveRule(rule string) {
    r.drules.Remove(rule)
    r.updateDispatchKeys()
}

// 设置打包规则
func (r *Router) SetPatchRule(rule, replace string) {
    r.prules.Set(rule, replace)
    r.updatePatchKeys()
}

// 批量设置打包规则
func (r *Router) SetPatchRules(rules map[string]string) {
    r.prules.BatchSet(rules)
    r.updatePatchKeys()
}

// 删除打包规则
func (r *Router) RemovePatchRule(rule string) {
    r.prules.Remove(rule)
    r.updatePatchKeys()
}

func (r *Router) updateDispatchKeys() {
    r.dmu.Lock()
    defer r.dmu.Unlock()
    r.dkeys = r.drules.Keys()
    sort.Slice(r.dkeys, func(i, j int) bool { return len(r.dkeys[i]) > len(r.dkeys[j]) })
}

func (r *Router) updatePatchKeys() {
    r.pmu.Lock()
    defer r.pmu.Unlock()
    r.pkeys = r.prules.Keys()
    sort.Slice(r.pkeys, func(i, j int) bool { return len(r.pkeys[i]) > len(r.pkeys[j]) })
}

// 解析URI
func (r *Router) Dispatch(uri string) (string, error) {
    r.dmu.RLock()
    defer r.dmu.RUnlock()
    if len(r.dkeys) == 0 {
        return uri, errors.New("no dispatch rules found")
    }
    for _, rule := range r.dkeys {
        if replace := r.drules.Get(rule); replace != "" {
            result, err := gregx.ReplaceString(rule, uri, replace)
            if err != nil {
                return result, err
            }
            if len(uri) != len(result) || strings.Compare(result, uri) != 0 {
                return result, nil
            }
        }
    }
    return uri, nil
}

// 打包内容
func (r *Router) Patch(content []byte) ([]byte, error) {
    r.pmu.RLock()
    defer r.pmu.RUnlock()
    if len(r.pkeys) == 0 {
        return content, errors.New("no patch rules found")
    }
    for _, rule := range r.pkeys {
        if replace := r.prules.Get(rule); replace != "" {
            result, err := gregx.Replace(rule, content, []byte(replace))
            if err != nil {
                return result, err
            }
            if len(content) != len(result) || bytes.Compare(result, content) != 0 {
                return result, nil
            }
        }
    }
    return content, nil
}