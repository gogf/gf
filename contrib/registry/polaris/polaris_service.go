// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
)

// Service for wrapping gsvc.Server and extends extra attributes for polaris purpose.
type Service struct {
	gsvc.Service        // Common service object.
	ID           string // ID is the unique instance ID as registered, for some registrar server.
}

// GetKey overwrites the GetKey function of gsvc.Service for replacing separator string.
func (s *Service) GetKey() string {
	key := s.Service.GetKey()
	key = gstr.Replace(key, gsvc.DefaultSeparator, instanceIDSeparator)
	key = gstr.TrimLeft(key, instanceIDSeparator)
	return key
}

// GetPrefix overwrites the GetPrefix function of gsvc.Service for replacing separator string.
func (s *Service) GetPrefix() string {
	prefix := s.Service.GetPrefix()
	prefix = gstr.Replace(prefix, gsvc.DefaultSeparator, instanceIDSeparator)
	prefix = gstr.TrimLeft(prefix, instanceIDSeparator)
	return prefix
}
