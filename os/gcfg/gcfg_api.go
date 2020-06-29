// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import (
	"errors"
	"time"

	"github.com/gogf/gf/encoding/gjson"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/os/gtime"
)

func (c *Config) Get(pattern string, def ...interface{}) interface{} {
	if j := c.getJson(); j != nil {
		return j.Get(pattern, def...)
	}
	return nil
}

func (c *Config) GetVar(pattern string, def ...interface{}) *gvar.Var {
	if j := c.getJson(); j != nil {
		return gvar.New(j.Get(pattern, def...))
	}
	return gvar.New(nil)
}

func (c *Config) Contains(pattern string) bool {
	if j := c.getJson(); j != nil {
		return j.Contains(pattern)
	}
	return false
}

func (c *Config) GetMap(pattern string, def ...interface{}) map[string]interface{} {
	if j := c.getJson(); j != nil {
		return j.GetMap(pattern, def...)
	}
	return nil
}

func (c *Config) GetMapStrStr(pattern string, def ...interface{}) map[string]string {
	if j := c.getJson(); j != nil {
		return j.GetMapStrStr(pattern, def...)
	}
	return nil
}

func (c *Config) GetArray(pattern string, def ...interface{}) []interface{} {
	if j := c.getJson(); j != nil {
		return j.GetArray(pattern, def...)
	}
	return nil
}

func (c *Config) GetBytes(pattern string, def ...interface{}) []byte {
	if j := c.getJson(); j != nil {
		return j.GetBytes(pattern, def...)
	}
	return nil
}

func (c *Config) GetString(pattern string, def ...interface{}) string {
	if j := c.getJson(); j != nil {
		return j.GetString(pattern, def...)
	}
	return ""
}

func (c *Config) GetStrings(pattern string, def ...interface{}) []string {
	if j := c.getJson(); j != nil {
		return j.GetStrings(pattern, def...)
	}
	return nil
}

func (c *Config) GetInterfaces(pattern string, def ...interface{}) []interface{} {
	if j := c.getJson(); j != nil {
		return j.GetInterfaces(pattern, def...)
	}
	return nil
}

func (c *Config) GetBool(pattern string, def ...interface{}) bool {
	if j := c.getJson(); j != nil {
		return j.GetBool(pattern, def...)
	}
	return false
}

func (c *Config) GetFloat32(pattern string, def ...interface{}) float32 {
	if j := c.getJson(); j != nil {
		return j.GetFloat32(pattern, def...)
	}
	return 0
}

func (c *Config) GetFloat64(pattern string, def ...interface{}) float64 {
	if j := c.getJson(); j != nil {
		return j.GetFloat64(pattern, def...)
	}
	return 0
}

func (c *Config) GetFloats(pattern string, def ...interface{}) []float64 {
	if j := c.getJson(); j != nil {
		return j.GetFloats(pattern, def...)
	}
	return nil
}

func (c *Config) GetInt(pattern string, def ...interface{}) int {
	if j := c.getJson(); j != nil {
		return j.GetInt(pattern, def...)
	}
	return 0
}

func (c *Config) GetInt8(pattern string, def ...interface{}) int8 {
	if j := c.getJson(); j != nil {
		return j.GetInt8(pattern, def...)
	}
	return 0
}

func (c *Config) GetInt16(pattern string, def ...interface{}) int16 {
	if j := c.getJson(); j != nil {
		return j.GetInt16(pattern, def...)
	}
	return 0
}

func (c *Config) GetInt32(pattern string, def ...interface{}) int32 {
	if j := c.getJson(); j != nil {
		return j.GetInt32(pattern, def...)
	}
	return 0
}

func (c *Config) GetInt64(pattern string, def ...interface{}) int64 {
	if j := c.getJson(); j != nil {
		return j.GetInt64(pattern, def...)
	}
	return 0
}

func (c *Config) GetInts(pattern string, def ...interface{}) []int {
	if j := c.getJson(); j != nil {
		return j.GetInts(pattern, def...)
	}
	return nil
}

func (c *Config) GetUint(pattern string, def ...interface{}) uint {
	if j := c.getJson(); j != nil {
		return j.GetUint(pattern, def...)
	}
	return 0
}

func (c *Config) GetUint8(pattern string, def ...interface{}) uint8 {
	if j := c.getJson(); j != nil {
		return j.GetUint8(pattern, def...)
	}
	return 0
}

func (c *Config) GetUint16(pattern string, def ...interface{}) uint16 {
	if j := c.getJson(); j != nil {
		return j.GetUint16(pattern, def...)
	}
	return 0
}

func (c *Config) GetUint32(pattern string, def ...interface{}) uint32 {
	if j := c.getJson(); j != nil {
		return j.GetUint32(pattern, def...)
	}
	return 0
}

func (c *Config) GetUint64(pattern string, def ...interface{}) uint64 {
	if j := c.getJson(); j != nil {
		return j.GetUint64(pattern, def...)
	}
	return 0
}

func (c *Config) GetTime(pattern string, format ...string) time.Time {
	if j := c.getJson(); j != nil {
		return j.GetTime(pattern, format...)
	}
	return time.Time{}
}

func (c *Config) GetDuration(pattern string, def ...interface{}) time.Duration {
	if j := c.getJson(); j != nil {
		return j.GetDuration(pattern, def...)
	}
	return 0
}

func (c *Config) GetGTime(pattern string, format ...string) *gtime.Time {
	if j := c.getJson(); j != nil {
		return j.GetGTime(pattern, format...)
	}
	return nil
}

func (c *Config) GetJson(pattern string, def ...interface{}) *gjson.Json {
	if j := c.getJson(); j != nil {
		return j.GetJson(pattern, def...)
	}
	return nil
}

func (c *Config) GetJsons(pattern string, def ...interface{}) []*gjson.Json {
	if j := c.getJson(); j != nil {
		return j.GetJsons(pattern, def...)
	}
	return nil
}

func (c *Config) GetJsonMap(pattern string, def ...interface{}) map[string]*gjson.Json {
	if j := c.getJson(); j != nil {
		return j.GetJsonMap(pattern, def...)
	}
	return nil
}

func (c *Config) GetStruct(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.GetStruct(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) GetStructDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.GetStructDeep(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) GetStructs(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.GetStructs(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) GetStructsDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.GetStructsDeep(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) GetMapToMap(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.GetMapToMap(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) GetMapToMapDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.GetMapToMapDeep(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) GetMapToMaps(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.GetMapToMaps(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) GetMapToMapsDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.GetMapToMapsDeep(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) ToMap() map[string]interface{} {
	if j := c.getJson(); j != nil {
		return j.ToMap()
	}
	return nil
}

func (c *Config) ToArray() []interface{} {
	if j := c.getJson(); j != nil {
		return j.ToArray()
	}
	return nil
}

func (c *Config) ToStruct(pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.ToStruct(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) ToStructDeep(pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.ToStructDeep(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) ToStructs(pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.ToStructs(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) ToStructsDeep(pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.ToStructsDeep(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) ToMapToMap(pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.ToMapToMap(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) ToMapToMapDeep(pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.ToMapToMapDeep(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) ToMapToMaps(pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.ToMapToMaps(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

func (c *Config) ToMapToMapsDeep(pointer interface{}, mapping ...map[string]string) error {
	if j := c.getJson(); j != nil {
		return j.ToMapToMapsDeep(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// Clear removes all parsed configuration files content cache,
// which will force reload configuration content from file.
func (c *Config) Clear() {
	c.jsons.Clear()
}

// Dump prints current Json object with more manually readable.
func (c *Config) Dump() {
	if j := c.getJson(); j != nil {
		j.Dump()
	}
}
