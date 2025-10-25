package main

import (
	"context"
	"queue-service/bootstrap"
	"queue-service/infrastructure/grpc_client"
	"queue-service/infrastructure/handler_mail"

	grpc_service "github.com/anhvanhoa/service-core/bootstrap/grpc"
	"google.golang.org/grpc"

	gc "github.com/anhvanhoa/service-core/domain/grpc_client"
	"github.com/anhvanhoa/service-core/domain/log"
	"github.com/hibiken/asynq"
	"go.uber.org/zap/zapcore"
)

func main() {
	var env = bootstrap.Env{}
	logConfig := log.NewConfig()
	log := log.InitLogGRPC(logConfig, zapcore.DebugLevel, env.IsProduction())
	bootstrap.NewEnv(&env)
	clientFactory := gc.NewClientFactory(env.GrpcClients...)
	mailService := grpc_client.NewMailService(clientFactory.GetClient(env.MailServiceAddr))

	cf := asynq.Config{
		Concurrency: env.Queue.Concurrency,
		Queues:      env.Queue.Queues,
	}
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     env.Queue.Addr,
			DB:       env.Queue.Db,
			Password: env.Queue.Password,
			Network:  env.Queue.Network,
		},
		cf,
	)
	mux := asynq.NewServeMux()
	handler_mail.NewEmailHandler(mux, &env, log, mailService)
	config := &grpc_service.GRPCServerConfig{
		IsProduction: env.IsProduction(),
		PortGRPC:     env.PortGrpc,
		NameService:  env.NameService,
	}
	go func() {
		grpcSrv := grpc_service.NewGRPCServer(config, log, func(server *grpc.Server) {})
		if err := grpcSrv.Start(context.Background()); err != nil {
			log.Fatal("Không thể khởi động gRPC server: " + err.Error())
		}
	}()

	if err := srv.Run(mux); err != nil {
		log.Fatal("Không thể khởi động máy chủ: " + err.Error())
	}
}
