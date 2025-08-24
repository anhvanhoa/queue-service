package main

import (
	"queue-service/bootstrap"
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
	for name, client := range env.GRPC_CLIENTS {
		clientConfig = append(clientConfig, &grpc_client.Config{
			Name:          name,
			ServerAddress: client.ServerAddress,
			Timeout:       client.Timeout,
			MaxRetries:    client.MaxRetries,
			KeepAlive:     client.KeepAlive,
		})
	}
	clientFactory := grpc_client.NewClientFactory(clientConfig...)
	mailService := grpc_client.NewMailService(clientFactory.GetClient("mail_service"))
	cf := asynq.Config{
		Concurrency: env.QUEUE.Concurrency,
		Queues:      env.QUEUE.Queues,
	}
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     env.QUEUE.Addr,
			DB:       env.QUEUE.DB,
			Password: env.QUEUE.Password,
			Network:  env.QUEUE.Network,
		},
		cf,
	)
	mux := asynq.NewServeMux()
	handler_mail.NewEmailHandler(mux, &env, log, mailService)

	if err := srv.Run(mux); err != nil {
		log.Fatal("Không thể khởi động máy chủ: " + err.Error())
	}
}
