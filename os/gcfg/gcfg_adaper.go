// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import "context"

// Adapter is the interface for configuration retrieving.
type Adapter interface {
	// Available checks and returns the backend configuration service is available.
	// The optional parameter `resource` specifies certain configuration resource.
	//
	// Note that this function does not return error as it just does simply check for
	// backend configuration service.
	Available(ctx context.Context, resource ...string) (ok bool)

	// Get retrieves and returns value by specified `pattern` in current resource.
	// Pattern like:
	// "x.y.z" for map item.
	// "x.0.y" for slice item.
	Get(ctx context.Context, pattern string) (value interface{}, err error)

	// Data retrieves and returns all configuration data in current resource as map.
	// Note that this function may lead lots of memory usage if configuration data is too large,
	// you can implement this function if necessary.
	Data(ctx context.Context) (data map[string]interface{}, err error)
}
