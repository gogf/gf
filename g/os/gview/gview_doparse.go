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
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/os/gfcache"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gfsnotify"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gspath"
	"github.com/gogf/gf/g/text/gstr"
	"text/template"
)

var (
	// Templates cache map for template folder.
	templates = gmap.NewStrAnyMap()
)

// getTemplate returns the template object associated with given template folder <path>.
// It uses template cache to enhance performance, that is, it will return the same template object
// with the same given <path>. It will also refresh the template cache
// if the template files under <path> changes (recursively).
func (view *View) getTemplate(path string, pattern string) (tpl *template.Template, err error) {
	r := templates.GetOrSetFuncLock(path, func() interface {} {
		files     := ([]string)(nil)
		files, err = gfile.ScanDir(path, pattern, true)
		if err != nil {
			return nil
		}
		tpl = template.New(path).Delims(view.delimiters[0], view.delimiters[1]).Funcs(view.funcMap)
		if tpl, err = tpl.ParseFiles(files...); err != nil {
			return nil
		}
		gfsnotify.Add(path, func(event *gfsnotify.Event) {
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
	view.paths.RLockFunc(func(array []string) {
		for _, v := range array {
			if path, _ = gspath.Search(v, file); path != "" {
				folder = v
				break
			}
			if path, _ = gspath.Search(v + gfile.Separator + "template", file); path != "" {
				folder = v + gfile.Separator + "template"
				break
			}
		}
	})
	if path == "" {
		buffer := bytes.NewBuffer(nil)
		if view.paths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("[gview] cannot find template file \"%s\" in following paths:", file))
			view.paths.RLockFunc(func(array []string) {
				index := 1
				for _, v := range array {
					buffer.WriteString(fmt.Sprintf("\n%d. %s", index,  v))
					index++
					buffer.WriteString(fmt.Sprintf("\n%d. %s", index,  v + gfile.Separator + "template"))
					index++
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf("[gview] cannot find template file \"%s\" with no path set/add", file))
		}
		glog.Error(buffer.String())
		err = errors.New(fmt.Sprintf(`template file "%s" not found`, file))
	}
	return
}

// ParseContent parses given template file <file>
// with given template parameters <params> and function map <funcMap>
// and returns the parsed string content.
func (view *View) Parse(file string, params...Params) (parsed string, err error) {
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
	tpl, err = tpl.Parse(gfcache.GetContents(path))
	if err != nil {
		return "", err
	}
	// Note that the template variable assignment cannot change the value
	// of the existing <params> or view.data because both variables are pointers.
	// It's need to merge the values of the two maps into a new map.
	vars := (map[string]interface{})(nil)
	length := len(view.data)
	if len(params) > 0 {
		length += len(params[0])
	}
	if length > 0 {
		vars = make(map[string]interface{}, length)
	}
	if len(view.data) > 0 {
		if len(params) > 0 {
			for k, v := range params[0] {
				vars[k] = v
			}
			for k, v := range view.data {
				vars[k] = v
			}
		} else {
			vars = view.data
		}
	} else {
		if len(params) > 0 {
			vars = params[0]
		}
	}
	buffer := bytes.NewBuffer(nil)
	if err := tpl.Execute(buffer, vars); err != nil {
		return "", err
	}
	return gstr.Replace(buffer.String(), "<no value>", ""), nil
}

// ParseContent parses given template content <content>
// with given template parameters <params> and function map <funcMap>
// and returns the parsed content in []byte.
func (view *View) ParseContent(content string, params...Params) (string, error) {
	view.mu.RLock()
	defer view.mu.RUnlock()
	tpl := template.New("template content").Delims(view.delimiters[0], view.delimiters[1]).Funcs(view.funcMap)
	tpl, err := tpl.Parse(content)
	if err != nil {
		return "", err
	}
	// Note that the template variable assignment cannot change the value
	// of the existing <params> or view.data because both variables are pointers.
	// It's need to merge the values of the two maps into a new map.
	vars := (map[string]interface{})(nil)
	length := len(view.data)
	if len(params) > 0 {
		length += len(params[0])
	}
	if length > 0 {
		vars = make(map[string]interface{}, length)
	}
	if len(view.data) > 0 {
		if len(params) > 0 {
			for k, v := range params[0] {
				vars[k] = v
			}
			for k, v := range view.data {
				vars[k] = v
			}
		} else {
			vars = view.data
		}
	} else {
		if len(params) > 0 {
			vars = params[0]
		}
	}
	buffer := bytes.NewBuffer(nil)
	if err := tpl.Execute(buffer, vars); err != nil {
		return "", err
	}
	return gstr.Replace(buffer.String(), "<no value>", ""), nil
}
