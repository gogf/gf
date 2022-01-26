// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	DefaultPrefix     = `goframe`
	DefaultDeployment = `default`
	DefaultNamespace  = `default`
	DefaultVersion    = `latest`
)

// NewServiceFromKV creates and returns service from `key` and `value`.
func NewServiceFromKV(key, value []byte) (s *Service, err error) {
	array := gstr.Split(string(key), "/")
	if len(array) < 6 {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid service key "%s"`, key)
	}
	s = &Service{
		Prefix:     array[0],
		Deployment: array[1],
		Namespace:  array[2],
		Name:       array[3],
		Version:    array[4],
		Address:    array[5],
		Metadata:   make(map[string]interface{}),
	}
	s.autoFillDefaultAttributes()
	if len(value) > 0 {
		if err = gjson.Unmarshal(value, &s.Metadata); err != nil {
			return nil, gerror.WrapCodef(gcode.CodeInvalidParameter, err, `invalid service value "%s"`, value)
		}
	}
	return s, nil
}

// Key formats the service information and returns the Service as registering key.
func (s *Service) Key() string {
	s.autoFillDefaultAttributes()
	return gstr.Join([]string{
		s.Prefix,
		s.Deployment,
		s.Namespace,
		s.Name,
		s.Version,
		s.Address,
	}, "/")
}

func (s *Service) Value() string {
	b, err := gjson.Marshal(s.Metadata)
	if err != nil {
		intlog.Error(context.TODO(), err)
	}
	return string(b)
}

func (s *Service) autoFillDefaultAttributes() {
	if s.Prefix == "" {
		s.Prefix = DefaultPrefix
	}
	if s.Deployment == "" {
		s.Deployment = DefaultDeployment
	}
	if s.Namespace == "" {
		s.Namespace = DefaultNamespace
	}
	if s.Version == "" {
		s.Version = DefaultVersion
	}
}
