// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import (
	"sync"

	"github.com/gogf/gf/v2/text/gregex"
)

// SetGlobalAttributesOption binds the global attributes to certain instrument.
type SetGlobalAttributesOption struct {
	// Instrument specifies the instrument name.
	Instrument string

	// Instrument specifies the instrument version.
	InstrumentVersion string

	// InstrumentPattern specifies instrument by regular expression on Instrument name.
	// Example:
	// 1. given `.+` will match all instruments.
	// 2. given `github.com/gogf/gf.+` will match all goframe instruments.
	InstrumentPattern string
}

// GetGlobalAttributesOption binds the global attributes to certain instrument.
type GetGlobalAttributesOption struct {
	Instrument        string // Instrument specifies the instrument name.
	InstrumentVersion string // Instrument specifies the instrument version.
}

type globalAttributeItem struct {
	Attributes
	SetGlobalAttributesOption
}

var (
	globalAttributesMu sync.Mutex
	// globalAttributes stores the global attributes to a map.
	globalAttributes = make([]globalAttributeItem, 0)
)

// SetGlobalAttributes appends global attributes according `SetGlobalAttributesOption`.
// It appends global attributes to all metrics if given `SetGlobalAttributesOption` is empty.
// It appends global attributes to certain instrument by given `SetGlobalAttributesOption`.
func SetGlobalAttributes(attrs Attributes, option SetGlobalAttributesOption) {
	globalAttributesMu.Lock()
	defer globalAttributesMu.Unlock()
	globalAttributes = append(
		globalAttributes, globalAttributeItem{
			Attributes:                attrs,
			SetGlobalAttributesOption: option,
		},
	)
}

// GetGlobalAttributes retrieves and returns the global attributes by `GetGlobalAttributesOption`.
// It returns the global attributes if given `GetGlobalAttributesOption` is empty.
// It returns global attributes of certain instrument if `GetGlobalAttributesOption` is not empty.
func GetGlobalAttributes(option GetGlobalAttributesOption) Attributes {
	globalAttributesMu.Lock()
	defer globalAttributesMu.Unlock()
	var attributes = make(Attributes, 0)
	for _, attrItem := range globalAttributes {
		// instrument name.
		if attrItem.InstrumentPattern != "" {
			if !gregex.IsMatchString(attrItem.InstrumentPattern, option.Instrument) {
				continue
			}
		} else {
			if (attrItem.Instrument != "" || option.Instrument != "") &&
				attrItem.Instrument != option.Instrument {
				continue
			}
		}
		// instrument version.
		if (attrItem.InstrumentVersion != "" || option.InstrumentVersion != "") &&
			attrItem.InstrumentVersion != option.InstrumentVersion {
			continue
		}
		attributes = append(attributes, attrItem.Attributes...)
	}
	return attributes
}
