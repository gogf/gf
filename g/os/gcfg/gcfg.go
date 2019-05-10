// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcfg provides reading, caching and managing for configuration.
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
	"github.com/gogf/gf/g/internal/cmdenv"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gfsnotify"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gspath"
	"github.com/gogf/gf/g/os/gtime"
	"time"
)

const (
    // Default configuration file name.
    DEFAULT_CONFIG_FILE = "config.toml"
)

// Configuration struct.
type Config struct {
    name   *gtype.String            // Default configuration file name.
    paths  *garray.StringArray      // Searching path array.
    jsons  *gmap.StrAnyMap          // The pared JSON objects for configuration files.
    vc     *gtype.Bool              // Whether do violence check in value index searching.
                                    // It affects the performance when set true(false in default).
}

// New returns a new configuration management object.
// The param <file> specifies the default configuration file name for reading.
func New(file...string) *Config {
    name := DEFAULT_CONFIG_FILE
    if len(file) > 0 {
        name = file[0]
    }
    c := &Config {
        name   : gtype.NewString(name),
        paths  : garray.NewStringArray(),
        jsons  : gmap.NewStrAnyMap(),
        vc     : gtype.NewBool(),
    }
	// Customized dir path from env/cmd.
	if envPath := cmdenv.Get("gf.gcfg.path").String(); envPath != "" {
		if gfile.Exists(envPath) {
			c.SetPath(envPath)
		} else {
			glog.Errorfln("Configuration directory path does not exist: %s", envPath)
		}
	} else {
		// Dir path of working dir.
		c.SetPath(gfile.Pwd())
		// Dir path of binary.
		if selfPath := gfile.SelfDir(); selfPath != "" && gfile.Exists(selfPath) {
			c.AddPath(selfPath)
		}
		// Dir path of main package.
		if mainPath := gfile.MainPkgPath(); mainPath != "" && gfile.Exists(mainPath) {
			c.AddPath(mainPath)
		}
	}
    return c
}

// filePath returns the absolute configuration file path for the given filename by <file>.
func (c *Config) filePath(file...string) (path string) {
    name := c.name.Val()
    if len(file) > 0 {
        name = file[0]
    }
    path = c.FilePath(name)
    if path == "" {
        buffer := bytes.NewBuffer(nil)
        if c.paths.Len() > 0 {
            buffer.WriteString(fmt.Sprintf("[gcfg] cannot find config file \"%s\" in following paths:", name))
            c.paths.RLockFunc(func(array []string) {
            	index := 1
                for _, v := range array {
                    buffer.WriteString(fmt.Sprintf("\n%d. %s", index,  v))
                    index++
                    buffer.WriteString(fmt.Sprintf("\n%d. %s", index,  v + gfile.Separator + "config"))
	                index++
                }
            })
        } else {
            buffer.WriteString(fmt.Sprintf("[gcfg] cannot find config file \"%s\" with no path set/add", name))
        }
        glog.Error(buffer.String())
    }
    return path
}

// SetPath sets the configuration directory path for file search.
// The param <path> can be absolute or relative path,
// but absolute path is strongly recommended.
func (c *Config) SetPath(path string) error {
    // Absolute path.
    realPath := gfile.RealPath(path)
    if realPath == "" {
        // Relative path.
        c.paths.RLockFunc(func(array []string) {
            for _, v := range array {
                if path, _ := gspath.Search(v, path); path != "" {
                    realPath = path
                    break
                }
            }
        })
    }
    // Path not exist.
    if realPath == "" {
        buffer := bytes.NewBuffer(nil)
        if c.paths.Len() > 0 {
            buffer.WriteString(fmt.Sprintf("[gcfg] SetPath failed: cannot find directory \"%s\" in following paths:", path))
            c.paths.RLockFunc(func(array []string) {
                for k, v := range array {
                    buffer.WriteString(fmt.Sprintf("\n%d. %s",k + 1,  v))
                }
            })
        } else {
            buffer.WriteString(fmt.Sprintf(`[gcfg] SetPath failed: path "%s" does not exist`, path))
        }
        err := errors.New(buffer.String())
        glog.Error(err)
        return err
    }
    // Should be a directory.
    if !gfile.IsDir(realPath) {
        err := errors.New(fmt.Sprintf(`[gcfg] SetPath failed: path "%s" should be directory type`, path))
        glog.Error(err)
        return err
    }
    // Repeated path check.
    if c.paths.Search(realPath) != -1 {
        return nil
    }
    c.jsons.Clear()
    c.paths.Clear()
    c.paths.Append(realPath)
    return nil
}

// SetViolenceCheck sets whether to perform hierarchical conflict check.
// This feature needs to be enabled when there is a level symbol in the key name.
// The default is off.
// Turning on this feature is quite expensive,
// and it is not recommended to allow separators in the key names.
// It is best to avoid this on the application side.
func (c *Config) SetViolenceCheck(check bool) {
    c.vc.Set(check)
    c.Clear()
}

// AddPath adds a absolute or relative path to the search paths.
func (c *Config) AddPath(path string) error {
    // Absolute path.
    realPath := gfile.RealPath(path)
    if realPath == "" {
	    // Relative path.
        c.paths.RLockFunc(func(array []string) {
            for _, v := range array {
                if path, _ := gspath.Search(v, path); path != "" {
                    realPath = path
                    break
                }
            }
        })
    }
    if realPath == "" {
        buffer := bytes.NewBuffer(nil)
        if c.paths.Len() > 0 {
            buffer.WriteString(fmt.Sprintf("[gcfg] AddPath failed: cannot find directory \"%s\" in following paths:", path))
            c.paths.RLockFunc(func(array []string) {
                for k, v := range array {
                    buffer.WriteString(fmt.Sprintf("\n%d. %s", k + 1,  v))
                }
            })
        } else {
            buffer.WriteString(fmt.Sprintf(`[gcfg] AddPath failed: path "%s" does not exist`, path))
        }
        err := errors.New(buffer.String())
        glog.Error(err)
        return err
    }
    if !gfile.IsDir(realPath) {
        err := errors.New(fmt.Sprintf(`[gcfg] AddPath failed: path "%s" should be directory type`, path))
        glog.Error(err)
        return err
    }
    // Repeated path check.
    if c.paths.Search(realPath) != -1 {
        return nil
    }
    c.paths.Append(realPath)
    //glog.Debug("[gcfg] AddPath:", realPath)
    return nil
}

// Deprecated.
// Alias of FilePath.
func (c *Config) GetFilePath(file...string) (path string) {
	return c.FilePath(file...)
}

// GetFilePath returns the absolute path of the specified configuration file.
// If <file> is not passed, it returns the configuration file path of the default name.
// If the specified configuration file does not exist,
// an empty string is returned.
func (c *Config) FilePath(file...string) (path string) {
    name := c.name.Val()
    if len(file) > 0 {
        name = file[0]
    }
    c.paths.RLockFunc(func(array []string) {
        for _, v := range array {
            if path, _ = gspath.Search(v, name); path != "" {
                break
            }
            if path, _ = gspath.Search(v + gfile.Separator + "config", name); path != "" {
                break
            }
        }
    })
    return
}

// SetFileName sets the default configuration file name.
func (c *Config) SetFileName(name string) {
    c.name.Set(name)
}

// GetFileName returns the default configuration file name.
func (c *Config) GetFileName() string {
    return c.name.Val()
}

// getJson returns a gjson.Json object for the specified <file> content.
// It would print error if file reading fails.
// If any error occurs, it return nil.
func (c *Config) getJson(file...string) *gjson.Json {
    name := c.name.Val()
    if len(file) > 0 {
        name = file[0]
    }
    r := c.jsons.GetOrSetFuncLock(name, func() interface{} {
        content  := ""
        filePath := ""
        if content = GetContent(name); content == "" {
            filePath = c.filePath(name)
            if filePath == "" {
                return nil
            }
            content = gfile.GetContents(filePath)
        }
        if j, err := gjson.LoadContent(content); err == nil {
            j.SetViolenceCheck(c.vc.Val())
            // Add monitor for this configuration file,
            // any changes of this file will refresh its cache in Config object.
            if filePath != "" {
                gfsnotify.Add(filePath, func(event *gfsnotify.Event) {
                    c.jsons.Remove(name)
                })
            }
            return j
        } else {
            if filePath != "" {
                glog.Criticalfln(`[gcfg] Load config file "%s" failed: %s`, filePath, err.Error())
            } else {
                glog.Criticalfln(`[gcfg] Load configuration failed: %s`, err.Error())
            }
        }
        return nil
    })
    if r != nil {
        return r.(*gjson.Json)
    }
    return nil
}

func (c *Config) Get(pattern string, def...interface{}) interface{} {
    if j := c.getJson(); j != nil {
        return j.Get(pattern, def...)
    }
    return nil
}

func (c *Config) GetVar(pattern string, def...interface{}) *gvar.Var {
    if j := c.getJson(); j != nil {
        return gvar.New(j.Get(pattern, def...), true)
    }
    return gvar.New(nil, true)
}

func (c *Config) Contains(pattern string) bool {
    if j := c.getJson(); j != nil {
        return j.Contains(pattern)
    }
    return false
}

func (c *Config) GetMap(pattern string, def...interface{}) map[string]interface{} {
    if j := c.getJson(); j != nil {
        return j.GetMap(pattern, def...)
    }
    return nil
}

func (c *Config) GetArray(pattern string, def...interface{}) []interface{} {
    if j := c.getJson(); j != nil {
        return j.GetArray(pattern, def...)
    }
    return nil
}

func (c *Config) GetString(pattern string, def...interface{}) string {
    if j := c.getJson(); j != nil {
        return j.GetString(pattern, def...)
    }
    return ""
}

func (c *Config) GetStrings(pattern string, def...interface{}) []string {
    if j := c.getJson(); j != nil {
        return j.GetStrings(pattern, def...)
    }
    return nil
}

func (c *Config) GetInterfaces(pattern string, def...interface{}) []interface{} {
    if j := c.getJson(); j != nil {
        return j.GetInterfaces(pattern, def...)
    }
    return nil
}

func (c *Config) GetBool(pattern string, def...interface{}) bool {
    if j := c.getJson(); j != nil {
        return j.GetBool(pattern, def...)
    }
    return false
}

func (c *Config) GetFloat32(pattern string, def...interface{}) float32 {
    if j := c.getJson(); j != nil {
        return j.GetFloat32(pattern, def...)
    }
    return 0
}

func (c *Config) GetFloat64(pattern string, def...interface{}) float64 {
    if j := c.getJson(); j != nil {
        return j.GetFloat64(pattern, def...)
    }
    return 0
}

func (c *Config) GetFloats(pattern string, def...interface{}) []float64 {
    if j := c.getJson(); j != nil {
        return j.GetFloats(pattern, def...)
    }
    return nil
}

func (c *Config) GetInt(pattern string, def...interface{}) int {
    if j := c.getJson(); j != nil {
        return j.GetInt(pattern, def...)
    }
    return 0
}


func (c *Config) GetInt8(pattern string, def...interface{}) int8 {
    if j := c.getJson(); j != nil {
        return j.GetInt8(pattern, def...)
    }
    return 0
}

func (c *Config) GetInt16(pattern string, def...interface{}) int16 {
    if j := c.getJson(); j != nil {
        return j.GetInt16(pattern, def...)
    }
    return 0
}

func (c *Config) GetInt32(pattern string, def...interface{}) int32 {
    if j := c.getJson(); j != nil {
        return j.GetInt32(pattern, def...)
    }
    return 0
}

func (c *Config) GetInt64(pattern string, def...interface{}) int64 {
    if j := c.getJson(); j != nil {
        return j.GetInt64(pattern, def...)
    }
    return 0
}

func (c *Config) GetInts(pattern string, def...interface{}) []int {
    if j := c.getJson(); j != nil {
        return j.GetInts(pattern, def...)
    }
    return nil
}

func (c *Config) GetUint(pattern string, def...interface{}) uint {
    if j := c.getJson(); j != nil {
        return j.GetUint(pattern, def...)
    }
    return 0
}

func (c *Config) GetUint8(pattern string, def...interface{}) uint8 {
    if j := c.getJson(); j != nil {
        return j.GetUint8(pattern, def...)
    }
    return 0
}

func (c *Config) GetUint16(pattern string, def...interface{}) uint16 {
    if j := c.getJson(); j != nil {
        return j.GetUint16(pattern, def...)
    }
    return 0
}

func (c *Config) GetUint32(pattern string, def...interface{}) uint32 {
    if j := c.getJson(); j != nil {
        return j.GetUint32(pattern, def...)
    }
    return 0
}

func (c *Config) GetUint64(pattern string, def...interface{}) uint64 {
    if j := c.getJson(); j != nil {
        return j.GetUint64(pattern, def...)
    }
    return 0
}

func (c *Config) GetTime(pattern string, format...string) time.Time {
	if j := c.getJson(); j != nil {
		return j.GetTime(pattern, format...)
	}
	return time.Time{}
}

func (c *Config) GetDuration(pattern string, def...interface{}) time.Duration {
	if j := c.getJson(); j != nil {
		return j.GetDuration(pattern, def...)
	}
	return 0
}

func (c *Config) GetGTime(pattern string, format...string) *gtime.Time {
	if j := c.getJson(); j != nil {
		return j.GetGTime(pattern, format...)
	}
	return nil
}

func (c *Config) GetToStruct(pattern string, pointer interface{}, def...interface{}) error {
    if j := c.getJson(); j != nil {
        return j.GetToStruct(pattern, pointer)
    }
    return errors.New("config file not found")
}

// Deprecated. See Clear.
func (c *Config) Reload() {
    c.jsons.Clear()
}

// Clear removes all parsed configuration files content cache,
// which will force reload configuration content from file.
func (c *Config) Clear() {
    c.jsons.Clear()
}

