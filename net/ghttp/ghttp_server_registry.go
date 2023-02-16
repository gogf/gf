// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"

	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// doServiceRegister registers current service to Registry.
func (s *Server) doServiceRegister() {
	if s.registrar == nil {
		return
	}
	var (
		ctx      = gctx.GetInitCtx()
		protocol = gsvc.DefaultProtocol
		insecure = true
		err      error
	)
	if s.config.TLSConfig != nil {
		protocol = `https`
		insecure = false
	}
	metadata := gsvc.Metadata{
		gsvc.MDProtocol: protocol,
		gsvc.MDInsecure: insecure,
	}
	s.service = &gsvc.LocalService{
		Name:      s.GetName(),
		Endpoints: s.calculateListenedEndpoints(),
		Metadata:  metadata,
	}
	s.Logger().Debugf(ctx, `service register: %+v`, s.service)
	if s.service, err = s.registrar.Register(ctx, s.service); err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}
}

// doServiceDeregister de-registers current service from Registry.
func (s *Server) doServiceDeregister() {
	if s.registrar == nil {
		return
	}
	var ctx = gctx.GetInitCtx()
	s.Logger().Debugf(ctx, `service deregister: %+v`, s.service)
	if err := s.registrar.Deregister(ctx, s.service); err != nil {
		s.Logger().Errorf(ctx, `%+v`, err)
	}
}

func (s *Server) calculateListenedEndpoints() gsvc.Endpoints {
	var (
		address       = s.config.Address
		endpoints     = make(gsvc.Endpoints, 0)
		listenedIps   []string
		listenedPorts []int
	)
	if address == "" {
		address = s.config.HTTPSAddr
	}
	var addrArray = gstr.Split(address, ":")
	switch addrArray[0] {
	case "0.0.0.0", "":
		listenedIps, _ = gipv4.GetIntranetIpArray()
	default:
		listenedIps = []string{addrArray[0]}
	}
	switch addrArray[1] {
	case "0":
		listenedPorts = s.GetListenedPorts()
	default:
		listenedPorts = []int{gconv.Int(addrArray[1])}
	}
	for _, ip := range listenedIps {
		for _, port := range listenedPorts {
			endpoints = append(endpoints, gsvc.NewEndpoint(fmt.Sprintf(`%s:%d`, ip, port)))
		}
	}
	return endpoints
}
