package grpcservice

import (
	"context"
	"fmt"
	"net"

	"queue-service/bootstrap"
	loggerI "queue-service/domain/service/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server    *grpc.Server
	healthSrv *health.Server
	env       *bootstrap.Env
	log       loggerI.Log
}

func NewGRPCServer(env *bootstrap.Env, log loggerI.Log) *GRPCServer {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggingInterceptor(log),
		),
	)
	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthSrv)

	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthSrv.SetServingStatus(env.NameService, grpc_health_v1.HealthCheckResponse_SERVING)

	if !env.IsProduction() {
		log.Info("Reflection is enabled")
		reflection.Register(server)
	}

	return &GRPCServer{
		server:    server,
		healthSrv: healthSrv,
		log:       log,
		env:       env,
	}
}

func (s *GRPCServer) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.env.PortGrpc))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	go func() {
		s.log.Info(fmt.Sprintf("gRPC server starting to serve on port %d", s.env.PortGrpc))
		if err := s.server.Serve(lis); err != nil {
			s.log.Error(fmt.Sprintf("gRPC server failed to serve: %v", err))
		}
	}()

	<-ctx.Done()

	s.log.Info("Shutting down gRPC server...")
	s.server.GracefulStop()

	return nil
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

func (s *GRPCServer) SetHealthStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	if s.healthSrv != nil {
		s.healthSrv.SetServingStatus(service, status)
	}
}
