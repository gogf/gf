// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gp.

// 数据文件编码/解析.
package gparser

import (
    "gitee.com/johng/gf/g/encoding/gjson"
)

type Parser struct {
    json *gjson.Json
}

func New (values...interface{}) *Parser {
    if len(values) > 0 {
        return &Parser{gjson.NewJson(values[0])}
    }
    return &Parser{gjson.NewJson(nil)}
}

func Load (path string) (*Parser, error) {
    if j, e := gjson.Load(path); e == nil {
        return &Parser{j}, nil
    } else {
        return nil, e
    }
}

// 支持的配置文件格式：xml, json, yaml/yml, toml
func LoadContent (data []byte, fileType string) (*Parser, error) {
    if j, e := gjson.LoadContent(data, fileType); e == nil {
        return &Parser{j}, nil
    } else {
        return nil, e
    }
}

// 将指定的json内容转换为指定结构返回，查找失败或者转换失败，目标对象转换为nil
// 注意第二个参数需要给的是变量地址
func (p *Parser) GetToVar(pattern string, v interface{}) error {
    return p.json.GetToVar(pattern, v)
}

// 获得一个键值对关联数组/哈希表，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (p *Parser) GetMap(pattern string) map[string]interface{} {
    return p.json.GetMap(pattern)
}

// 获得一个数组[]interface{}，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (p *Parser) GetArray(pattern string) []interface{} {
    return p.json.GetArray(pattern)
}

// 返回指定json中的string
func (p *Parser) GetString(pattern string) string {
    return p.json.GetString(pattern)
}

// 返回指定json中的bool(false:"", 0, false, off)
func (p *Parser) GetBool(pattern string) bool {
    return p.json.GetBool(pattern)
}

func (p *Parser) GetInt(pattern string) int {
    return p.json.GetInt(pattern)
}

func (p *Parser) GetUint(pattern string) uint {
    return p.json.GetUint(pattern)
}

func (p *Parser) GetFloat32(pattern string) float32 {
    return p.json.GetFloat32(pattern)
}

func (p *Parser) GetFloat64(pattern string) float64 {
    return p.json.GetFloat64(pattern)
}

// 根据pattern查找并设置数据
// 注意：写入的时候"."符号只能表示层级，不能使用带"."符号的键名
func (p *Parser) Set(pattern string, value interface{}) error {
    return p.json.Set(pattern, value)
}

// 动态删除变量节点
func (p *Parser) Remove(pattern string) error {
    return p.json.Remove(pattern)
}

// 根据约定字符串方式访问json解析数据，参数形如： "items.name.first", "list.0"
// 返回的结果类型的interface{}，因此需要自己做类型转换
// 如果找不到对应节点的数据，返回nil
func (p *Parser) Get(pattern string) interface{} {
    return p.json.Get(pattern)
}

// 转换为map[string]interface{}类型,如果转换失败，返回nil
func (p *Parser) ToMap() map[string]interface{} {
    return p.json.ToMap()
}

// 转换为[]interface{}类型,如果转换失败，返回nil
func (p *Parser) ToArray() []interface{} {
    return p.json.ToArray()
}

/* 以下为数据文件格式转换，支持类型：xml, json, yaml/yml, toml */

func (p *Parser) ToXml(rootTag...string) ([]byte, error) {
    return p.json.ToXml(rootTag...)
}

func (p *Parser) ToXmlIndent(rootTag...string) ([]byte, error) {
    return p.json.ToXmlIndent(rootTag...)
}

func (p *Parser) ToJson() ([]byte, error) {
    return p.json.ToJson()
}

func (p *Parser) ToJsonIndent() ([]byte, error) {
    return p.json.ToJsonIndent()
}

func (p *Parser) ToYaml() ([]byte, error) {
    return p.json.ToYaml()
}

func (p *Parser) ToToml() ([]byte, error) {
    return p.json.ToToml()
}

func VarToXml(value interface{}, rootTag...string) ([]byte, error) {
    return gjson.NewJson(value).ToXml(rootTag...)
}

func VarToXmlIndent(value interface{}, rootTag...string) ([]byte, error) {
    return gjson.NewJson(value).ToXmlIndent(rootTag...)
}

func VarToJson(value interface{}) ([]byte, error) {
    return gjson.NewJson(value).ToJson()
}

func VarToJsonIndent(value interface{}) ([]byte, error) {
    return gjson.NewJson(value).ToJsonIndent()
}

func VarToYaml(value interface{}) ([]byte, error) {
    return gjson.NewJson(value).ToYaml()
}

func VarToToml(value interface{}) ([]byte, error) {
    return gjson.NewJson(value).ToToml()
}