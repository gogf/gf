// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
)

// GrpcServer is the server for GRPC protocol.
type GrpcServer struct {
	Server    *grpc.Server
	config    *GrpcServerConfig
	listener  net.Listener
	services  []gsvc.Service
	waitGroup sync.WaitGroup
	registrar gsvc.Registrar
}

// Service implements gsvc.Service interface.
type Service struct {
	gsvc.Service
	Endpoints gsvc.Endpoints
}

// New creates and returns a grpc server.
func (s modServer) New(conf ...*GrpcServerConfig) *GrpcServer {
	autoLoadAndRegisterFileRegistry()

	var (
		ctx    = gctx.GetInitCtx()
		config *GrpcServerConfig
	)
	if len(conf) > 0 {
		config = conf[0]
	} else {
		config = s.NewConfig()
	}
	if config.Address == "" {
		randomPort, err := gtcp.GetFreePort()
		if err != nil {
			g.Log().Fatalf(ctx, `%+v`, err)
		}
		config.Address = fmt.Sprintf(`:%d`, randomPort)
	}
	if !gstr.Contains(config.Address, ":") {
		g.Log().Fatal(ctx, "invalid service address, should contain listening port")
	}
	if config.Logger == nil {
		config.Logger = glog.New()
	}
	grpcServer := &GrpcServer{
		config:    config,
		registrar: gsvc.GetRegistry(),
	}
	grpcServer.config.Options = append([]grpc.ServerOption{
		s.ChainUnary(
			s.UnaryRecover,
			s.UnaryTracing,
			s.UnaryError,
			grpcServer.internalUnaryLogger,
		),
		s.ChainStream(
			s.StreamTracing,
		),
	}, grpcServer.config.Options...)
	grpcServer.Server = grpc.NewServer(grpcServer.config.Options...)
	return grpcServer
}

// Service binds service list to current server.
// Server will automatically register the service list after it starts.
func (s *GrpcServer) Service(services ...gsvc.Service) {
	s.services = append(s.services, services...)
}

// Run starts the server in blocking way.
func (s *GrpcServer) Run() {
	var (
		err error
		ctx = gctx.GetInitCtx()
	)
	// Create listener to bind listening ip and port.
	s.listener, err = net.Listen("tcp", s.config.Address)
	if err != nil {
		s.config.Logger.Fatalf(ctx, `%+v`, err)
	}

	// Start listening.
	go func() {
		if err = s.Server.Serve(s.listener); err != nil {
			s.config.Logger.Fatalf(ctx, `%+v`, err)
		}
	}()

	// Service register.
	s.doServiceRegister()
	s.config.Logger.Infof(
		ctx,
		"pid[%d]: grpc server started listening on [%s]",
		gproc.Pid(), s.GetListenedAddress(),
	)
	s.doSignalListen()
}

// doSignalListen does signal listening and handling for gracefully shutdown.
func (s *GrpcServer) doSignalListen() {
	var (
		ctx     = gctx.GetInitCtx()
		sigChan = make(chan os.Signal, 1)
	)
	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGABRT,
	)
	for {
		sig := <-sigChan
		switch sig {
		case
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGABRT:
			s.config.Logger.Infof(ctx, "signal received: %s, gracefully shutting down", sig.String())
			s.doServiceDeregister()
			time.Sleep(time.Second)
			s.Stop()
			return
		}
	}
}

// doServiceRegister registers current service to Registry.
func (s *GrpcServer) doServiceRegister() {
	if s.registrar == nil {
		return
	}
	if len(s.services) == 0 {
		s.services = []gsvc.Service{&gsvc.LocalService{
			Name:     s.config.Name,
			Metadata: gsvc.Metadata{},
		}}
	}
	var (
		err      error
		ctx      = gctx.GetInitCtx()
		protocol = `grpc`
	)
	// Register service list after server starts.
	for i, service := range s.services {
		service = &gsvc.LocalService{
			Name:      service.GetName(),
			Endpoints: s.calculateListenedEndpoints(),
			Metadata:  service.GetMetadata(),
		}
		service.GetMetadata().Sets(gsvc.Metadata{
			gsvc.MDProtocol: protocol,
		})
		s.config.Logger.Debugf(ctx, `service register: %+v`, service)
		if service, err = s.registrar.Register(ctx, service); err != nil {
			s.config.Logger.Fatalf(ctx, `%+v`, err)
		}
		s.services[i] = service
	}
}

// doServiceDeregister de-registers current service from Registry.
func (s *GrpcServer) doServiceDeregister() {
	if s.registrar == nil {
		return
	}
	var ctx = gctx.GetInitCtx()
	for _, service := range s.services {
		s.config.Logger.Debugf(ctx, `service deregister: %+v`, service)
		if err := s.registrar.Deregister(ctx, service); err != nil {
			s.config.Logger.Errorf(ctx, `%+v`, err)
		}
	}
}

// Start starts the server in no-blocking way.
func (s *GrpcServer) Start() {
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		s.Run()
	}()
}

// Wait works with Start, which blocks current goroutine until the server stops.
func (s *GrpcServer) Wait() {
	s.waitGroup.Wait()
}

// Stop gracefully stops the server.
func (s *GrpcServer) Stop() {
	s.Server.GracefulStop()
}

// GetListenedAddress retrieves and returns the address string which are listened by current server.
func (s *GrpcServer) GetListenedAddress() string {
	if !gstr.Contains(s.config.Address, FreePortAddress) {
		return s.config.Address
	}
	var (
		address      = s.config.Address
		listenedPort = s.GetListenedPort()
	)
	address = gstr.Replace(address, FreePortAddress, fmt.Sprintf(`:%d`, listenedPort))
	return address
}

// GetListenedPort retrieves and returns one port which is listened to by current server.
func (s *GrpcServer) GetListenedPort() int {
	if ln := s.listener; ln != nil {
		return ln.Addr().(*net.TCPAddr).Port
	}
	return -1
}

func (s *GrpcServer) calculateListenedEndpoints() gsvc.Endpoints {
	var (
		address      = s.config.Address
		endpoints    = make(gsvc.Endpoints, 0)
		listenedPort = s.GetListenedPort()
		listenedIps  []string
	)
	var addrArray = gstr.Split(address, ":")
	switch addrArray[0] {
	case "0.0.0.0", "":
		listenedIps = []string{gipv4.MustGetIntranetIp()}
	default:
		listenedIps = []string{addrArray[0]}
	}
	for _, ip := range listenedIps {
		endpoints = append(endpoints, gsvc.NewEndpoint(fmt.Sprintf(`%s:%d`, ip, listenedPort)))
	}
	return endpoints
}
