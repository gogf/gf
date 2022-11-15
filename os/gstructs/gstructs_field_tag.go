// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstructs

import (
	"strings"

	"github.com/gogf/gf/v2/util/gtag"
)

// TagJsonName returns the `json` tag name string of the field.
func (f *Field) TagJsonName() string {
	if jsonTag := f.Tag(gtag.Json); jsonTag != "" {
		return strings.Split(jsonTag, ",")[0]
	}
	return ""
}

// TagDefault returns the most commonly used tag `default/d` value of the field.
func (f *Field) TagDefault() string {
	v := f.Tag(gtag.Default)
	if v == "" {
		v = f.Tag(gtag.DefaultShort)
	}
	return v
}

// TagParam returns the most commonly used tag `param/p` value of the field.
func (f *Field) TagParam() string {
	v := f.Tag(gtag.Param)
	if v == "" {
		v = f.Tag(gtag.ParamShort)
	}
	return v
}

// TagValid returns the most commonly used tag `valid/v` value of the field.
func (f *Field) TagValid() string {
	v := f.Tag(gtag.Valid)
	if v == "" {
		v = f.Tag(gtag.ValidShort)
	}
	return v
}

// TagDescription returns the most commonly used tag `description/des/dc` value of the field.
func (f *Field) TagDescription() string {
	v := f.Tag(gtag.Description)
	if v == "" {
		v = f.Tag(gtag.DescriptionShort)
	}
	if v == "" {
		v = f.Tag(gtag.DescriptionShort2)
	}
	return v
}

// TagSummary returns the most commonly used tag `summary/sum/sm` value of the field.
func (f *Field) TagSummary() string {
	v := f.Tag(gtag.Summary)
	if v == "" {
		v = f.Tag(gtag.SummaryShort)
	}
	if v == "" {
		v = f.Tag(gtag.SummaryShort2)
	}
	return v
}

// TagAdditional returns the most commonly used tag `additional/ad` value of the field.
func (f *Field) TagAdditional() string {
	v := f.Tag(gtag.Additional)
	if v == "" {
		v = f.Tag(gtag.AdditionalShort)
	}
	return v
}

// TagExample returns the most commonly used tag `example/eg` value of the field.
func (f *Field) TagExample() string {
	v := f.Tag(gtag.Example)
	if v == "" {
		v = f.Tag(gtag.ExampleShort)
	}
	return v
}
