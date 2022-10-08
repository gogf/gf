// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"context"
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
}

// Service implements gsvc.Service interface.
type Service struct {
	gsvc.Service
	Endpoints gsvc.Endpoints
}

// NewGrpcServer creates and returns a grpc server.
func (s modServer) NewGrpcServer(conf ...*GrpcServerConfig) *GrpcServer {
	var (
		ctx    = context.TODO()
		config *GrpcServerConfig
	)
	if len(conf) > 0 {
		config = conf[0]
	} else {
		config = s.NewGrpcServerConfig()
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
		config: config,
	}
	grpcServer.config.Options = append([]grpc.ServerOption{
		s.ChainUnary(
			s.UnaryRecover,
			s.UnaryTracing,
			s.UnaryError,
			grpcServer.UnaryLogger,
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
	var (
		serviceAddress string
		ctx            = context.TODO()
		array          = gstr.Split(s.config.Address, ":")
	)
	if array[0] == "0.0.0.0" || array[0] == "" {
		intraIP, err := gipv4.GetIntranetIp()
		if err != nil {
			s.config.Logger.Fatalf(
				ctx,
				"retrieving intranet ip failed, please check your net card or manually assign the service address: %+v",
				err,
			)
		}
		serviceAddress = fmt.Sprintf(`%s:%s`, intraIP, array[1])
	} else {
		serviceAddress = s.config.Address
	}
	for i, service := range services {
		if len(service.GetEndpoints()) == 0 {
			newService := &Service{
				Service:   service,
				Endpoints: make(gsvc.Endpoints, 0),
			}
			newService.Endpoints = gsvc.NewEndpoints(serviceAddress)
			services[i] = newService
		}
	}
	s.services = append(s.services, services...)
}

// Run starts the server in blocking way.
func (s *GrpcServer) Run() {
	autoLoadAndRegisterEtcdRegistry()

	var (
		err error
		ctx = context.TODO()
	)
	// Initialize services configured.
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
		ctx     = context.Background()
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
	if gsvc.GetRegistry() == nil {
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
		ctx      = context.Background()
		protocol = `grpc`
		address  = s.GetListenedAddress()
		array    = gstr.Split(address, ":")
		ip       = array[0]
		port     = array[1]
	)
	if ip == "" {
		ip = gipv4.MustGetIntranetIp()
	}
	// Register service list after server starts.
	for i, service := range s.services {
		service = &gsvc.LocalService{
			Name:      service.GetName(),
			Endpoints: gsvc.NewEndpoints(fmt.Sprintf(`%s:%s`, ip, port)),
			Metadata:  service.GetMetadata(),
		}
		service.GetMetadata().Sets(gsvc.Metadata{
			gsvc.MDProtocol: protocol,
		})
		s.config.Logger.Debugf(ctx, `service register: %+v`, service)
		if service, err = gsvc.Register(ctx, service); err != nil {
			s.config.Logger.Fatalf(ctx, `%+v`, err)
		}
		s.services[i] = service
	}
}

// doServiceDeregister de-registers current service from Registry.
func (s *GrpcServer) doServiceDeregister() {
	if gsvc.GetRegistry() == nil {
		return
	}
	var ctx = context.Background()
	for _, service := range s.services {
		s.config.Logger.Debugf(ctx, `service deregister: %+v`, service)
		if err := gsvc.Deregister(ctx, service); err != nil {
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
