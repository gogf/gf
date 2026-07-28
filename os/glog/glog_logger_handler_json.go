// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"context"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

// reservedJsonHandlerKeys defines the reserved keys that cannot be overwritten by context values or content flattening.
var reservedJsonHandlerKeys = gset.NewStrSetFrom([]string{
	"Time",
	"TraceId",
	"CtxStr",
	"Level",
	"CallerFunc",
	"CallerPath",
	"Prefix",
	"Content",
	"Stack",
})

// HandlerJson is a handler for output logging content as a single json string.
func HandlerJson(ctx context.Context, in *HandlerInput) {
	output := gmap.NewStrAnyMap()
	output.Set("Time", in.TimeFormat)
	setJsonHandlerOutputValue(output, "TraceId", in.TraceId)
	setJsonHandlerOutputCtxValues(ctx, in, output)
	output.Set("Level", in.LevelFormat)
	setJsonHandlerOutputValue(output, "CallerFunc", in.CallerFunc)
	setJsonHandlerOutputValue(output, "CallerPath", in.CallerPath)
	setJsonHandlerOutputValue(output, "Prefix", in.Prefix)
	// Deprecated: CtxStr will be removed in the next major version.
	// Context values are now exported as individual top-level JSON fields.
	setJsonHandlerOutputValue(output, "CtxStr", in.CtxStr)
	setJsonHandlerOutputContent(output, in)
	setJsonHandlerOutputValue(output, "Stack", in.Stack)

	// Output json content.
	jsonBytes, err := json.Marshal(output.Map())
	if err != nil {
		panic(err)
	}
	in.Buffer.Write(jsonBytes)
	in.Buffer.Write([]byte("\n"))
	in.Next(ctx)
}

// setJsonHandlerOutputValue sets non-empty value to handler output.
func setJsonHandlerOutputValue(output *gmap.StrAnyMap, key string, value any) {
	if gconv.String(value) != "" {
		output.Set(key, value)
	}
}

// setJsonHandlerOutputCtxValues appends configured context values to handler output.
// It skips reserved keys to prevent overwriting core metadata fields.
func setJsonHandlerOutputCtxValues(ctx context.Context, in *HandlerInput, output *gmap.StrAnyMap) {
	if ctx == nil || in.Logger == nil {
		return
	}
	for _, ctxKey := range in.Logger.GetCtxKeys() {
		key := gconv.String(ctxKey)
		// Skip reserved keys to prevent overwriting core metadata fields.
		if reservedJsonHandlerKeys.Contains(key) {
			continue
		}
		ctxValue := ctx.Value(ctxKey)
		if ctxValue == nil {
			ctxValue = ctx.Value(gctx.StrKey(key))
		}
		if ctxValue != nil {
			output.Set(key, ctxValue)
		}
	}
}

// setJsonHandlerOutputContent appends logging content or flattens JSON object content.
// It skips reserved keys to prevent injecting/spoofing core metadata fields.
func setJsonHandlerOutputContent(output *gmap.StrAnyMap, in *HandlerInput) {
	content := in.Content
	if len(in.Values) > 0 {
		if content != "" {
			content += " "
		}
		content += in.ValuesContent()
	}

	var contentMap map[string]any
	if content != "" && json.Unmarshal([]byte(content), &contentMap) == nil {
		for key, value := range contentMap {
			// Skip reserved keys to prevent injecting/spoofing core metadata fields.
			if reservedJsonHandlerKeys.Contains(key) {
				continue
			}
			// Skip keys that already exist in output (e.g., from context values).
			if output.Get(key) != nil {
				continue
			}
			output.Set(key, value)
		}
		return
	}
	output.Set("Content", content)
}
