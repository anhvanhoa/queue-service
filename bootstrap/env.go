package bootstrap

import (
	"log"
	"os"
	"path/filepath"
	"queue-service/infrastructure/grpc_client"
	"strings"

	"github.com/spf13/viper"
)

type queue struct {
	Addr        string
	DB          int
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Network     string
	Concurrency int
	Queues      map[string]int
}

type Env struct {
	MODE_ENV          string
	NAME_SERVICE      string
	HOST_GRPC         string
	PORT_GRPC         int
	INTERVAL_CHECK    string
	TIMEOUT_CHECK     string
	QUEUE             *queue
	MAIL_SERVICE_ADDR string
	GRPC_CLIENTS      []*grpc_client.Config
}

func NewEnv(env *Env) {
	absPath, err := filepath.Abs("./")
	if err != nil {
		log.Fatal("Lỗi khi lấy đường dẫn tuyệt đối:", err)
	}

	mode := os.Getenv("ENV_MODE")
	viper.SetConfigType("yaml")
	if mode == "production" {
		viper.SetConfigName("prod.config")
	} else {
		viper.SetConfigName("dev.config")
	}
	viper.AddConfigPath(absPath)
	err = viper.ReadInConfig()
	if err != nil {
		panic("Lỗi khi đọc file cấu hình, " + err.Error())
	}

	err = viper.UnmarshalExact(env)
	if err != nil {
		panic("Lỗi khi phân tích file cấu hình, " + err.Error())
	}
}

func (env *Env) IsProduction() bool {
	return strings.ToLower(env.MODE_ENV) == "production"
}
