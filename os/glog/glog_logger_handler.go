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

	"github.com/gogf/gf/v2/util/gconv"
)

// Handler is function handler for custom logging content outputs.
type Handler func(ctx context.Context, in *HandlerInput)

// HandlerInput is the input parameter struct for logging Handler.
type HandlerInput struct {
	internalHandlerInfo
	Logger      *Logger       // Current Logger object.
	Buffer      *bytes.Buffer // Buffer for logging content outputs.
	Time        time.Time     // Logging time, which is the time that logging triggers.
	TimeFormat  string        // Formatted time string, like "2016-01-09 12:00:00".
	Color       int           // Using color, like COLOR_RED, COLOR_BLUE, etc. Eg: 34
	Level       int           // Using level, like LEVEL_INFO, LEVEL_ERRO, etc. Eg: 256
	LevelFormat string        // Formatted level string, like "DEBU", "ERRO", etc. Eg: ERRO
	CallerFunc  string        // The source function name that calls logging, only available if F_CALLER_FN set.
	CallerPath  string        // The source file path and its line number that calls logging, only available if F_FILE_SHORT or F_FILE_LONG set.
	CtxStr      string        // The retrieved context value string from context, only available if Config.CtxKeys configured.
	TraceId     string        // Trace id, only available if OpenTelemetry is enabled.
	Prefix      string        // Custom prefix string for logging content.
	Content     string        // Content is the main logging content without error stack string produced by logger.
	Values      []any         // The passed un-formatted values array to logger.
	Stack       string        // Stack string produced by logger, only available if Config.StStatus configured.
	IsAsync     bool          // IsAsync marks it is in asynchronous logging.
}

type internalHandlerInfo struct {
	index    int       // Middleware handling index for internal usage.
	handlers []Handler // Handler array calling bu index.
}

// defaultHandler is the default handler for package.
var defaultHandler Handler

// doFinalPrint is a handler for logging content printing.
// This handler outputs logging content to file/stdout/write if any of them configured.
func doFinalPrint(ctx context.Context, in *HandlerInput) {
	buffer := in.Logger.doFinalPrint(ctx, in)
	if in.Buffer.Len() == 0 {
		in.Buffer = buffer
	}
}

// SetDefaultHandler sets default handler for package.
func SetDefaultHandler(handler Handler) {
	defaultHandler = handler
}

// GetDefaultHandler returns the default handler of package.
func GetDefaultHandler() Handler {
	return defaultHandler
}

// Next calls the next logging handler in middleware way.
func (in *HandlerInput) Next(ctx context.Context) {
	in.index++
	if in.index < len(in.handlers) {
		in.handlers[in.index](ctx, in)
	}
}

// String returns the logging content formatted by default logging handler.
func (in *HandlerInput) String(withColor ...bool) string {
	formatWithColor := false
	if len(withColor) > 0 {
		formatWithColor = withColor[0]
	}
	return in.getDefaultBuffer(formatWithColor).String()
}

func (in *HandlerInput) getDefaultBuffer(withColor bool) *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	if in.Logger.config.HeaderPrint {
		if in.TimeFormat != "" {
			buffer.WriteString(in.TimeFormat)
		}
		if in.Logger.config.LevelPrint && in.LevelFormat != "" {
			var levelStr = "[" + in.LevelFormat + "]"
			if withColor {
				in.addStringToBuffer(buffer, in.Logger.getColoredStr(
					in.Logger.getColorByLevel(in.Level), levelStr,
				))
			} else {
				in.addStringToBuffer(buffer, levelStr)
			}
		}
	}
	if in.TraceId != "" {
		in.addStringToBuffer(buffer, "{"+in.TraceId+"}")
	}
	if in.CtxStr != "" {
		in.addStringToBuffer(buffer, "{"+in.CtxStr+"}")
	}
	if in.Logger.config.HeaderPrint {
		if in.Prefix != "" {
			in.addStringToBuffer(buffer, in.Prefix)
		}
		if in.CallerFunc != "" {
			in.addStringToBuffer(buffer, in.CallerFunc)
		}
		if in.CallerPath != "" {
			in.addStringToBuffer(buffer, in.CallerPath)
		}
	}

	if in.Content != "" {
		in.addStringToBuffer(buffer, in.Content)
	}

	// Convert values string content.
	var valueContent string
	for _, v := range in.Values {
		valueContent = gconv.String(v)
		if len(valueContent) == 0 {
			continue
		}
		if buffer.Len() > 0 {
			if buffer.Bytes()[buffer.Len()-1] == '\n' {
				// Remove one blank line(\n\n).
				if valueContent[0] == '\n' {
					valueContent = valueContent[1:]
				}
				buffer.WriteString(valueContent)
			} else {
				buffer.WriteString(" " + valueContent)
			}
		} else {
			buffer.WriteString(valueContent)
		}
	}

	if in.Stack != "" {
		in.addStringToBuffer(buffer, "\nStack:\n"+in.Stack)
	}
	// avoid a single space at the end of a line.
	buffer.WriteString("\n")
	return buffer
}

func (in *HandlerInput) getRealBuffer(withColor bool) *bytes.Buffer {
	if in.Buffer.Len() > 0 {
		return in.Buffer
	}
	return in.getDefaultBuffer(withColor)
}

func (in *HandlerInput) addStringToBuffer(buffer *bytes.Buffer, strings ...string) {
	for _, s := range strings {
		if buffer.Len() > 0 {
			buffer.WriteByte(' ')
		}
		buffer.WriteString(s)
	}
}
