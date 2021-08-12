// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"context"
	"time"
)

// Handler is function handler for custom logging content outputs.
type Handler func(ctx context.Context, in *HandlerInput)

type HandlerInput struct {
	logger      *Logger         // Logger.
	index       int             // Middleware handling index for internal usage.
	Ctx         context.Context // Context.
	Time        time.Time       // Logging time, which is the time that logging triggers.
	TimeFormat  string          // Formatted time string, like "2016-01-09 12:00:00".
	Color       int             // Using color, like COLOR_RED, COLOR_BLUE, etc.
	Level       int             // Using level, like LEVEL_INFO, LEVEL_ERRO, etc.
	LevelFormat string          // Formatted level string, like "DEBU", "ERRO", etc.
	CallerFunc  string          // The source function name that calls logging.
	CallerPath  string          // The source file path and its line number that calls logging.
	CtxStr      string          // The retrieved context value string from context.
	Prefix      string          // Custom prefix string for logging content.
	Content     string          // Content is the main logging content that passed by you.
	IsAsync     bool            // IsAsync marks it is in asynchronous logging.
}

// defaultHandler is the default handler for logger.
func defaultHandler(ctx context.Context, in *HandlerInput) {
	in.logger.doPrint(ctx, in)
}

func (i *HandlerInput) addStringToBuffer(buffer *bytes.Buffer, strings ...string) {
	for _, s := range strings {
		if buffer.Len() > 0 {
			buffer.WriteByte(' ')
		}
		buffer.WriteString(s)
	}
}

// Buffer creates and returns a buffer that handled by default logging content handler.
func (i *HandlerInput) Buffer() *bytes.Buffer {
	return i.getBuffer(false)
}

func (i *HandlerInput) getBuffer(withColor bool) *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	if i.TimeFormat != "" {
		buffer.WriteString(i.TimeFormat)
	}
	if i.LevelFormat != "" {
		if withColor {
			i.addStringToBuffer(buffer, i.logger.getColoredStr(
				i.logger.getColorByLevel(i.Level), i.LevelFormat,
			))
		} else {
			i.addStringToBuffer(buffer, i.LevelFormat)
		}
	}
	if i.Prefix != "" {
		i.addStringToBuffer(buffer, i.Prefix)
	}
	if i.CtxStr != "" {
		i.addStringToBuffer(buffer, i.CtxStr)
	}
	if i.CallerFunc != "" {
		i.addStringToBuffer(buffer, i.CallerFunc)
	}
	if i.CallerPath != "" {
		i.addStringToBuffer(buffer, i.CallerPath)
	}
	if i.Content != "" {
		i.addStringToBuffer(buffer, i.Content)
	}
	i.addStringToBuffer(buffer, "\n")
	return buffer
}

// String retrieves and returns the logging content handled by default handler.
func (i *HandlerInput) String() string {
	return i.Buffer().String()
}

// Next calls the next logging handler in middleware way.
func (i *HandlerInput) Next() {
	if len(i.logger.config.Handlers)-1 > i.index {
		i.index++
		i.logger.config.Handlers[i.index](i.Ctx, i)
	} else {
		// The last handler is the default handler.
		defaultHandler(i.Ctx, i)
	}
}
