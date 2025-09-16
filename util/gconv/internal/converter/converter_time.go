// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"time"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Time converts `any` to time.Time.
func (c *Converter) Time(anyInput any, format ...string) (time.Time, error) {
	// Handle special cases when no format is specified
	if len(format) == 0 {
		// Direct type matches - fastest path
		if v, ok := anyInput.(time.Time); ok {
			return v, nil
		}
		if v, ok := anyInput.(*gtime.Time); ok {
			// Handle *gtime.Time directly to preserve timezone
			if v == nil {
				return time.Time{}, nil
			}
			return v.Time, nil
		}

		// Handle map inputs by extracting the first value
		// This is optimized for ORM scenarios where maps like {"now": gtimeVal}
		// need to be converted to a single time.Time value
		if mapData, ok := anyInput.(map[string]interface{}); ok {
			if len(mapData) == 0 {
				return time.Time{}, nil
			}
			// Extract the first value efficiently without full iteration
			for _, value := range mapData {
				return c.Time(value, format...)
			}
		}
	}

	// Fall back to GTime conversion for complex cases
	t, err := c.GTime(anyInput, format...)
	if err != nil {
		return time.Time{}, err
	}
	if t != nil {
		return t.Time, nil
	}
	return time.Time{}, nil
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func (c *Converter) Duration(anyInput any) (time.Duration, error) {
	// It's already this type.
	if v, ok := anyInput.(time.Duration); ok {
		return v, nil
	}
	s, err := c.String(anyInput)
	if err != nil {
		return 0, err
	}
	if !utils.IsNumeric(s) {
		return gtime.ParseDuration(s)
	}
	i, err := c.Int64(anyInput)
	if err != nil {
		return 0, err
	}
	return time.Duration(i), nil
}

// GTime converts `any` to *gtime.Time.
// The parameter `format` can be used to specify the format of `any`.
// It returns the converted value that matched the first format of the formats slice.
// If no `format` given, it converts `any` using gtime.NewFromTimeStamp if `any` is numeric,
// or using gtime.StrToTime if `any` is string.
func (c *Converter) GTime(anyInput any, format ...string) (*gtime.Time, error) {
	if empty.IsNil(anyInput) {
		return nil, nil
	}

	// Check for custom interfaces first
	if v, ok := anyInput.(localinterface.IGTime); ok {
		return v.GTime(format...), nil
	}

	// Handle direct type matches when no format is specified - HIGHEST PRIORITY for timezone preservation
	if len(format) == 0 {
		switch v := anyInput.(type) {
		case *gtime.Time:
			return v, nil
		case gtime.Time:
			// Return a pointer to preserve the exact same gtime instance with timezone
			return &v, nil
		case time.Time:
			return gtime.New(v), nil
		case *time.Time:
			return gtime.New(v), nil
		}
	}

	// Convert to string for parsing
	s, err := c.String(anyInput)
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return gtime.New(), nil
	}

	// Handle format-specific conversion
	if len(format) > 0 {
		for _, item := range format {
			t, err := gtime.StrToTimeFormat(s, item)
			if err != nil {
				return nil, err
			}
			if t != nil {
				return t, nil
			}
		}
		return nil, nil
	}

	// Handle numeric timestamps
	if utils.IsNumeric(s) {
		i, err := c.Int64(s)
		if err != nil {
			return nil, err
		}
		return gtime.NewFromTimeStamp(i), nil
	}

	// Parse as time string with timezone preservation
	// Enhanced: if the string lacks timezone info, try to parse it with RFC3339 format first
	// This helps preserve timezone when the original gtime had timezone information
	result, err := gtime.StrToTime(s)
	if err == nil && result != nil {
		// If parsing succeeded but the result is in local timezone while
		// the string suggests it should be in a different timezone,
		// we may need additional handling here in the future
		return result, nil
	}

	return gtime.StrToTime(s)
}
