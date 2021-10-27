// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import "context"

// Adapter is the interface for configuration retrieving.
type Adapter interface {
	// Get retrieves and returns value by specified `pattern`.
	Get(ctx context.Context, pattern string) (value interface{}, err error)

	// Data retrieves and returns all configuration data as map type.
	// Note that this function may lead lots of memory usage if configuration data is too large,
	// you can implement this function if necessary.
	Data(ctx context.Context) (data map[string]interface{}, err error)
}
