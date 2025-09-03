package main

import (
	"context"
	"queue-service/bootstrap"
	"queue-service/infrastructure/discovery"
	grpcservice "queue-service/infrastructure/discovery/grpc_service"
	"queue-service/infrastructure/grpc_client"
	"queue-service/infrastructure/handler_mail"
	logger "queue-service/infrastructure/service/logger"

	"github.com/hibiken/asynq"
	"go.uber.org/zap/zapcore"
)

func main() {
	var env = bootstrap.Env{}
	bootstrap.NewEnv(&env)
	logConfig := logger.NewConfig()
	log := logger.InitLogger(logConfig, zapcore.DebugLevel, env.IsProduction())

	discovery, err := discovery.NewDiscovery(log, &env)
	if err != nil {
		log.Fatal("Không thể khởi động discovery: " + err.Error())
	}
	discovery.Register(env.NameService)
	clientFactory := grpc_client.NewClientFactory(log, env.GrpcClients)
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

	go func() {
		grpcSrv := grpcservice.NewGRPCServer(&env, log)
		if err := grpcSrv.Start(context.Background()); err != nil {
			log.Fatal("Không thể khởi động gRPC server: " + err.Error())
		}
	}()

	if err := srv.Run(mux); err != nil {
		log.Fatal("Không thể khởi động máy chủ: " + err.Error())
	}
}
