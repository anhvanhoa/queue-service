package main

import (
	"queue-service/bootstrap"
	"queue-service/constants"
	"queue-service/infrastructure/grpc_client"
	"queue-service/infrastructure/handler_mail"
	pkglog "queue-service/infrastructure/service/logger"

	"github.com/hibiken/asynq"
	"go.uber.org/zap/zapcore"
)

func main() {
	var env = bootstrap.Env{}
	bootstrap.NewEnv(&env)
	logConfig := pkglog.NewConfig()
	log := pkglog.InitLogger(logConfig, zapcore.DebugLevel, env.IsProduction())
	clientConfig := []*grpc_client.Config{}
	for name, client := range env.GrpcClients {
		clientConfig = append(clientConfig, &grpc_client.Config{
			Name:          name,
			ServerAddress: client.ServerAddress,
			Timeout:       client.Timeout,
			MaxRetries:    client.MaxRetries,
			KeepAlive:     client.KeepAlive,
		})
	}
	clientFactory := grpc_client.NewClientFactory(log, clientConfig...)
	mailService := grpc_client.NewMailService(clientFactory.GetClient(constants.MailService))
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

	if err := srv.Run(mux); err != nil {
		log.Fatal("Không thể khởi động máy chủ: " + err.Error())
	}
}
