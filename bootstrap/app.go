package bootstrap

import (
	loggerI "queue-service/domain/service/logger"
	logger "queue-service/infrastructure/service/logger"

	"github.com/go-pg/pg/v10"
	"go.uber.org/zap/zapcore"
)

type Application struct {
	Env *Env
	DB  *pg.DB
	Log loggerI.Log
}

func App() *Application {
	env := Env{}
	NewEnv(&env)

	logConfig := logger.NewConfig()
	log := logger.InitLogger(logConfig, zapcore.DebugLevel, env.IsProduction())

	db := NewPostgresDB(&env, log)
	return &Application{
		Env: &env,
		DB:  db,
		Log: log,
	}
}
