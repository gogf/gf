// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"context"
	"log/slog"
	"slices"
	"strings"
)

// SlogHandler is a slog.Handler that forwards records to a glog Logger.
// It maps slog levels to glog levels and appends the record attributes to the
// message as "key=value" pairs, with attribute groups prefixed by "group.".
type SlogHandler struct {
	logger *Logger
	attrs  []string
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
	var b strings.Builder
	b.WriteString(record.Message)
	for _, attr := range h.attrs {
		b.WriteString(" ")
		b.WriteString(attr)
	}
	var prefix string
	if len(h.groups) > 0 {
		prefix = strings.Join(h.groups, ".") + "."
	}
	record.Attrs(func(attr slog.Attr) bool {
		if line := formatSlogAttr(prefix, attr); line != "" {
			b.WriteString(" ")
			b.WriteString(line)
		}
		return true
	})
	message := b.String()
	switch slogLevelToGlogLevel(record.Level) {
	case LEVEL_DEBU:
		h.logger.Debug(ctx, message)
	case LEVEL_INFO:
		h.logger.Info(ctx, message)
	case LEVEL_NOTI:
		h.logger.Notice(ctx, message)
	case LEVEL_WARN:
		h.logger.Warning(ctx, message)
	default:
		h.logger.Error(ctx, message)
	}
	return nil
}

// WithAttrs returns a new handler whose records include the given attributes.
func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	prefix := ""
	if len(h.groups) > 0 {
		prefix = strings.Join(h.groups, ".") + "."
	}
	formatted := make([]string, 0, len(attrs))
	for _, attr := range attrs {
		if line := formatSlogAttr(prefix, attr); line != "" {
			formatted = append(formatted, line)
		}
	}
	if len(formatted) == 0 {
		return h
	}
	newHandler := *h
	newHandler.attrs = append(slices.Clone(h.attrs), formatted...)
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

// formatSlogAttr formats one attribute as "key=value", expanding nested groups
// into "group.key=value" pairs.
func formatSlogAttr(prefix string, attr slog.Attr) string {
	var b strings.Builder
	appendSlogAttr(&b, prefix, attr)
	return b.String()
}

func appendSlogAttr(b *strings.Builder, prefix string, attr slog.Attr) {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return
	}
	if attr.Value.Kind() == slog.KindGroup {
		if attr.Key != "" {
			prefix += attr.Key + "."
		}
		for _, child := range attr.Value.Group() {
			appendSlogAttr(b, prefix, child)
		}
		return
	}
	if b.Len() > 0 {
		b.WriteString(" ")
	}
	b.WriteString(prefix)
	b.WriteString(attr.Key)
	b.WriteString("=")
	b.WriteString(attr.Value.String())
}
