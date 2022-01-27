// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
)

// doServiceRegister registers current service to Registry.
func (s *Server) doServiceRegister() {
	var (
		ctx      = context.Background()
		protocol = `http`
		insecure = true
		address  = s.config.Address
	)
	if address == "" {
		address = s.config.HTTPSAddr
	}
	var (
		array = gstr.Split(address, ":")
		ip    = array[0]
		port  = array[1]
	)
	if ip == "" {
		ip = gipv4.MustGetIntranetIp()
	}
	if s.config.TLSConfig != nil {
		protocol = `https`
		insecure = false
	}
	metadata := gsvc.Metadata{
		gsvc.MDProtocol: protocol,
		gsvc.MDInsecure: insecure,
	}
	s.service = &gsvc.Service{
		Name:      s.GetName(),
		Endpoints: []string{fmt.Sprintf(`%s:%s`, ip, port)},
		Metadata:  metadata,
	}
	_ = gsvc.Register(ctx, s.service)
}

// doServiceDeregister de-registers current service from Registry.
func (s *Server) doServiceDeregister() {
	_ = gsvc.Deregister(context.Background(), s.service)
}
