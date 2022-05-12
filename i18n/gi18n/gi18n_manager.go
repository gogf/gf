// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gi18n

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
)

// Manager for i18n contents, it is concurrent safe, supporting hot reload.
type Manager struct {
	mu      sync.RWMutex
	data    map[string]map[string]string // Translating map.
	pattern string                       // Pattern for regex parsing.
	options Options                      // configuration options.
}

// Options is used for i18n object configuration.
type Options struct {
	Path       string   // I18n files storage path.
	Language   string   // Default local language.
	Delimiters []string // Delimiters for variable parsing.
}

var (
	defaultLanguage   = "en"                // defaultDelimiters defines the default language if user does not specified in options.
	defaultDelimiters = []string{"{#", "}"} // defaultDelimiters defines the default key variable delimiters.
)

// New creates and returns a new i18n manager.
// The optional parameter `option` specifies the custom options for i18n manager.
// It uses a default one if it's not passed.
func New(options ...Options) *Manager {
	var opts Options
	if len(options) > 0 {
		opts = options[0]
	} else {
		opts = DefaultOptions()
	}
	if len(opts.Language) == 0 {
		opts.Language = defaultLanguage
	}
	if len(opts.Delimiters) == 0 {
		opts.Delimiters = defaultDelimiters
	}
	m := &Manager{
		options: opts,
		pattern: fmt.Sprintf(
			`%s(\w+)%s`,
			gregex.Quote(opts.Delimiters[0]),
			gregex.Quote(opts.Delimiters[1]),
		),
	}
	intlog.Printf(context.TODO(), `New: %#v`, m)
	return m
}

// DefaultOptions creates and returns a default options for i18n manager.
func DefaultOptions() Options {
	var (
		path        = "i18n"
		realPath, _ = gfile.Search(path)
	)
	if realPath != "" {
		path = realPath
		// To avoid of the source path of GF: github.com/gogf/i18n/gi18n
		if gfile.Exists(path + gfile.Separator + "gi18n") {
			path = ""
		}
	}
	return Options{
		Path:       path,
		Language:   "en",
		Delimiters: defaultDelimiters,
	}
}

// SetPath sets the directory path storing i18n files.
func (m *Manager) SetPath(path string) error {
	if gres.Contains(path) {
		m.options.Path = path
	} else {
		realPath, _ := gfile.Search(path)
		if realPath == "" {
			return gerror.NewCodef(gcode.CodeInvalidParameter, `%s does not exist`, path)
		}
		m.options.Path = realPath
	}
	intlog.Printf(context.TODO(), `SetPath: %s`, m.options.Path)
	return nil
}

// SetLanguage sets the language for translator.
func (m *Manager) SetLanguage(language string) {
	m.options.Language = language
	intlog.Printf(context.TODO(), `SetLanguage: %s`, m.options.Language)
}

// SetDelimiters sets the delimiters for translator.
func (m *Manager) SetDelimiters(left, right string) {
	m.pattern = fmt.Sprintf(`%s(\w+)%s`, gregex.Quote(left), gregex.Quote(right))
	intlog.Printf(context.TODO(), `SetDelimiters: %v`, m.pattern)
}

// T is alias of Translate for convenience.
func (m *Manager) T(ctx context.Context, content string) string {
	return m.Translate(ctx, content)
}

// Tf is alias of TranslateFormat for convenience.
func (m *Manager) Tf(ctx context.Context, format string, values ...interface{}) string {
	return m.TranslateFormat(ctx, format, values...)
}

// TranslateFormat translates, formats and returns the `format` with configured language
// and given `values`.
func (m *Manager) TranslateFormat(ctx context.Context, format string, values ...interface{}) string {
	return fmt.Sprintf(m.Translate(ctx, format), values...)
}

// Translate translates `content` with configured language.
func (m *Manager) Translate(ctx context.Context, content string) string {
	m.init(ctx)
	m.mu.RLock()
	defer m.mu.RUnlock()
	transLang := m.options.Language
	if lang := LanguageFromCtx(ctx); lang != "" {
		transLang = lang
	}
	data := m.data[transLang]
	if data == nil {
		return content
	}
	// Parse content as name.
	if v, ok := data[content]; ok {
		return v
	}
	// Parse content as variables container.
	result, _ := gregex.ReplaceStringFuncMatch(
		m.pattern, content,
		func(match []string) string {
			if v, ok := data[match[1]]; ok {
				return v
			}
			return match[0]
		})
	intlog.Printf(ctx, `Translate for language: %s`, transLang)
	return result
}

// GetContent retrieves and returns the configured content for given key and specified language.
// It returns an empty string if not found.
func (m *Manager) GetContent(ctx context.Context, key string) string {
	m.init(ctx)
	m.mu.RLock()
	defer m.mu.RUnlock()
	transLang := m.options.Language
	if lang := LanguageFromCtx(ctx); lang != "" {
		transLang = lang
	}
	if data, ok := m.data[transLang]; ok {
		return data[key]
	}
	return ""
}

// init initializes the manager for lazy initialization design.
// The i18n manager is only initialized once.
func (m *Manager) init(ctx context.Context) {
	m.mu.RLock()
	// If the data is not nil, means it's already initialized.
	if m.data != nil {
		m.mu.RUnlock()
		return
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if gres.Contains(m.options.Path) {
		files := gres.ScanDirFile(m.options.Path, "*.*", true)
		if len(files) > 0 {
			var (
				path  string
				name  string
				lang  string
				array []string
			)
			m.data = make(map[string]map[string]string)
			for _, file := range files {
				name = file.Name()
				path = name[len(m.options.Path)+1:]
				array = strings.Split(path, "/")
				if len(array) > 1 {
					lang = array[0]
				} else {
					lang = gfile.Name(array[0])
				}
				if m.data[lang] == nil {
					m.data[lang] = make(map[string]string)
				}
				if j, err := gjson.LoadContent(file.Content()); err == nil {
					for k, v := range j.Var().Map() {
						m.data[lang][k] = gconv.String(v)
					}
				} else {
					intlog.Errorf(ctx, "load i18n file '%s' failed: %+v", name, err)
				}
			}
		}
	} else if m.options.Path != "" {
		files, _ := gfile.ScanDirFile(m.options.Path, "*.*", true)
		if len(files) == 0 {
			return
		}
		var (
			path  string
			lang  string
			array []string
		)
		m.data = make(map[string]map[string]string)
		for _, file := range files {
			path = file[len(m.options.Path)+1:]
			array = strings.Split(path, gfile.Separator)
			if len(array) > 1 {
				lang = array[0]
			} else {
				lang = gfile.Name(array[0])
			}
			if m.data[lang] == nil {
				m.data[lang] = make(map[string]string)
			}
			if j, err := gjson.LoadContent(gfile.GetBytes(file)); err == nil {
				for k, v := range j.Var().Map() {
					m.data[lang][k] = gconv.String(v)
				}
			} else {
				intlog.Errorf(ctx, "load i18n file '%s' failed: %+v", file, err)
			}
		}
		// Monitor changes of i18n files for hot reload feature.
		_, _ = gfsnotify.Add(path, func(event *gfsnotify.Event) {
			// Any changes of i18n files, clear the data.
			m.mu.Lock()
			m.data = nil
			m.mu.Unlock()
			gfsnotify.Exit()
		}, ``)
	}
}
