// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/gogf/gf/os/gfcache"

	"github.com/gogf/gf/os/gres"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/encoding/ghash"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gfsnotify"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gmlock"
	"github.com/gogf/gf/os/gspath"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

const (
	// Template name for content parsing.
	gCONTENT_TEMPLATE_NAME = "template content"
)

var (
	// Templates cache map for template folder.
	// TODO Note that there's no expiring logic for this map.
	templates        = gmap.NewStrAnyMap(true)
	resourceTryFiles = []string{"template/", "template", "/template", "/template/"}
)

// getTemplate returns the template object associated with given template folder <path>.
// It uses template cache to enhance performance, that is, it will return the same template object
// with the same given <path>. It will also refresh the template cache
// if the template files under <path> changes (recursively).
func (view *View) getTemplate(path string, pattern string) (tpl *template.Template, err error) {
	r := templates.GetOrSetFuncLock(path, func() interface{} {
		tpl = template.New(path).Delims(view.delimiters[0], view.delimiters[1]).Funcs(view.funcMap)
		// Firstly checking the resource manager.
		if files := gres.ScanDirFile(path, pattern, true); len(files) > 0 {
			var err error
			for _, v := range files {
				_, err = tpl.New(v.FileInfo().Name()).Parse(string(v.Content()))
				if err != nil {
					glog.Error(err)
				}
			}
			return tpl
		}
		// Secondly checking the file system.
		files := ([]string)(nil)
		files, err = gfile.ScanDir(path, pattern, true)
		if err != nil {
			return nil
		}
		if tpl, err = tpl.ParseFiles(files...); err != nil {
			return nil
		}
		_, _ = gfsnotify.Add(path, func(event *gfsnotify.Event) {
			templates.Remove(path)
			gfsnotify.Exit()
		})
		return tpl
	})
	if r != nil {
		return r.(*template.Template), nil
	}
	return
}

// searchFile returns the found absolute path for <file>, and its template folder path.
func (view *View) searchFile(file string) (path string, folder string, err error) {
	// Firstly checking the resource manager.
	view.paths.RLockFunc(func(array []string) {
		f := (*gres.File)(nil)
		for _, v := range array {
			v = strings.TrimRight(v, "/")
			if f = gres.Get(v + "/" + file); f != nil {
				path = f.Name()
				folder = gfile.Dir(path)
				break
			}
			if f = gres.Get(v + "/template/" + file); f != nil {
				path = f.Name()
				folder = gfile.Dir(path)
				break
			}
		}
	})
	// Secondly checking the file system.
	if path == "" {
		view.paths.RLockFunc(func(array []string) {
			for _, v := range array {
				v = strings.TrimRight(v, gfile.Separator)
				if path, _ = gspath.Search(v, file); path != "" {
					folder = v
					break
				}
				if path, _ = gspath.Search(v+gfile.Separator+"template", file); path != "" {
					folder = v + gfile.Separator + "template"
					break
				}
			}
		})
	}
	// Checking the configuration file in default paths.
	if path == "" && !gres.IsEmpty() {
		for _, v := range resourceTryFiles {
			if file := gres.Get(v + file); file != nil {
				path = file.Name()
				folder = gfile.Dir(path)
				return
			}
		}
	}
	// Error checking.
	if path == "" {
		buffer := bytes.NewBuffer(nil)
		if view.paths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("[gview] cannot find template file \"%s\" in following paths:", file))
			view.paths.RLockFunc(func(array []string) {
				index := 1
				for _, v := range array {
					v = strings.TrimRight(v, "/")
					if v == "" {
						v = "/"
					}
					buffer.WriteString(fmt.Sprintf("\n%d. %s", index, v))
					index++
					buffer.WriteString(fmt.Sprintf("\n%d. %s", index, strings.TrimRight(v, "/")+gfile.Separator+"template"))
					index++
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf("[gview] cannot find template file \"%s\" with no path set/add", file))
		}
		if errorPrint() {
			glog.Error(buffer.String())
		}
		err = errors.New(fmt.Sprintf(`template file "%s" not found`, file))
	}
	return
}

// ParseContent parses given template file <file>
// with given template parameters <params> and function map <funcMap>
// and returns the parsed string content.
func (view *View) Parse(file string, params ...Params) (parsed string, err error) {
	view.mu.RLock()
	defer view.mu.RUnlock()
	path, folder, err := view.searchFile(file)
	if err != nil {
		return "", err
	}
	tpl, err := view.getTemplate(folder, fmt.Sprintf(`*%s`, gfile.Ext(path)))
	if err != nil {
		return "", err
	}
	// Using memory lock to ensure concurrent safety for template parsing.
	gmlock.LockFunc("gview-parsing:"+folder, func() {
		if file := gres.Get(path); file != nil {
			tpl, err = tpl.Parse(string(file.Content()))
		} else {
			tpl, err = tpl.Parse(gfcache.GetContents(path))
		}
	})
	if err != nil {
		return "", err
	}
	// Note that the template variable assignment cannot change the value
	// of the existing <params> or view.data because both variables are pointers.
	// It's need to merge the values of the two maps into a new map.
	variables := (map[string]interface{})(nil)
	length := len(view.data)
	if len(params) > 0 {
		length += len(params[0])
	}
	if length > 0 {
		variables = make(map[string]interface{}, length)
	}
	if len(view.data) > 0 {
		if len(params) > 0 {
			for k, v := range params[0] {
				variables[k] = v
			}
			for k, v := range view.data {
				variables[k] = v
			}
		} else {
			variables = view.data
		}
	} else {
		if len(params) > 0 {
			variables = params[0]
		}
	}
	buffer := bytes.NewBuffer(nil)
	if err := tpl.Execute(buffer, variables); err != nil {
		return "", err
	}
	result := gstr.Replace(buffer.String(), "<no value>", "")
	result = view.i18nTranslate(result, variables)
	return result, nil
}

// ParseContent parses given template content <content>
// with given template parameters <params> and function map <funcMap>
// and returns the parsed content in []byte.
func (view *View) ParseContent(content string, params ...Params) (string, error) {
	view.mu.RLock()
	defer view.mu.RUnlock()
	err := (error)(nil)
	tpl := templates.GetOrSetFuncLock(gCONTENT_TEMPLATE_NAME, func() interface{} {
		return template.New(gCONTENT_TEMPLATE_NAME).Delims(view.delimiters[0], view.delimiters[1]).Funcs(view.funcMap)
	}).(*template.Template)
	// Using memory lock to ensure concurrent safety for content parsing.
	hash := gconv.String(ghash.DJBHash64([]byte(content)))
	gmlock.LockFunc("gview-parsing-content:"+hash, func() {
		tpl, err = tpl.Parse(content)
	})
	if err != nil {
		return "", err
	}
	// Note that the template variable assignment cannot change the value
	// of the existing <params> or view.data because both variables are pointers.
	// It's need to merge the values of the two maps into a new map.
	variables := (map[string]interface{})(nil)
	length := len(view.data)
	if len(params) > 0 {
		length += len(params[0])
	}
	if length > 0 {
		variables = make(map[string]interface{}, length)
	}
	if len(view.data) > 0 {
		if len(params) > 0 {
			for k, v := range params[0] {
				variables[k] = v
			}
			for k, v := range view.data {
				variables[k] = v
			}
		} else {
			variables = view.data
		}
	} else {
		if len(params) > 0 {
			variables = params[0]
		}
	}
	buffer := bytes.NewBuffer(nil)
	if err := tpl.Execute(buffer, variables); err != nil {
		return "", err
	}
	result := gstr.Replace(buffer.String(), "<no value>", "")
	result = view.i18nTranslate(result, variables)
	return result, nil
}
