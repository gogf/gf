// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

// SetConfig set configurations for the logger.
func SetConfig(config Config) error {
	return logger.SetConfig(config)
}

// SetConfigWithMap set configurations with map for the logger.
func SetConfigWithMap(m map[string]interface{}) error {
	return logger.SetConfigWithMap(m)
}
