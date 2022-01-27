// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	separator = "/"
)

// NewServiceWithName creates and returns service from `name`.
func NewServiceWithName(name string) (s *Service) {
	s = &Service{
		Name:     name,
		Metadata: make(Metadata),
	}
	s.autoFillDefaultAttributes()
	return s
}

// NewServiceWithKV creates and returns service from `key` and `value`.
func NewServiceWithKV(key, value []byte) (s *Service, err error) {
	array := gstr.Split(gstr.Trim(string(key), separator), separator)
	if len(array) < 6 {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid service key "%s"`, key)
	}
	s = &Service{
		Prefix:     array[0],
		Deployment: array[1],
		Namespace:  array[2],
		Name:       array[3],
		Version:    array[4],
		Endpoints:  gstr.Split(array[5], ","),
		Metadata:   make(Metadata),
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
	serviceNameUnique := s.KeyWithoutEndpoints()
	serviceNameUnique += separator + gstr.Join(s.Endpoints, ",")
	return serviceNameUnique
}

// KeyWithSchema formats the service information and returns the Service as dialing target key.
func (s *Service) KeyWithSchema() string {
	return fmt.Sprintf(`%s://%s`, Schema, s.Key())
}

// KeyWithoutEndpoints formats the service information and returns a string as unique name of service.
func (s *Service) KeyWithoutEndpoints() string {
	s.autoFillDefaultAttributes()
	return "/" + gstr.Join([]string{
		s.Prefix,
		s.Deployment,
		s.Namespace,
		s.Name,
		s.Version,
	}, separator)
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
		s.Prefix = gcmd.GetOptWithEnv(EnvPrefix, DefaultPrefix).String()
	}
	if s.Deployment == "" {
		s.Deployment = gcmd.GetOptWithEnv(EnvDeployment, DefaultDeployment).String()
	}
	if s.Namespace == "" {
		s.Namespace = gcmd.GetOptWithEnv(EnvNamespace, DefaultNamespace).String()
	}
	if s.Version == "" {
		s.Version = gcmd.GetOptWithEnv(EnvVersion, DefaultVersion).String()
	}
	if s.Name == "" {
		s.Name = gcmd.GetOptWithEnv(EnvName).String()
	}
}
