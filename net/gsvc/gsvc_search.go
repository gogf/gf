// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc

// Key formats and returns a string for prefix searching purpose.
func (s *SearchInput) Key() string {
	if s.Separator == "" {
		s.Separator = defaultSeparator
	}
	keyPrefix := ""
	if s.Prefix != "" {
		if s.Separator == defaultSeparator {
			keyPrefix += s.Separator + s.Prefix
		} else {
			keyPrefix += s.Prefix
		}
	}
	if s.Deployment != "" {
		keyPrefix += s.Separator + s.Deployment
		if s.Namespace != "" {
			keyPrefix += s.Separator + s.Namespace
			if s.Name != "" {
				keyPrefix += s.Separator + s.Name
				if s.Version != "" {
					keyPrefix += s.Separator + s.Version
				}
			}
		}
	}

	return keyPrefix
}
