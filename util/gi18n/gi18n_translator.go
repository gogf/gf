// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gi18n

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gogf/gf/os/gfsnotify"

	"github.com/gogf/gf/text/gregex"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/encoding/gjson"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gres"
)

// Translator, it is concurrent safe, supporting hot reload.
type Translator struct {
	mu      sync.RWMutex
	data    map[string]map[string]string // Translating map.
	pattern string                       // Pattern for regex parsing.
	options Options                      // configuration options.
}

type Options struct {
	Path       string   // I18n files storage path.
	Language   string   // Local language.
	Delimiters []string // Delimiters for variable parsing.
}

var (
	defaultDelimiters = []string{"{#", "}"}
)

func New(options ...Options) *Translator {
	var opts Options
	if len(options) > 0 {
		opts = options[0]
	} else {
		opts = DefaultOptions()
	}
	if len(opts.Delimiters) == 0 {
		opts.Delimiters = defaultDelimiters
	}
	return &Translator{
		options: opts,
		pattern: fmt.Sprintf(
			`%s(\w+)%s`,
			gregex.Quote(opts.Delimiters[0]),
			gregex.Quote(opts.Delimiters[1]),
		),
	}
}

func DefaultOptions() Options {
	path := "i18n"
	realPath, _ := gfile.Search(path)
	if realPath != "" {
		path = realPath
	} else {
		path = "/" + path
	}
	return Options{
		Path:       path,
		Delimiters: []string{"{#", "}"},
	}
}

// SetPath sets the directory path storing i18n files.
func (t *Translator) SetPath(path string) error {
	if gres.Contains(path) {
		t.options.Path = path
	} else {
		realPath, _ := gfile.Search(path)
		if realPath == "" {
			return errors.New(fmt.Sprintf(`%s does not exist`, path))
		}
		t.options.Path = realPath
	}
	return nil
}

// SetLanguage sets the language for translator.
func (t *Translator) SetLanguage(language string) {
	t.options.Language = language
}

// SetDelimiters sets the delimiters for translator.
func (t *Translator) SetDelimiters(left, right string) {
	t.pattern = fmt.Sprintf(`%s(\w+)%s`, gregex.Quote(left), gregex.Quote(right))
}

// T is alias of Translate.
func (t *Translator) T(content string, language ...string) string {
	return t.Translate(content, language...)
}

// Translate translates <content> with configured language.
// The parameter <language> specifies custom translation language ignoring configured language.
func (t *Translator) Translate(content string, language ...string) string {
	t.init()
	t.mu.RLock()
	defer t.mu.RUnlock()
	var data map[string]string
	if len(language) > 0 {
		data = t.data[language[0]]
	} else {
		data = t.data[t.options.Language]
	}
	if data == nil {
		return content
	}
	// Parse content as name.
	if v, ok := data[content]; ok {
		return v
	}
	// Parse content as variables container.
	result, _ := gregex.ReplaceStringFuncMatch(t.pattern, content, func(match []string) string {
		if v, ok := data[match[1]]; ok {
			return v
		}
		return match[0]
	})
	return result
}

func (t *Translator) init() {
	t.mu.RLock()
	if t.data != nil {
		t.mu.RUnlock()
		return
	}
	t.mu.RUnlock()

	t.mu.Lock()
	defer t.mu.Unlock()
	if gres.Contains(t.options.Path) {
		files := gres.ScanDirFile(t.options.Path, "*.*", true)
		if len(files) > 0 {
			var path string
			var name string
			var lang string
			var array []string
			t.data = make(map[string]map[string]string)
			for _, file := range files {
				name = file.Name()
				path = name[len(t.options.Path)+1:]
				array = strings.Split(path, "/")
				if len(array) > 1 {
					lang = array[0]
				} else {
					lang = gfile.Name(array[0])
				}
				if t.data[lang] == nil {
					t.data[lang] = make(map[string]string)
				}
				j, _ := gjson.LoadContent(file.Content())
				if j != nil {
					for k, v := range j.ToMap() {
						t.data[lang][k] = gconv.String(v)
					}
				}
			}
		}
	} else {
		files, _ := gfile.ScanDirFile(t.options.Path, "*.*", true)
		if len(files) > 0 {
			var path string
			var lang string
			var array []string
			t.data = make(map[string]map[string]string)
			for _, file := range files {
				path = file[len(t.options.Path)+1:]
				array = strings.Split(path, "/")
				if len(array) > 1 {
					lang = array[0]
				} else {
					lang = gfile.Name(array[0])
				}
				if t.data[lang] == nil {
					t.data[lang] = make(map[string]string)
				}
				j, _ := gjson.LoadContent(gfile.GetBytes(file))
				if j != nil {
					for k, v := range j.ToMap() {
						t.data[lang][k] = gconv.String(v)
					}
				}
			}
			_, _ = gfsnotify.Add(path, func(event *gfsnotify.Event) {
				t.mu.Lock()
				t.data = nil
				t.mu.Unlock()
				gfsnotify.Exit()
			})
		}
	}
}
