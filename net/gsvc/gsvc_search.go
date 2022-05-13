// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc

// Key formats and returns a string for prefix searching purpose.
func (s *SearchInput) Key() string {
	separator := DefaultSeparator
	if s.Separator != "" {
		separator = s.Separator
	}
	keyPrefix := ""
	if s.Prefix != "" {
		if s.Separator == DefaultSeparator {
			keyPrefix += separator + s.Prefix
		} else {
			keyPrefix += s.Prefix
		}
	}
	if s.Deployment != "" {
		keyPrefix += separator + s.Deployment
		if s.Namespace != "" {
			keyPrefix += separator + s.Namespace
			if s.Name != "" {
				keyPrefix += separator + s.Name
				if s.Version != "" {
					keyPrefix += separator + s.Version
				}
			}
		}
	}

	return keyPrefix
}
