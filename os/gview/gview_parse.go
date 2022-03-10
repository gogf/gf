// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

import (
	"bytes"
	"context"
	"fmt"
	htmltpl "html/template"
	"strconv"
	texttpl "text/template"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/ghash"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gmlock"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/os/gspath"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gutil"
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
	resourceTryFolders = []string{
		"template/", "template", "/template", "/template/",
		"resource/template/", "resource/template", "/resource/template", "/resource/template/",
	}

	// Prefix array for trying searching in local system.
	localSystemTryFolders = []string{"", "template/", "resource/template"}
)

// Parse parses given template file `file` with given template variables `params`
// and returns the parsed template content.
func (view *View) Parse(ctx context.Context, file string, params ...Params) (result string, err error) {
	var tpl interface{}
	// It caches the file, folder and content to enhance performance.
	r := view.fileCacheMap.GetOrSetFuncLock(file, func() interface{} {
		var (
			path     string
			folder   string
			content  string
			resource *gres.File
		)
		// Searching the absolute file path for `file`.
		path, folder, resource, err = view.searchFile(ctx, file)
		if err != nil {
			return nil
		}
		if resource != nil {
			content = string(resource.Content())
		} else {
			content = gfile.GetContentsWithCache(path)
		}
		// Monitor template files changes using fsnotify asynchronously.
		if resource == nil {
			if _, err = gfsnotify.AddOnce("gview.Parse:"+folder, folder, func(event *gfsnotify.Event) {
				// CLEAR THEM ALL.
				view.fileCacheMap.Clear()
				templates.Clear()
				gfsnotify.Exit()
			}); err != nil {
				intlog.Errorf(ctx, `%+v`, err)
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
	// Get the template object instance for `folder`.
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
	// of the existing `params` or view.data because both variables are pointers.
	// It needs to merge the values of the two maps into a new map.
	variables := gutil.MapMergeCopy(params...)
	if len(view.data) > 0 {
		gutil.MapMerge(variables, view.data)
	}
	view.setI18nLanguageFromCtx(ctx, variables)

	buffer := bytes.NewBuffer(nil)
	if view.config.AutoEncode {
		newTpl, err := tpl.(*htmltpl.Template).Clone()
		if err != nil {
			return "", err
		}
		if err = newTpl.Execute(buffer, variables); err != nil {
			return "", err
		}
	} else {
		if err = tpl.(*texttpl.Template).Execute(buffer, variables); err != nil {
			return "", err
		}
	}

	// TODO any graceful plan to replace "<no value>"?
	result = gstr.Replace(buffer.String(), "<no value>", "")
	result = view.i18nTranslate(ctx, result, variables)
	return result, nil
}

// ParseDefault parses the default template file with params.
func (view *View) ParseDefault(ctx context.Context, params ...Params) (result string, err error) {
	return view.Parse(ctx, view.config.DefaultFile, params...)
}

// ParseContent parses given template content `content`  with template variables `params`
// and returns the parsed content in []byte.
func (view *View) ParseContent(ctx context.Context, content string, params ...Params) (string, error) {
	// It's not necessary continuing parsing if template content is empty.
	if content == "" {
		return "", nil
	}
	var (
		err error
		key = fmt.Sprintf("%s_%v_%v", templateNameForContentParsing, view.config.Delimiters, view.config.AutoEncode)
		tpl = templates.GetOrSetFuncLock(key, func() interface{} {
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
	)
	// Using memory lock to ensure concurrent safety for content parsing.
	hash := strconv.FormatUint(ghash.DJB64([]byte(content)), 10)
	gmlock.LockFunc("gview.ParseContent:"+hash, func() {
		if view.config.AutoEncode {
			tpl, err = tpl.(*htmltpl.Template).Parse(content)
		} else {
			tpl, err = tpl.(*texttpl.Template).Parse(content)
		}
	})
	if err != nil {
		err = gerror.Wrapf(err, `template parsing failed`)
		return "", err
	}
	// Note that the template variable assignment cannot change the value
	// of the existing `params` or view.data because both variables are pointers.
	// It needs to merge the values of the two maps into a new map.
	variables := gutil.MapMergeCopy(params...)
	if len(view.data) > 0 {
		gutil.MapMerge(variables, view.data)
	}
	view.setI18nLanguageFromCtx(ctx, variables)

	buffer := bytes.NewBuffer(nil)
	if view.config.AutoEncode {
		var newTpl *htmltpl.Template
		newTpl, err = tpl.(*htmltpl.Template).Clone()
		if err != nil {
			err = gerror.Wrapf(err, `template clone failed`)
			return "", err
		}
		if err = newTpl.Execute(buffer, variables); err != nil {
			err = gerror.Wrapf(err, `template parsing failed`)
			return "", err
		}
	} else {
		if err = tpl.(*texttpl.Template).Execute(buffer, variables); err != nil {
			err = gerror.Wrapf(err, `template parsing failed`)
			return "", err
		}
	}
	// TODO any graceful plan to replace "<no value>"?
	result := gstr.Replace(buffer.String(), "<no value>", "")
	result = view.i18nTranslate(ctx, result, variables)
	return result, nil
}

// getTemplate returns the template object associated with given template file `path`.
// It uses template cache to enhance performance, that is, it will return the same template object
// with the same given `path`. It will also automatically refresh the template cache
// if the template files under `path` changes (recursively).
func (view *View) getTemplate(filePath, folderPath, pattern string) (tpl interface{}, err error) {
	var (
		mapKey  = fmt.Sprintf("%s_%v", filePath, view.config.Delimiters)
		mapFunc = func() interface{} {
			tplName := filePath
			if view.config.AutoEncode {
				tpl = htmltpl.New(tplName).Delims(
					view.config.Delimiters[0],
					view.config.Delimiters[1],
				).Funcs(view.funcMap)
			} else {
				tpl = texttpl.New(tplName).Delims(
					view.config.Delimiters[0],
					view.config.Delimiters[1],
				).Funcs(view.funcMap)
			}
			// Firstly checking the resource manager.
			if !gres.IsEmpty() {
				if files := gres.ScanDirFile(folderPath, pattern, true); len(files) > 0 {
					if view.config.AutoEncode {
						var t = tpl.(*htmltpl.Template)
						for _, v := range files {
							_, err = t.New(v.FileInfo().Name()).Parse(string(v.Content()))
							if err != nil {
								err = view.formatTemplateObjectCreatingError(v.Name(), tplName, err)
								return nil
							}
						}
					} else {
						var t = tpl.(*texttpl.Template)
						for _, v := range files {
							_, err = t.New(v.FileInfo().Name()).Parse(string(v.Content()))
							if err != nil {
								err = view.formatTemplateObjectCreatingError(v.Name(), tplName, err)
								return nil
							}
						}
					}
					return tpl
				}
			}

			// Secondly checking the file system,
			// and then automatically parsing all its sub-files recursively.
			var files []string
			files, err = gfile.ScanDir(folderPath, pattern, true)
			if err != nil {
				return nil
			}
			if view.config.AutoEncode {
				t := tpl.(*htmltpl.Template)
				for _, file := range files {
					if _, err = t.Parse(gfile.GetContents(file)); err != nil {
						err = view.formatTemplateObjectCreatingError(file, tplName, err)
						return nil
					}
				}
			} else {
				t := tpl.(*texttpl.Template)
				for _, file := range files {
					if _, err = t.Parse(gfile.GetContents(file)); err != nil {
						err = view.formatTemplateObjectCreatingError(file, tplName, err)
						return nil
					}
				}
			}
			return tpl
		}
	)
	result := templates.GetOrSetFuncLock(mapKey, mapFunc)
	if result != nil {
		return result, nil
	}
	return
}

// formatTemplateObjectCreatingError formats the error that created from creating template object.
func (view *View) formatTemplateObjectCreatingError(filePath, tplName string, err error) error {
	if err != nil {
		return gerror.NewSkip(1, gstr.Replace(err.Error(), tplName, filePath))
	}
	return nil
}

// searchFile returns the found absolute path for `file` and its template folder path.
// Note that, the returned `folder` is the template folder path, but not the folder of
// the returned template file `path`.
func (view *View) searchFile(ctx context.Context, file string) (path string, folder string, resource *gres.File, err error) {
	var tempPath string
	// Firstly checking the resource manager.
	if !gres.IsEmpty() {
		// Try folders.
		for _, tryFolder := range resourceTryFolders {
			tempPath = tryFolder + file
			if resource = gres.Get(tempPath); resource != nil {
				path = resource.Name()
				folder = tryFolder
				return
			}
		}
		// Search folders.
		view.searchPaths.RLockFunc(func(array []string) {
			for _, searchPath := range array {
				for _, tryFolder := range resourceTryFolders {
					tempPath = searchPath + tryFolder + file
					if resFile := gres.Get(tempPath); resFile != nil {
						path = resFile.Name()
						folder = searchPath + tryFolder
						return
					}
				}
			}
		})
	}

	// Secondly checking the file system.
	if path == "" {
		// Absolute path.
		path = gfile.RealPath(file)
		if path != "" {
			folder = gfile.Dir(path)
			return
		}
		// In search paths.
		view.searchPaths.RLockFunc(func(array []string) {
			for _, searchPath := range array {
				searchPath = gstr.TrimRight(searchPath, `\/`)
				for _, tryFolder := range localSystemTryFolders {
					relativePath := gstr.TrimRight(
						gfile.Join(tryFolder, file),
						`\/`,
					)
					if path, _ = gspath.Search(searchPath, relativePath); path != "" {
						folder = gfile.Join(searchPath, tryFolder)
						return
					}
				}
			}
		})
	}

	// Error checking.
	if path == "" {
		buffer := bytes.NewBuffer(nil)
		if view.searchPaths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("cannot find template file \"%s\" in following paths:", file))
			view.searchPaths.RLockFunc(func(array []string) {
				index := 1
				for _, searchPath := range array {
					searchPath = gstr.TrimRight(searchPath, `\/`)
					for _, tryFolder := range localSystemTryFolders {
						buffer.WriteString(fmt.Sprintf(
							"\n%d. %s",
							index, gfile.Join(searchPath, tryFolder),
						))
						index++
					}
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf("cannot find template file \"%s\" with no path set/add", file))
		}
		if errorPrint() {
			glog.Error(ctx, buffer.String())
		}
		err = gerror.NewCodef(gcode.CodeInvalidParameter, `template file "%s" not found`, file)
	}
	return
}
