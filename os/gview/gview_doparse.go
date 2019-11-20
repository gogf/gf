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
	"github.com/gogf/gf/encoding/ghash"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gfcache"
	"github.com/gogf/gf/os/gfsnotify"
	"github.com/gogf/gf/os/gmlock"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"strconv"
	"strings"
	"text/template"

	"github.com/gogf/gf/os/gres"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gspath"
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

// fileCacheItem is the cache item for template file.
type fileCacheItem struct {
	path    string
	folder  string
	content string
}

// Parse parses given template file <file> with given template variables <params>
// and returns the parsed template content.
func (view *View) Parse(file string, params ...Params) (result string, err error) {
	var tpl *template.Template
	// It caches the file, folder and its content to enhance performance.
	r := view.fileCacheMap.GetOrSetFuncLock(file, func() interface{} {
		var path, folder, content string
		var resource *gres.File
		// Searching the absolute file path for <file>.
		path, folder, resource, err = view.searchFile(file)
		if err != nil {
			return nil
		}
		if resource != nil {
			content = gconv.UnsafeBytesToStr(resource.Content())
		} else {
			content = gfcache.GetContents(path)
		}
		// Monitor template files changes using fsnotify asynchronously.
		if resource == nil {
			if _, err := gfsnotify.AddOnce("gview.Parse:"+folder, folder, func(event *gfsnotify.Event) {
				// CLEAR THEM ALL.
				view.fileCacheMap.Clear()
				templates.Clear()
				gfsnotify.Exit()
			}); err != nil {
				intlog.Error(err)
			}
		}
		return &fileCacheItem{
			path:    path,
			folder:  folder,
			content: content,
		}
	})
	if r == nil {
		return
	}
	item := r.(*fileCacheItem)
	// Get the template object instance for <folder>.
	tpl, err = view.getTemplate(item.folder, fmt.Sprintf(`*%s`, gfile.Ext(item.path)))
	if err != nil {
		return "", err
	}
	// Using memory lock to ensure concurrent safety for template parsing.
	gmlock.LockFunc("gview.Parse:"+item.folder, func() {
		tpl, err = tpl.Parse(item.content)
	})

	// Note that the template variable assignment cannot change the value
	// of the existing <params> or view.data because both variables are pointers.
	// It needs to merge the values of the two maps into a new map.
	var variables map[string]interface{}
	length := len(view.data)
	if len(params) > 0 {
		length += len(params[0])
	}
	if length > 0 {
		variables = make(map[string]interface{}, length)
	}
	if len(view.data) > 0 {
		if len(params) > 0 {
			if variables == nil {
				variables = make(map[string]interface{})
			}
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
	// TODO any graceful plan to replace "<no value>"?
	result = gstr.Replace(buffer.String(), "<no value>", "")
	result = view.i18nTranslate(result, variables)
	return result, nil
}

// ParseDefault parses the default template file with params.
func (view *View) ParseDefault(params ...Params) (result string, err error) {
	return view.Parse(view.defaultFile, params...)
}

// ParseContent parses given template content <content>  with template variables <params>
// and returns the parsed content in []byte.
func (view *View) ParseContent(content string, params ...Params) (string, error) {
	err := (error)(nil)
	tpl := templates.GetOrSetFuncLock(gCONTENT_TEMPLATE_NAME, func() interface{} {
		return template.New(gCONTENT_TEMPLATE_NAME).Delims(view.delimiters[0], view.delimiters[1]).Funcs(view.funcMap)
	}).(*template.Template)
	// Using memory lock to ensure concurrent safety for content parsing.
	hash := strconv.FormatUint(ghash.DJBHash64([]byte(content)), 10)
	gmlock.LockFunc("gview.ParseContent:"+hash, func() {
		tpl, err = tpl.Parse(content)
	})
	if err != nil {
		return "", err
	}
	// Note that the template variable assignment cannot change the value
	// of the existing <params> or view.data because both variables are pointers.
	// It needs to merge the values of the two maps into a new map.
	var variables map[string]interface{}
	length := len(view.data)
	if len(params) > 0 {
		length += len(params[0])
	}
	if length > 0 {
		variables = make(map[string]interface{}, length)
	}
	if len(view.data) > 0 {
		if len(params) > 0 {
			if variables == nil {
				variables = make(map[string]interface{})
			}
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
	// TODO any graceful plan to replace "<no value>"?
	result := gstr.Replace(buffer.String(), "<no value>", "")
	result = view.i18nTranslate(result, variables)
	return result, nil
}

// getTemplate returns the template object associated with given template folder <path>.
// It uses template cache to enhance performance, that is, it will return the same template object
// with the same given <path>. It will also automatically refresh the template cache
// if the template files under <path> changes (recursively).
func (view *View) getTemplate(path string, pattern string) (tpl *template.Template, err error) {
	r := templates.GetOrSetFuncLock(path, func() interface{} {
		tpl = template.New(path).Delims(view.delimiters[0], view.delimiters[1]).Funcs(view.funcMap)
		// Firstly checking the resource manager.
		if !gres.IsEmpty() {
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
		return tpl
	})
	if r != nil {
		return r.(*template.Template), nil
	}
	return
}

// searchFile returns the found absolute path for <file> and its template folder path.
func (view *View) searchFile(file string) (path string, folder string, resource *gres.File, err error) {
	// Firstly checking the resource manager.
	if !gres.IsEmpty() {
		for _, v := range resourceTryFiles {
			if resource = gres.Get(v + file); resource != nil {
				path = resource.Name()
				folder = v
				return
			}
		}

		view.paths.RLockFunc(func(array []string) {
			for _, v := range array {
				v = strings.TrimRight(v, "/"+gfile.Separator)
				if resource = gres.Get(v + "/" + file); resource != nil {
					path = resource.Name()
					folder = v
					break
				}
				if resource = gres.Get(v + "/template/" + file); resource != nil {
					path = resource.Name()
					folder = v + "/template"
					break
				}
			}
		})
	}

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
