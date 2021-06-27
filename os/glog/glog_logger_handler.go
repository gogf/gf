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
	"time"
)

type Handler func(ctx context.Context, input *HandlerInput)

type HandlerInput struct {
	logger      *Logger
	index       int
	Ctx         context.Context
	Time        time.Time
	TimeFormat  string
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
	input.logger.printToWriter(ctx, input)
}

func (i *HandlerInput) addStringToBuffer(buffer *bytes.Buffer, s string) {
	if buffer.Len() > 0 {
		buffer.WriteByte(' ')
	}
	buffer.WriteString(s)
}

func (i *HandlerInput) Buffer(withColor ...bool) *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(i.TimeFormat)
	if i.LevelFormat != "" {
		if i.logger.config.FileColor || (len(withColor) > 0 && withColor[0] == mustWithColor) {
			i.addStringToBuffer(buffer, i.getLevelFormatWithColor())
		} else {
			i.addStringToBuffer(buffer, i.LevelFormat)
		}
	}
	if i.CallerFunc != "" {
		i.addStringToBuffer(buffer, i.CallerFunc)
	}
	if i.CallerPath != "" {
		i.addStringToBuffer(buffer, i.CallerPath)
	}
	if i.Prefix != "" {
		i.addStringToBuffer(buffer, i.Prefix)
	}
	if i.CtxStr != "" {
		i.addStringToBuffer(buffer, i.CtxStr)
	}
	if i.Content != "" {
		i.addStringToBuffer(buffer, i.Content)
	}
	i.addStringToBuffer(buffer, "\n")
	return buffer
}

// getLevelFormatWithColor returns the prefix string with color.
func (i *HandlerInput) getLevelFormatWithColor() string {
	s := i.LevelFormat
	color := defaultLevelColor[i.Level]
	if i.logger.config.color != 0 {
		color = i.logger.config.color
	}
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, s)
}

func (i *HandlerInput) String() string {
	return i.Buffer().String()
}

func (i *HandlerInput) Next() {
	if len(i.logger.config.Handlers)-1 > i.index {
		i.index++
		i.logger.config.Handlers[i.index](i.Ctx, i)
	}
}
