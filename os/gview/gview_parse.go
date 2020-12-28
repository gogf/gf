// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
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
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gfsnotify"
	"github.com/gogf/gf/os/gmlock"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	htmltpl "html/template"
	"strconv"
	"strings"
	texttpl "text/template"

	"github.com/gogf/gf/os/gres"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gspath"
)

const (
	// Template name for content parsing.
	templateNameForContentParsing = "TemplateContent"
)

// fileCacheItem is the cache item for template file.
type fileCacheItem struct {
	path    string
	folder  string
	content string
}

var (
	// Templates cache map for template folder.
	// Note that there's no expiring logic for this map.
	templates = gmap.NewStrAnyMap(true)

	// Try-folders for resource template file searching.
	resourceTryFolders = []string{"template/", "template", "/template", "/template/"}
)

// Parse parses given template file <file> with given template variables <params>
// and returns the parsed template content.
func (view *View) Parse(file string, params ...Params) (result string, err error) {
	var tpl interface{}
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
			content = gfile.GetContentsWithCache(path)
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
	// It's not necessary continuing parsing if template content is empty.
	if item.content == "" {
		return "", nil
	}
	// Get the template object instance for <folder>.
	tpl, err = view.getTemplate(item.path, item.folder, fmt.Sprintf(`*%s`, gfile.Ext(item.path)))
	if err != nil {
		return "", err
	}
	// Using memory lock to ensure concurrent safety for template parsing.
	gmlock.LockFunc("gview.Parse:"+item.path, func() {
		if view.config.AutoEncode {
			tpl, err = tpl.(*htmltpl.Template).Parse(item.content)
		} else {
			tpl, err = tpl.(*texttpl.Template).Parse(item.content)
		}
		if err != nil && item.path != "" {
			err = gerror.Wrap(err, item.path)
		}
	})
	if err != nil {
		return "", err
	}
	// Note that the template variable assignment cannot change the value
	// of the existing <params> or view.data because both variables are pointers.
	// It needs to merge the values of the two maps into a new map.
	variables := gutil.MapMergeCopy(params...)
	if len(view.data) > 0 {
		gutil.MapMerge(variables, view.data)
	}
	buffer := bytes.NewBuffer(nil)
	if view.config.AutoEncode {
		newTpl, err := tpl.(*htmltpl.Template).Clone()
		if err != nil {
			return "", err
		}
		if err := newTpl.Execute(buffer, variables); err != nil {
			return "", err
		}
	} else {
		if err := tpl.(*texttpl.Template).Execute(buffer, variables); err != nil {
			return "", err
		}
	}

	// TODO any graceful plan to replace "<no value>"?
	result = gstr.Replace(buffer.String(), "<no value>", "")
	result = view.i18nTranslate(result, variables)
	return result, nil
}

// ParseDefault parses the default template file with params.
func (view *View) ParseDefault(params ...Params) (result string, err error) {
	return view.Parse(view.config.DefaultFile, params...)
}

// ParseContent parses given template content <content>  with template variables <params>
// and returns the parsed content in []byte.
func (view *View) ParseContent(content string, params ...Params) (string, error) {
	// It's not necessary continuing parsing if template content is empty.
	if content == "" {
		return "", nil
	}
	err := (error)(nil)
	key := fmt.Sprintf("%s_%v_%v", templateNameForContentParsing, view.config.Delimiters, view.config.AutoEncode)
	tpl := templates.GetOrSetFuncLock(key, func() interface{} {
		if view.config.AutoEncode {
			return htmltpl.New(templateNameForContentParsing).Delims(
				view.config.Delimiters[0],
				view.config.Delimiters[1],
			).Funcs(view.funcMap)
		}
		return texttpl.New(templateNameForContentParsing).Delims(
			view.config.Delimiters[0],
			view.config.Delimiters[1],
		).Funcs(view.funcMap)
	})
	// Using memory lock to ensure concurrent safety for content parsing.
	hash := strconv.FormatUint(ghash.DJBHash64([]byte(content)), 10)
	gmlock.LockFunc("gview.ParseContent:"+hash, func() {
		if view.config.AutoEncode {
			tpl, err = tpl.(*htmltpl.Template).Parse(content)
		} else {
			tpl, err = tpl.(*texttpl.Template).Parse(content)
		}
	})
	if err != nil {
		return "", err
	}
	// Note that the template variable assignment cannot change the value
	// of the existing <params> or view.data because both variables are pointers.
	// It needs to merge the values of the two maps into a new map.
	variables := gutil.MapMergeCopy(params...)
	if len(view.data) > 0 {
		gutil.MapMerge(variables, view.data)
	}
	buffer := bytes.NewBuffer(nil)
	if view.config.AutoEncode {
		newTpl, err := tpl.(*htmltpl.Template).Clone()
		if err != nil {
			return "", err
		}
		if err := newTpl.Execute(buffer, variables); err != nil {
			return "", err
		}
	} else {
		if err := tpl.(*texttpl.Template).Execute(buffer, variables); err != nil {
			return "", err
		}
	}
	// TODO any graceful plan to replace "<no value>"?
	result := gstr.Replace(buffer.String(), "<no value>", "")
	result = view.i18nTranslate(result, variables)
	return result, nil
}

// getTemplate returns the template object associated with given template file <path>.
// It uses template cache to enhance performance, that is, it will return the same template object
// with the same given <path>. It will also automatically refresh the template cache
// if the template files under <path> changes (recursively).
func (view *View) getTemplate(filePath, folderPath, pattern string) (tpl interface{}, err error) {
	// Key for template cache.
	key := fmt.Sprintf("%s_%v", filePath, view.config.Delimiters)
	result := templates.GetOrSetFuncLock(key, func() interface{} {
		// Do not use <key> but the <filePath> as the parameter <name> for function New,
		// because when error occurs the <name> will be printed out for error locating.
		if view.config.AutoEncode {
			tpl = htmltpl.New(filePath).Delims(
				view.config.Delimiters[0],
				view.config.Delimiters[1],
			).Funcs(view.funcMap)
		} else {
			tpl = texttpl.New(filePath).Delims(
				view.config.Delimiters[0],
				view.config.Delimiters[1],
			).Funcs(view.funcMap)
		}
		// Firstly checking the resource manager.
		if !gres.IsEmpty() {
			if files := gres.ScanDirFile(folderPath, pattern, true); len(files) > 0 {
				var err error
				for _, v := range files {
					if view.config.AutoEncode {
						_, err = tpl.(*htmltpl.Template).New(v.FileInfo().Name()).Parse(string(v.Content()))
						if err != nil {
							intlog.Error(err)
						}
					} else {
						_, err = tpl.(*texttpl.Template).New(v.FileInfo().Name()).Parse(string(v.Content()))
						if err != nil {
							intlog.Error(err)
						}
					}
				}
				return tpl
			}
		}

		// Secondly checking the file system.
		var (
			files []string
		)
		files, err = gfile.ScanDir(folderPath, pattern, true)
		if err != nil {
			return nil
		}
		if view.config.AutoEncode {
			t := tpl.(*htmltpl.Template)
			for _, file := range files {
				_, err = t.Parse(gfile.GetContents(file))
				if err != nil {
					return nil
				}
			}
		} else {
			t := tpl.(*texttpl.Template)
			for _, file := range files {
				_, err = t.Parse(gfile.GetContents(file))
				if err != nil {
					return nil
				}
			}
		}
		return tpl
	})
	if result != nil {
		return result, nil
	}
	return
}

// searchFile returns the found absolute path for <file> and its template folder path.
// Note that, the returned <folder> is the template folder path, but not the folder of
// the returned template file <path>.
func (view *View) searchFile(file string) (path string, folder string, resource *gres.File, err error) {
	// Firstly checking the resource manager.
	if !gres.IsEmpty() {
		// Try folders.
		for _, folderPath := range resourceTryFolders {
			if resource = gres.Get(folderPath + file); resource != nil {
				path = resource.Name()
				folder = folderPath
				return
			}
		}
		// Search folders.
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
			for _, folderPath := range array {
				folderPath = strings.TrimRight(folderPath, gfile.Separator)
				if path, _ = gspath.Search(folderPath, file); path != "" {
					folder = folderPath
					break
				}
				if path, _ = gspath.Search(folderPath+gfile.Separator+"template", file); path != "" {
					folder = folderPath + gfile.Separator + "template"
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
				for _, folderPath := range array {
					folderPath = strings.TrimRight(folderPath, "/")
					if folderPath == "" {
						folderPath = "/"
					}
					buffer.WriteString(fmt.Sprintf("\n%d. %s", index, folderPath))
					index++
					buffer.WriteString(fmt.Sprintf("\n%d. %s", index, strings.TrimRight(folderPath, "/")+gfile.Separator+"template"))
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
