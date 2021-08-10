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

type Handler func(ctx context.Context, input *HandlerInput)

type HandlerInput struct {
	logger      *Logger
	index       int
	Ctx         context.Context
	Time        time.Time
	TimeFormat  string
	Color       int
	Level       int
	LevelFormat string
	CallerFunc  string
	CallerPath  string
	CtxStr      string
	Prefix      string
	Content     string
	IsAsync     bool
}

// defaultHandler is the default handler for logger.
func defaultHandler(ctx context.Context, input *HandlerInput) {
	input.logger.doPrint(ctx, input)
}

func (i *HandlerInput) addStringToBuffer(buffer *bytes.Buffer, strings ...string) {
	for _, s := range strings {
		if buffer.Len() > 0 {
			buffer.WriteByte(' ')
		}
		buffer.WriteString(s)
	}
}

func (i *HandlerInput) Buffer() *bytes.Buffer {
	return i.getBuffer(false)
}

func (i *HandlerInput) getBuffer(withColor bool) *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(i.TimeFormat)
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

func (i *HandlerInput) String() string {
	return i.Buffer().String()
}

func (i *HandlerInput) Next() {
	if len(i.logger.config.Handlers)-1 > i.index {
		i.index++
		i.logger.config.Handlers[i.index](i.Ctx, i)
	} else {
		// The last handler is the default handler.
		defaultHandler(i.Ctx, i)
	}
}
