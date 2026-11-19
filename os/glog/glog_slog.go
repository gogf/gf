// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"slices"
	"strings"
)

// SlogHandler is a slog.Handler that forwards records to a glog Logger.
// It maps slog levels to glog levels and forwards the record attributes as
// key-value pairs of the logging values, so the configured glog handlers format
// them the same way as direct key-value logging calls.
type SlogHandler struct {
	logger *Logger
	attrs  []any
	groups []string
}

// NewSlogHandler creates and returns a slog.Handler that forwards records to `logger`.
// It uses the default logger if `logger` is not given or is nil.
func NewSlogHandler(logger ...*Logger) *SlogHandler {
	l := DefaultLogger()
	if len(logger) > 0 && logger[0] != nil {
		l = logger[0]
	}
	return &SlogHandler{logger: l}
}

// Enabled reports whether the handler handles records at the given level,
// according to the level configuration of the underlying logger.
func (h *SlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.logger.checkLevel(slogLevelToGlogLevel(level))
}

// Handle forwards the record to the underlying logger.
func (h *SlogHandler) Handle(ctx context.Context, record slog.Record) error {
	level := slogLevelToGlogLevel(record.Level)
	if !h.logger.checkLevel(level) {
		return nil
	}
	values := make([]any, 0, 1+len(h.attrs)+2*record.NumAttrs())
	values = append(values, record.Message)
	values = append(values, h.attrs...)
	var prefix string
	if len(h.groups) > 0 {
		prefix = strings.Join(h.groups, ".") + "."
	}
	record.Attrs(func(attr slog.Attr) bool {
		values = appendSlogAttr(values, prefix, attr)
		return true
	})
	var stack string
	if level == LEVEL_ERRO && h.logger.config.StStatus == 1 {
		if record.PC != 0 {
			stack = slogStackFromPC(record.PC, h.logger.stackFilters())
		} else {
			stack = h.logger.GetStack()
		}
	}
	h.logger.printWithCaller(ctx, level, record.PC, stack, values...)
	return nil
}

// WithAttrs returns a new handler whose records include the given attributes.
func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	var prefix string
	if len(h.groups) > 0 {
		prefix = strings.Join(h.groups, ".") + "."
	}
	pairs := make([]any, 0, 2*len(attrs))
	for _, attr := range attrs {
		pairs = appendSlogAttr(pairs, prefix, attr)
	}
	if len(pairs) == 0 {
		return h
	}
	newHandler := *h
	newHandler.attrs = append(slices.Clone(h.attrs), pairs...)
	return &newHandler
}

// WithGroup returns a new handler that prefixes the following attributes with the group name.
func (h *SlogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newHandler := *h
	newHandler.groups = append(slices.Clone(h.groups), name)
	return &newHandler
}

// slogLevelToGlogLevel maps a slog level to the level of the glog package.
func slogLevelToGlogLevel(level slog.Level) int {
	switch {
	case level < slog.LevelInfo:
		return LEVEL_DEBU
	case level < slog.LevelWarn:
		if level == slog.LevelInfo {
			return LEVEL_INFO
		}
		return LEVEL_NOTI
	case level < slog.LevelError:
		return LEVEL_WARN
	default:
		return LEVEL_ERRO
	}
}

// appendSlogAttr appends `attr` to `values` as key-value pairs, expanding nested
// groups into "group.key" keys.
func appendSlogAttr(values []any, prefix string, attr slog.Attr) []any {
	if attr.Equal(slog.Attr{}) {
		return values
	}
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return values
	}
	if attr.Value.Kind() == slog.KindGroup {
		if attr.Key != "" {
			prefix += attr.Key + "."
		}
		for _, child := range attr.Value.Group() {
			values = appendSlogAttr(values, prefix, child)
		}
		return values
	}
	return append(values, prefix+attr.Key, attr.Value.String())
}

// slogStackFromPC returns the stack content beginning at the frame of `pc`, which is
// the user call site recorded in the slog record. Starting from the record PC keeps
// the log/slog and SlogHandler frames out of the stack even in -trimpath builds, where
// runtime.GOROOT() is empty and GOROOT based filtering takes no effect.
func slogStackFromPC(pc uintptr, filters []string) string {
	pcs := make([]uintptr, 64)
	for {
		if n := runtime.Callers(2, pcs); n < len(pcs) {
			pcs = pcs[:n]
			break
		}
		pcs = make([]uintptr, 2*len(pcs))
	}
	start := slices.Index(pcs, pc)
	if start < 0 {
		return ""
	}
	var (
		index  = 1
		space  = "  "
		buffer = bytes.NewBuffer(nil)
		goRoot = strings.ReplaceAll(runtime.GOROOT(), "\\", "/")
		frames = runtime.CallersFrames(pcs[start:])
	)
	for {
		frame, more := frames.Next()
		if !slogStackFrameFiltered(frame.File, filters, goRoot) {
			function := frame.Function
			if function == "" {
				function = "unknown"
			}
			if index > 9 {
				space = " "
			}
			fmt.Fprintf(buffer, "%d.%s%s\n    %s:%d\n", index, space, function, frame.File, frame.Line)
			index++
		}
		if !more {
			break
		}
	}
	return buffer.String()
}

// slogStackFrameFiltered reports whether the stack frame of `file` should be filtered out.
func slogStackFrameFiltered(file string, filters []string, goRoot string) bool {
	if file == "" {
		return true
	}
	for _, filter := range filters {
		if filter != "" && strings.Contains(file, filter) {
			return true
		}
	}
	if goRoot != "" {
		file = strings.ReplaceAll(file, "\\", "/")
		if strings.HasPrefix(file, goRoot) && len(file) > len(goRoot) {
			if separator := file[len(goRoot)]; separator == '/' || separator == '\\' {
				return true
			}
		}
	}
	return false
}
