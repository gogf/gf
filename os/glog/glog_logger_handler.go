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

// HandlerInput is the input parameter struct for logging Handler.
type HandlerInput struct {
	Logger       *Logger       // Logger.
	Buffer       *bytes.Buffer // Buffer for logging content outputs.
	Time         time.Time     // Logging time, which is the time that logging triggers.
	TimeFormat   string        // Formatted time string, like "2016-01-09 12:00:00".
	Color        int           // Using color, like COLOR_RED, COLOR_BLUE, etc.
	Level        int           // Using level, like LEVEL_INFO, LEVEL_ERRO, etc.
	LevelFormat  string        // Formatted level string, like "DEBU", "ERRO", etc.
	CallerFunc   string        // The source function name that calls logging.
	CallerPath   string        // The source file path and its line number that calls logging.
	CtxStr       string        // The retrieved context value string from context.
	Prefix       string        // Custom prefix string for logging content.
	Content      string        // Content is the main logging content that passed by you.
	IsAsync      bool          // IsAsync marks it is in asynchronous logging.
	handlerIndex int           // Middleware handling index for internal usage.
}

// Next calls the next logging handler in middleware way.
func (i *HandlerInput) Next(ctx context.Context) {
	if len(i.Logger.config.Handlers)-1 > i.handlerIndex {
		i.handlerIndex++
		i.Logger.config.Handlers[i.handlerIndex](ctx, i)
	} else {
		defaultHandler(ctx, i)
	}
}

// String returns the logging content formatted by default logging handler.
func (i *HandlerInput) String(withColor ...bool) string {
	formatWithColor := false
	if len(withColor) > 0 {
		formatWithColor = withColor[0]
	}
	return i.getDefaultBuffer(formatWithColor).String()
}

func (i *HandlerInput) getDefaultBuffer(withColor bool) *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	if i.TimeFormat != "" {
		buffer.WriteString(i.TimeFormat)
	}
	if i.LevelFormat != "" {
		if withColor {
			i.addStringToBuffer(buffer, i.Logger.getColoredStr(
				i.Logger.getColorByLevel(i.Level), i.LevelFormat,
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
	// avoid a single space at the end of a line.
	buffer.WriteString("\n")
	return buffer
}

func (i *HandlerInput) getRealBuffer(withColor bool) *bytes.Buffer {
	if i.Buffer.Len() > 0 {
		return i.Buffer
	}
	return i.getDefaultBuffer(withColor)
}

// defaultHandler is the default handler for logger.
func defaultHandler(ctx context.Context, in *HandlerInput) {
	buffer := in.Logger.doDefaultPrint(ctx, in)
	if in.Buffer.Len() == 0 {
		in.Buffer = buffer
	}
}

func (i *HandlerInput) addStringToBuffer(buffer *bytes.Buffer, strings ...string) {
	for _, s := range strings {
		if buffer.Len() > 0 {
			buffer.WriteByte(' ')
		}
		buffer.WriteString(s)
	}
}
