// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc

func (s *SearchInput) Key() string {
	keyPrefix := ""
	if s.Prefix != "" {
		keyPrefix += "/" + s.Prefix
	}
	if s.Deployment != "" {
		keyPrefix += "/" + s.Deployment
		if s.Namespace != "" {
			keyPrefix += "/" + s.Namespace
			if s.Name != "" {
				keyPrefix += "/" + s.Name
				if s.Version != "" {
					keyPrefix += "/" + s.Version
				}
			}
		}
	}
	return keyPrefix
}
