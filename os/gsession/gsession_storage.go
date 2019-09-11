// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

type Storage interface {
	// Get returns the session data bytes for given session id.
	Get(id string) map[string]interface{}
	// Set updates the content for session id.
	// Note that the parameter <content> is the serialized bytes for session map.
	Set(id string, data map[string]interface{}) error
	// UpdateTTL updates the TTL for session id.
	UpdateTTL(id string) error
}
